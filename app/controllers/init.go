package controllers

import (
	"github.com/revel/revel"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/Detegr/up/db"
)

var conn *gorm.DB

func InitDB() {
	c, err := gorm.Open("postgres", "user=postgres dbname=up sslmode=disable")
	if err != nil {
		panic("Could not open database")
	}
	c.AutoMigrate(&db.User{}, &db.File{})
	c.Model(&db.User{}).AddUniqueIndex("idx_name", "name")
	conn = &c
}

func init() {
	revel.OnAppStart(InitDB)
	revel.TemplateFuncs["uploadedFile"] = func(flash map[string]string) string { return flash["FileName"] }
	revel.TemplateFuncs["filePresent"] = func(flash map[string]string) bool { return flash["FileName"] != "" }
	revel.TemplateFuncs["userLoggedIn"] = func(session map[string]string) bool { return session["User"] != "" }
}
