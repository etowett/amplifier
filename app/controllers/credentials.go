package controllers

import (
	"amplifier/app/db"
	"amplifier/app/entities"
	"amplifier/app/models"
	"amplifier/app/webutils"

	"github.com/revel/revel"
)

type Credentials struct {
	App
}

func (c Credentials) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(App.Index)
	}
	return nil
}

func (c Credentials) All() revel.Result {
	result := entities.Response{}
	ctx := c.Request.Context()

	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		result.Success = false
		result.Message = "Failed to parse page filters"
		return c.Render(result)
	}

	newCredential := &models.Credential{}
	data, err := newCredential.All(ctx, db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get messages: %v", err)
		result.Success = false
		result.Message = "Could not get messages"
		return c.Render(result)
	}

	recordsCount, err := newCredential.Count(ctx, db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get messages count: %v", err)
		result.Success = false
		result.Message = "Could not get messages count"
		return c.Render(result)
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"Credentials": data,
		"Pagination":  models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
	}
	return c.Render(result)
}
