package controllers

import (
	"amplifier/app/db"
	"amplifier/app/models"
	"time"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	loggedInUser := c.connected()
	if loggedInUser != nil {
		return c.Redirect(App.Dash)
	}
	return c.Render()
}

func (c App) Dash() revel.Result {
	loggedInUser := c.connected()
	if loggedInUser == nil {
		return c.Redirect(App.Index)
	}
	return c.Render()
}

func (c App) Health() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"success":     true,
		"message":     "Ok",
		"server_time": time.Now().String(),
	})
}

func (c App) getUserFromUsername(username string) *models.User {
	user := &models.User{}
	c.Session.GetInto("user", user, false)
	if user.Username == username {
		return user
	}

	newUser := &models.User{}
	foundUser, err := newUser.ByUsername(c.Request.Context(), db.DB(), username)
	if err != nil {
		c.Log.Errorf("could not get user by username: %v", err)
		return nil
	}

	c.Session["user"] = foundUser
	return foundUser
}

func (c App) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["username"]; ok {
		return c.getUserFromUsername(username.(string))
	}
	return nil
}

func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}
