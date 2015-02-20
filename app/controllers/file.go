package controllers

import "github.com/revel/revel"
import "os"
import "path"

type File struct {
	*revel.Controller
}

func (c File) Upload(file []byte) revel.Result {
	filename := c.Params.Files["file"][0].Filename
	localfile, err := os.Create(path.Join("uploads", filename))
	if err != nil {
		panic("Could not create file")
	}
	_, err = localfile.Write(file)
	if err != nil {
		panic("Could not write to file")
	}
	return c.Redirect(App.Index)
}
