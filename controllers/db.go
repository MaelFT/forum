package forum

import (
    "database/sql"
    "errors"
    "fmt"
    "time"
    "golang.org/x/crypto/bcrypt"

    "github.com/mattn/go-sqlite3"
    models "forum/models"
)

var (
    ErrDuplicate    = errors.New("record already exists")
    ErrNotExists    = errors.New("row not exists")
    ErrUpdateFailed = errors.New("update failed")
    ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
    db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
    return &SQLiteRepository{
        db: db,
    }
}

/*---------------------
     Users Queries
----------------------*/

func (r *SQLiteRepository) TableUsers() error {
    query := `
    CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        mail TEXT NOT NULL UNIQUE,
        checkedmail INTEGER NOT NULL,
        password TEXT NOT NULL,
        categories TEXT,
        cookie TEXT NOT NULL UNIQUE,
        role TEXT,
        date DATE
    );
    `

    _, err := r.db.Exec(query)
    return err
}

func (r *SQLiteRepository) CreateUser(users models.Users) (*models.Users, error) {
    password := HashPassword(users.Password)
    res, err := r.db.Exec("INSERT INTO users(username, mail, checkedmail, password, categories, cookie, role, date) values(?,?,?,?,?,?,?,?)", 
    users.Username, users.Mail, users.CheckedMail, password, users.Categories, users.Cookie, users.Role, time.Now())
    if err != nil {
        var sqliteErr sqlite3.Error
        if errors.As(err, &sqliteErr) {
            if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
                return nil, ErrDuplicate
            }
        }
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    users.ID = id

    return &users, nil
}

func (r *SQLiteRepository) AllUsers() ([]models.Users, error) {
    rows, err := r.db.Query("SELECT * FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var all []models.Users
    for rows.Next() {
        var users models.Users
        if err := rows.Scan(&users.ID, &users.Username, &users.Mail, &users.CheckedMail, &users.Password, &users.Categories, &users.Cookie, &users.Role, &users.Date); err != nil {
            return nil, err
        }
        all = append(all, users)
    }
    return all, nil
}

func (r *SQLiteRepository) GetUserByCookie(c string) (*models.Users, error) {
    row := r.db.QueryRow("SELECT * FROM users WHERE cookie = ?", c)
    var users models.Users
    if err := row.Scan(&users.ID, &users.Username, &users.Mail, &users.CheckedMail, &users.Password, &users.Categories, &users.Cookie, &users.Role, &users.Date); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }
    return &users, nil
}

func (r *SQLiteRepository) CheckUser(username string, password string) (*models.Users, error) {
    row := r.db.QueryRow("SELECT * FROM users WHERE username = ?", username)
    var users models.Users
    if err := row.Scan(&users.ID, &users.Username, &users.Mail, &users.CheckedMail, &users.Password, &users.Categories, &users.Cookie, &users.Role, &users.Date); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }
    fmt.Println(users.Password, password, ComparePassword(users.Password, password))
    if ComparePassword(users.Password, password) != nil {
        return nil, ErrNotExists
    }
    return &users, nil
}

func (r *SQLiteRepository) UpdateUser(id int64, updated models.Users) (*models.Users, error) {
    if id == 0 {
        return nil, errors.New("invalid updated ID")
    }
    res, err := r.db.Exec("UPDATE users SET username = ?, mail = ?, checkedmail = ?, password = ?, categories = ?, cookie = ?, role = ?, date = ?, WHERE id = ?", updated.Username, updated.Mail, updated.CheckedMail, updated.Password, updated.Categories, updated.Cookie, updated.Role, updated.Date, id)
    if err != nil {
        return nil, err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return nil, err
    }

    if rowsAffected == 0 {
        return nil, ErrUpdateFailed
    }

    return &updated, nil
}

func (r *SQLiteRepository) DeleteUser(id int64) error {
    res, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrDeleteFailed
    }

    return err
}

func HashPassword(password string) string {
	pw := []byte(password)
	result, _ := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	return string(result) 
}

func ComparePassword(hashPassword string, password string) error {
	pw := []byte(password)
	hw := []byte(hashPassword)
	err := bcrypt.CompareHashAndPassword(hw, pw)
	return err
}

/*---------------------
     Posts Queries
----------------------*/

func (r *SQLiteRepository) TablePosts() error {
    query := `
    CREATE TABLE IF NOT EXISTS posts(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        content TEXT NOT NULL,
        categories TEXT,
        user_id INTEGER NOT NULL,
        date DATE,
        FOREIGN KEY (user_id)
            REFERENCES users (id)
    );
    `

    _, err := r.db.Exec(query)
    return err
}

func (r *SQLiteRepository) CreatePost(posts models.Posts) (*models.Posts, error) {
    res, err := r.db.Exec("INSERT INTO posts(name, content, categories, user_id, date) values(?,?,?,?,?)", 
    posts.Name, posts.Content, posts.Categories, posts.User_ID, time.Now())
    fmt.Println(res, err)
    if err != nil {
        var sqliteErr sqlite3.Error
        if errors.As(err, &sqliteErr) {
            if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
                return nil, ErrDuplicate
            }
        }
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    posts.ID = id

    return &posts, nil
}

func (r *SQLiteRepository) AllPosts() ([]models.Posts, error) {
    rows, err := r.db.Query("SELECT * FROM posts")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var all []models.Posts
    for rows.Next() {
        var posts models.Posts
        if err := rows.Scan(&posts.ID, &posts.Name, &posts.Content, &posts.Categories, &posts.User_ID, &posts.Date); err != nil {
            return nil, err
        }
        all = append(all, posts)
    }
    return all, nil
}

func (r *SQLiteRepository) AllPostsByDate() ([]models.Posts, error) {
    rows, err := r.db.Query("SELECT * FROM posts ORDER BY date")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var all []models.Posts
    for rows.Next() {
        var posts models.Posts
        if err := rows.Scan(&posts.ID, &posts.Name, &posts.Content, &posts.Categories, &posts.User_ID, &posts.Date); err != nil {
            return nil, err
        }
        all = append(all, posts)
    }
    return all, nil
}

func (r *SQLiteRepository) GetPostByID(id int64) (*models.Posts, error) {
    row := r.db.QueryRow("SELECT * FROM posts WHERE id = ?", id)

    var posts models.Posts
    if err := row.Scan(&posts.ID, &posts.Name, &posts.Content, &posts.Categories, &posts.User_ID, &posts.Date); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }
    return &posts, nil
}

func (r *SQLiteRepository) UpdatePost(id int64, updated models.Posts) (*models.Posts, error) {
    if id == 0 {
        return nil, errors.New("invalid updated ID")
    }
    res, err := r.db.Exec("UPDATE posts SET name = ?, content = ?, categories = ?, user_id = ?, date = ?, WHERE id = ?", updated.Name, updated.Content, updated.Categories, updated.User_ID, updated.Date, id)
    if err != nil {
        return nil, err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return nil, err
    }

    if rowsAffected == 0 {
        return nil, ErrUpdateFailed
    }

    return &updated, nil
}

func (r *SQLiteRepository) DeletePost(id int64) error {
    res, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrDeleteFailed
    }

    return err
}

/*---------------------
     Categories Queries
----------------------*/

func (r *SQLiteRepository) TableCategories() error {
    query := `
    CREATE TABLE IF NOT EXISTS categories(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        date DATE,
        FOREIGN KEY (user_id)
            REFERENCES users (id)
    );
    `

    _, err := r.db.Exec(query)
    return err
}

func (r *SQLiteRepository) CreateCategorie(category models.Categories) (*models.Categories, error) {
    res, err := r.db.Exec("INSERT INTO categories(title, description, user_id, date) values(?,?,?,?)", 
    category.Title, category.Description, category.User_ID, time.Now())
    if err != nil {
        var sqliteErr sqlite3.Error
        if errors.As(err, &sqliteErr) {
            if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
                return nil, ErrDuplicate
            }
        }
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    category.ID = id

    return &category, nil
}

func (r *SQLiteRepository) AllCategories() ([]models.Categories, error) {
    rows, err := r.db.Query("SELECT * FROM categories")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var all []models.Categories
    for rows.Next() {
        var category models.Categories
        if err := rows.Scan(&category.ID, &category.Title, &category.Description, &category.User_ID, &category.Date); err != nil {
            return nil, err
        }
        all = append(all, category)
    }
    return all, nil
}

func (r *SQLiteRepository) GetCategoryByID(id int64) (*models.Categories, error) {
    row := r.db.QueryRow("SELECT * FROM categories WHERE id = ?", id)

    var category models.Categories
    if err := row.Scan(&category.ID, &category.Title, &category.Description, &category.User_ID, &category.Date); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }
    return &category, nil
}

func (r *SQLiteRepository) UpdateCategory(id int64, updated models.Categories) (*models.Categories, error) {
    if id == 0 {
        return nil, errors.New("invalid updated ID")
    }
    res, err := r.db.Exec("UPDATE categories SET title = ?, description = ?, user_id = ?, date = ?, WHERE id = ?", updated.Title, updated.Description, updated.User_ID, updated.Date, id)
    if err != nil {
        return nil, err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return nil, err
    }

    if rowsAffected == 0 {
        return nil, ErrUpdateFailed
    }

    return &updated, nil
}

func (r *SQLiteRepository) DeleteCategorie(id int64) error {
    res, err := r.db.Exec("DELETE FROM categories WHERE id = ?", id)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrDeleteFailed
    }

    return err
}

/*---------------------
     Comments Queries
----------------------*/

func (r *SQLiteRepository) TableComments() error {
    query := `
    CREATE TABLE IF NOT EXISTS comments(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        date DATE,
        FOREIGN KEY (user_id)
            REFERENCES users (id),
        FOREIGN KEY (post_id)
            REFERENCES posts (id)
    );
    `

    _, err := r.db.Exec(query)
    return err
}

func (r *SQLiteRepository) CreateComment(comments models.Comments) (*models.Comments, error) {
    res, err := r.db.Exec("INSERT INTO comments(content, post_id, user_id, date) values(?,?,?,?)", 
    comments.Content, comments.Post_ID, comments.User_ID, time.Now())
    if err != nil {
        var sqliteErr sqlite3.Error
        if errors.As(err, &sqliteErr) {
            if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
                return nil, ErrDuplicate
            }
        }
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    comments.ID = id

    return &comments, nil
}

func (r *SQLiteRepository) AllComments() ([]models.Comments, error) {
    rows, err := r.db.Query("SELECT * FROM comments")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var all []models.Comments
    for rows.Next() {
        var comments models.Comments
        if err := rows.Scan(&comments.ID, &comments.Content, &comments.Post_ID, &comments.User_ID, &comments.Date); err != nil {
            return nil, err
        }
        all = append(all, comments)
    }
    return all, nil
}

func (r *SQLiteRepository) GetCommentByID(id int64) (*models.Comments, error) {
    row := r.db.QueryRow("SELECT * FROM comments WHERE id = ?", id)

    var comments models.Comments
    if err := row.Scan(&comments.ID, &comments.Content, &comments.Post_ID, &comments.User_ID, &comments.Date); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }
    return &comments, nil
}

func (r *SQLiteRepository) UpdateComment(id int64, updated models.Comments) (*models.Comments, error) {
    if id == 0 {
        return nil, errors.New("invalid updated ID")
    }
    res, err := r.db.Exec("UPDATE comments SET content = ?, post_id = ?, user_id = ?, date = ?, WHERE id = ?", updated.Content, updated.Post_ID, updated.User_ID, updated.Date, id)
    if err != nil {
        return nil, err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return nil, err
    }

    if rowsAffected == 0 {
        return nil, ErrUpdateFailed
    }

    return &updated, nil
}

func (r *SQLiteRepository) DeleteComment(id int64) error {
    res, err := r.db.Exec("DELETE FROM comments WHERE id = ?", id)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrDeleteFailed
    }

    return err
}

/*---------------------
     Like Queries
----------------------*/

func (r *SQLiteRepository) TableLike() error {
    query := `
    CREATE TABLE IF NOT EXISTS like(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        value INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        date DATE,
        FOREIGN KEY (user_id)
            REFERENCES users (id),
        FOREIGN KEY (post_id)
            REFERENCES posts (id)
    );
    `

    _, err := r.db.Exec(query)
    return err
}

func (r *SQLiteRepository) CreateLike(like models.Like) (*models.Like, error) {
    res, err := r.db.Exec("INSERT INTO like(value, post_id, user_id, date) values(?,?,?,?)", 
    like.Value, like.Post_ID, like.User_ID, time.Now())
    if err != nil {
        var sqliteErr sqlite3.Error
        if errors.As(err, &sqliteErr) {
            if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
                return nil, ErrDuplicate
            }
        }
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    like.ID = id

    return &like, nil
}

func (r *SQLiteRepository) UpdateLike(id int64, updated models.Like) (*models.Like, error) {
    if id == 0 {
        return nil, errors.New("invalid updated ID")
    }
    res, err := r.db.Exec("UPDATE posts SET value = ?, post_id = ?, user_id = ?, date = ?, WHERE id = ?", updated.Value, updated.Post_ID, updated.User_ID, updated.Date, id)
    if err != nil {
        return nil, err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return nil, err
    }

    if rowsAffected == 0 {
        return nil, ErrUpdateFailed
    }

    return &updated, nil
}

func (r *SQLiteRepository) DeleteLike(id int64) error {
    res, err := r.db.Exec("DELETE FROM like WHERE id = ?", id)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrDeleteFailed
    }

    return err
}