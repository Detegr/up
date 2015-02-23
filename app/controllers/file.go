package controllers

import (
	"github.com/revel/revel"
	"os"
	"path"
	"strings"
	"errors"
	"github.com/Detegr/up/db"
)

type File struct {
	App
}

func WriteFileToDisk(filename string, file []byte) error {
	localfile, err := os.Create(path.Join("uploads", filename))
	if err != nil {
		return err
	}
	_, err = localfile.Write(file)
	if err != nil {
		return err
	}
	return nil
}

func SaveFileToDb(c File, filename string) (string, error) {
	var existing db.File
	file := db.File { FileName: filename }
	if err := conn.Where(&file).First(&existing).Error; err == nil {
		ext := path.Ext(filename)
		newfile := strings.TrimSuffix(filename, ext) + "_" + ext
		if len(filename) > 64 {
			return "", errors.New("Too many files with same name") // I'm lazy
		}
		return SaveFileToDb(c, newfile)
	}
	if c.Session["User"] == "" {
		if err := conn.Create(&file).Error; err != nil {
			return "", err
		}
	} else {
		var user db.User
		err := conn.Where(db.User {Name: c.Session["User"]}).First(&user).Error;
		if err == nil {
			file.UserId = user.Id
			if err := conn.Model(&user).Association("Files").Append(&file).Error; err != nil {
				return "", err
			}
			return filename, nil
		}
		return "", err
	}
	return filename, nil
}

func (c File) Upload(file []byte) revel.Result {
	if c.Params.Files["file"] == nil {
		c.Flash.Error("No file was found, please try again")
		return c.Redirect(App.Index)
	}
	filename := c.Params.Files["file"][0].Filename
	filename, err := SaveFileToDb(c, filename)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(App.Index)
	}
	if err := WriteFileToDisk(filename, file); err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(App.Index)
	}
	c.Flash.Success("Success!")
	c.Flash.Out["FileName"] = filename
	return c.Redirect(App.Index)
}

func (c File) Serve(filename string) revel.Result {
	var dbfile db.File
	println(filename)
	var err error
	if err = conn.Where("file_name = ?", filename).First(&dbfile).Error; err == nil {
		file, err := os.Open(path.Join("uploads", dbfile.FileName))
		if err != nil {
			c.Flash.Error(err.Error())
			return c.Redirect(App.Index)
		}
		return c.RenderFile(file, "inline")
	}
	c.Flash.Error("File %s not found", filename)
	return c.Redirect(App.Index)
}
