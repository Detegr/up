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
	*revel.Controller
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

func SaveFileToDb(filename string) (string, error) {
	var existing db.File
	file := db.File { FileName: filename }
	if err := conn.Where(&file).First(&existing).Error; err == nil {
		ext := path.Ext(filename)
		newfile := strings.TrimSuffix(filename, ext) + "_" + ext
		if len(filename) > 64 {
			return "", errors.New("Too many files with same name") // I'm lazy
		}
		return SaveFileToDb(newfile)
	}
	if err := conn.Create(&file).Error; err != nil {
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
	filename, err := SaveFileToDb(filename)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(App.Index)
	}
	if err := WriteFileToDisk(filename, file); err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(App.Index)
	}
	c.Flash.Success("File %s uploaded!", filename)
	return c.Redirect(App.Index)
}
