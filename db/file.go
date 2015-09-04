package db

import "time"

type File struct {
	FileName string    `sql:"type:varchar(1024) PRIMARY KEY NOT NULL"`
	UserId int64
	ContentType string
	Created  time.Time `sql:"DEFAULT:current_timestamp"`
	Accessed time.Time `sql:"DEFAULT:current_timestamp"`
}
