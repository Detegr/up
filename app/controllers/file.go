package controllers

import (
	"github.com/revel/revel"
	"os"
	"path"
	"strings"
	"errors"
	"github.com/Detegr/up/db"
	"github.com/Detegr/up/app/routes"
	"mime"
	"net/http"
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

func SaveFileToDb(c File, contenttype string, filename string) (string, error) {
	var existing db.File
	if err := conn.Where("file_name = ?", filename).First(&existing).Error; err == nil {
		ext := path.Ext(filename)
		newfile := strings.TrimSuffix(filename, ext) + "_" + ext
		if len(filename) > 64 {
			return "", errors.New("Too many files with same name") // I'm lazy
		}
		return SaveFileToDb(c, contenttype, newfile)
	}
	file := db.File {
		FileName: filename,
		ContentType: contenttype,
	}
	user := c.CurrentUser()
	if user == nil {
		if err := conn.Create(&file).Error; err != nil {
			return "", err
		}
	} else {
		if err := conn.Model(&user).Association("Files").Append(&file).Error; err != nil {
			return "", err
		}
		return filename, nil
	}
	return filename, nil
}

func (c File) Upload(file []byte) revel.Result {
	if c.Params.Files["file"] == nil {
		c.Flash.Error("No file was found, please try again")
		return c.Redirect(App.Index)
	}
	filename := c.Params.Files["file"][0].Filename
	contenttype := mime.TypeByExtension(path.Ext(filename))
	if contenttype == "" {
		// Try to figure out the content type from the data
		contenttype = http.DetectContentType(file)
	}
	filename, err := SaveFileToDb(c, contenttype, filename)
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
		// For backwards compatibility, fetch and update the content type before serving the file
		if dbfile.ContentType == "" {
			contenttype := mime.TypeByExtension(path.Ext(filename))
			if contenttype == "" {
				buf := make([]byte, 512) // DetectContentType needs up to 512 bytes of buffer
				_, err := file.Read(buf)
				if err != nil {
					c.Flash.Error("Error retrieving the file, please try again.")
					return c.Redirect(App.Index)
				}
				contenttype = http.DetectContentType(buf)
			}
			dbfile.ContentType = contenttype
			// Gorm does not like my string-type primary key :( Gotta go raw SQL
			conn.Exec("UPDATE files SET content_type=? WHERE file_name=?", contenttype, filename)
		}
		c.Response.ContentType = dbfile.ContentType
		return c.RenderFile(file, revel.Inline)
	}
	c.Flash.Error("File %s not found", filename)
	return c.Redirect(App.Index)
}

func (c File) Delete(filename string) revel.Result {
	user := c.CurrentUser()
	if user == nil {
		return c.Redirect(routes.File.Serve(filename))
	}
	if err := conn.Where("file_name = ? AND user_id = ?", filename, user.Id).Delete(db.File{}).Error; err == nil {
		c.Flash.Success("File %s deleted.", filename)
		return c.Redirect(App.Index)
	}
	c.Flash.Error("Could not delete the file. Please try again.")
	return c.Redirect(App.Index)
}
