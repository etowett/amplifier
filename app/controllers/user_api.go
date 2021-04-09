package controllers

import (
	"amplifier/app/db"
	"amplifier/app/entities"
	"amplifier/app/forms"
	"amplifier/app/models"
	"database/sql"
	"net/http"
	"strings"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type (
	UsersAPI struct {
		App
	}
)

func (c UsersAPI) Save() revel.Result {
	status := http.StatusCreated
	ctx := c.Request.Context()
	userForm := forms.User{}
	err := c.Params.BindJSON(&userForm)
	if err != nil {
		c.Log.Errorf("failed to bind to json user create api: %v", err)
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Invalid form provided",
			Status:  status,
			Success: false,
		})
	}

	v := c.Validation
	userForm.Validate(v)
	if v.HasErrors() {
		retErrors := make([]string, 0)
		for _, theErr := range v.Errors {
			retErrors = append(retErrors, theErr.Message)
		}
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: strings.Join(retErrors, ", "),
			Status:  status,
			Success: false,
		})
	}

	c.Log.Infof("userForm api: =[%+v]", userForm)

	theUser := &models.User{}
	userByMail, err := theUser.ByEmail(ctx, db.DB(), userForm.Email)
	if err != nil && err != sql.ErrNoRows {
		c.Log.Errorf("error getting user by email: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Internal error occured!",
			Status:  status,
			Success: false,
		})
	}

	if userByMail.ID != 0 {
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "User with that email already exist!",
			Status:  status,
			Success: false,
		})
	}

	userByUsername, err := theUser.ByUsername(ctx, db.DB(), userForm.Username)
	if err != nil && err != sql.ErrNoRows {
		c.Log.Errorf("error getting user by username: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Internal error occured!",
			Status:  status,
			Success: false,
		})
	}

	if userByUsername.ID != 0 {
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "User with that Username already exist!",
			Status:  status,
			Success: false,
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Errorf("error generating password hash: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Internal error occured!",
			Status:  status,
			Success: false,
		})
	}

	newUser := &models.User{
		Username:     userForm.Username,
		FirstName:    userForm.FirstName,
		LastName:     userForm.LastName,
		Email:        userForm.Email,
		PasswordHash: string(passwordHash[:]),
	}

	err = newUser.Save(ctx, db.DB())
	if err != nil {
		c.Log.Errorf("error insert user: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Encountered an error saving request.",
			Status:  status,
			Success: false,
		})
	}

	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Data:    newUser,
		Status:  status,
		Success: true,
	})
}
