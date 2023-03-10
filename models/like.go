package forum

type Like struct {
	ID      int64
	Value   int
	User_ID int64
	Post_ID int64
	Date    string
}
