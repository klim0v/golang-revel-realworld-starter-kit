package controllers

import (
	"github.com/revel/revel"
	"net/http"
)

var requireAuth = map[string][]string{
	"UserController":    {"GET", "PUT"},
	"ArticleController": {"POST", "PUT", "DELETE"},
}

func authorize(c *revel.Controller) revel.Result {
	if methods, ok := requireAuth[c.Name]; ok {
		for _, v := range methods {
			if v == c.Request.Method {
				if c.Args[currentUserKey] == nil {
					c.Response.Status = http.StatusUnauthorized
					return c.Render(http.StatusText(c.Response.Status))
				}
			}
		}
	}
	return nil
}

func init() {
	revel.InterceptMethod((*ApplicationController).Init, revel.BEFORE)
	revel.InterceptMethod((*ApplicationController).AddUser, revel.BEFORE)
	revel.InterceptMethod((*ApplicationController).ExtractArticle, revel.BEFORE)
	revel.InterceptFunc(authorize, revel.BEFORE, &UserController{})
	revel.InterceptFunc(authorize, revel.BEFORE, &ArticleController{})
}
