package controllers

import "github.com/revel/revel"
import "net/http"
import "io/ioutil"
import "net/url"
import "encoding/json"
import "errors"
import "github.com/Detegr/up/db"
import "github.com/Detegr/up/secrets"
import "golang.org/x/crypto/bcrypt"

type User struct {
	App
}

func (c User) Register(username string, password string) revel.Result {
	return c.Render()
}

func VerifyNoRobot(recaptcha string) error {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{"secret": {secrets.Recaptcha}, "response": {recaptcha}})
	if err != nil {
		return errors.New("Could not verify that you aren't a robot. Please try again.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Could not verify that you aren't a robot. Please try again.")
	}
	type RecaptchaReply struct {
		Success bool `json:"success"`
	}
	var reply RecaptchaReply
	err = json.Unmarshal(body, &reply)
	if err != nil || reply.Success == false {
		return errors.New("Could not verify that you aren't a robot. Please try again.")
	}
	return nil
}

func CreateNewUser(username string, password string) error {
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("An error occurred when creating your account. Please try again.")
	}
	user := db.User{
		Name:     username,
		Password: string(pwhash),
		Files: make([]db.File, 0),
	}
	if err = conn.Save(&user).Error; err != nil {
		return errors.New("User already exists, please select another username.")
	}
	return nil
}

func (c User) ValidateRegistration(username string, password string) revel.Result {
	recaptcha := c.Params.Values["g-recaptcha-response"][0]
	err := VerifyNoRobot(recaptcha)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(User.Register)
	}
	err = CreateNewUser(username, password)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(User.Register)
	}
	c.Flash.Success("Account successfully created. Please login from below.")
	return c.Redirect(User.Login)
}

func (c User) Login() revel.Result {
	return c.Render()
}

func (c User) ValidateLogin(username string, password string) revel.Result {
	var user db.User
	if err := conn.Where(&db.User{Name: username}).First(&user).Error; err != nil {
		c.Flash.Error("Wrong username or password.")
		return c.Redirect(User.Login)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		c.Flash.Error("Wrong username or password.")
		return c.Redirect(User.Login)
	}
	c.Session["User"] = username
	return c.Redirect(App.Index)
}

func (c User) Logout() revel.Result {
	if c.Session["User"] == "" {
		return c.Redirect(App.Index)
	}
	c.Session["User"] = ""
	c.Flash.Success("Successfully logged out.")
	return c.Redirect(App.Index)
}
