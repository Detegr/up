package db

type User struct {
	Id int64
	Name string `sql:"not null; unique"`
	Password string `sql:"not null"`
	Files []File
}
