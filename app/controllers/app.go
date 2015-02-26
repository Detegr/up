package controllers

import (
	"github.com/revel/revel"
	"github.com/Detegr/up/db"
)

type App struct {
	*revel.Controller
}

func (c App) CurrentUser() *db.User {
	var user db.User
	if c.Session["User"] != "" {
		conn.Where("name = ?", c.Session["User"]).First(&user)
		return &user
	}
	return nil
}

func (c App) Index() revel.Result {
	var files []db.File
	user := c.CurrentUser()
	if user != nil {
		if err := conn.Model(&user).Order("created desc").Related(&files).Error; err != nil {
			c.Flash.Error("Sorry, could not fetch your uploaded files.")
		}
	}
	return c.Render(files)
}
