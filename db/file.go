package db

import "time"

type File struct {
	FileName string    `sql:"type:varchar(64) PRIMARY KEY NOT NULL"`
	Created  time.Time `sql:"DEFAULT:current_timestamp"`
	Accessed time.Time `sql:"DEFAULT:current_timestamp"`
}