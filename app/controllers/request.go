package controllers

import (
	"amplifier/app/db"
	"amplifier/app/entities"
	"amplifier/app/forms"
	"amplifier/app/models"
	"amplifier/app/webutils"

	"github.com/revel/revel"
)

type Requests struct {
	App
}

func (c Requests) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(App.Index)
	}
	return nil
}

func (c Requests) All() revel.Result {
	result := entities.Response{}
	ctx := c.Request.Context()

	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		result.Success = false
		result.Message = "Failed to parse page filters"
		return c.Render(result)
	}

	newRequest := &models.Request{}
	data, err := newRequest.All(ctx, db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get messages: %v", err)
		result.Success = false
		result.Message = "Could not get messages"
		return c.Render(result)
	}

	recordsCount, err := newRequest.Count(ctx, db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get messages count: %v", err)
		result.Success = false
		result.Message = "Could not get messages count"
		return c.Render(result)
	}

	revel.AppLog.Infof("data: %+v", data[0])

	result.Success = true
	result.Data = map[string]interface{}{
		"Requests":   data,
		"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
	}
	return c.Render(result)
}

func (c Requests) New() revel.Result {
	return c.Render()
}

func (c Requests) Save(request *forms.Request) revel.Result {

	v := c.Validation
	request.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(Requests.New)
	}

	newRequest := &models.Request{
		App:     request.App,
		Multi:   request.Multi,
		Number:  request.Count,
		Message: request.Message,
		Times:   request.Times,
	}
	err := newRequest.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("could not save credential: %v", err)
		c.Flash.Error("internal server error -  could not save credential")
		c.FlashParams()
		return c.Redirect(Requests.New)
	}

	c.Flash.Success("Request created for app - " + newRequest.App)
	return c.Redirect(Requests.All)
}
