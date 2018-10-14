package controllers

import (
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/revel/revel"
	"net/http"
	"reflect"
)

var requireAuth = map[string][]string{
	"UserController": {"GET", "PUT"},
}

func authorize(c *revel.Controller) revel.Result {
	if methods, ok := requireAuth[c.Name]; ok {
		for _, v := range methods {
			if v == c.Request.Method {
				if reflect.DeepEqual(c.Args[currentUserKey].(*models.User), &models.User{}) {
					c.Response.Status = http.StatusUnauthorized
					return c.RenderJSON(http.StatusText(c.Response.Status))
				}
			}
		}
	}
	return nil
}

func init() {
	revel.InterceptMethod((*ApplicationController).Init, revel.BEFORE)
	revel.InterceptMethod((*ApplicationController).AddUser, revel.BEFORE)
	revel.InterceptFunc(authorize, revel.BEFORE, &UserController{})
}
