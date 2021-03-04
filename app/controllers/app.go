package controllers

import (
	"time"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Health() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"success":     true,
		"message":     "Ok",
		"server_time": time.Now().String(),
	})
}
