package forum

type Posts struct {
	ID         int64
	Name       string
	Content    string
	Categories string
	User_ID    int64
	Date       string
	Nb_Like    int
	Is_Like    bool
}
