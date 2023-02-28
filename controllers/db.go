package forum

import (
    "database/sql"
    "errors"

    "github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/bcrypt"
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
     Users Query
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
        date DATE
    );
    `

    _, err := r.db.Exec(query)
    return err
}

func (r *SQLiteRepository) CreateUser(users models.Users) (*models.Users, error) {

    password, err := HashPassword(users.Password)

    res, err := r.db.Exec("INSERT INTO users(username, mail, checkedmail, password, categories, date) values(?,?,?,?,?,?)", 
    users.Username, users.Mail, users.CheckedMail, password, users.Categories, users.Date)
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

func HashPassword(password string) (string, error) {
    var passwordBytes = []byte(password)

    hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

    return string(hashedPasswordBytes), err
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
        if err := rows.Scan(&users.ID, &users.Username, &users.Mail, &users.CheckedMail, &users.Password, &users.Categories, &users.Date); err != nil {
            return nil, err
        }
        all = append(all, users)
    }
    return all, nil
}

func (r *SQLiteRepository) ConnectUser(username string, password string) (*models.Users, error) {
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return nil, err
    }

    row := r.db.QueryRow("SELECT * FROM users WHERE username = ?, password = ?", username, hashedPassword)

    var users models.Users
    if err := row.Scan(&users.ID, &users.Username, &users.Mail, &users.CheckedMail, &users.Password, &users.Categories, &users.Date); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }
    return &users, nil
}

func (r *SQLiteRepository) UpdateUser(id int64, updated models.Users) (*models.Users, error) {
    if id == 0 {
        return nil, errors.New("invalid updated ID")
    }
    res, err := r.db.Exec("UPDATE users SET username = ?, mail = ?, checkedmail = ?, password = ?, categories = ?, date = ?, WHERE id = ?", updated.Username, updated.Mail, updated.CheckedMail, updated.Password, updated.Categories, updated.Date, id)
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

/*---------------------
     Posts Query
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
    posts.Name, posts.Content, posts.Categories, posts.User_ID, posts.Date)
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
     Comments Query
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
    comments.Content, comments.Post_ID, comments.User_ID, comments.Date)
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
    res, err := r.db.Exec("UPDATE posts SET content = ?, post_id = ?, user_id = ?, date = ?, WHERE id = ?", updated.Content, updated.Post_ID, updated.User_ID, updated.Date, id)
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