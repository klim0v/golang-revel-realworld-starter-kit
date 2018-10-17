package controllers

import (
	"database/sql"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/services/auth"
	"github.com/revel/modules/orm/gorp/app/controllers"
	"github.com/revel/revel"
	"net/http"
)

const (
	currentUserKey    = "current_user"
	fetchedArticleKey = "article"
)

type ApplicationController struct {
	gorpController.Controller
	JWT auth.Tokener
}
type ValidationErrors map[string][]string

type errorJSON struct {
	Errors ValidationErrors `json:"errors"`
}

func (err *errorJSON) Build(errMap map[string]*revel.ValidationError) *errorJSON {
	err.Errors = ValidationErrors{}
	for _, validationError := range errMap {
		err.Errors[validationError.Key] = []string{validationError.Message}
	}
	return err
}

func (c *ApplicationController) Init() revel.Result {
	c.JWT = auth.NewJWT()
	return nil
}

func (c *ApplicationController) ExtractArticle() revel.Result {
	if slug := c.Params.Route.Get("slug"); slug != "" {
		article := c.FindArticleBySlug(slug)
		if article == nil {
			c.Response.Status = http.StatusNotFound
			return c.Render(http.StatusText(c.Response.Status))
		}
		c.Args[fetchedArticleKey] = article
	}
	return nil
}

func (c *ApplicationController) AddUser() revel.Result {
	claims, err := c.JWT.CheckRequest(c.Request)
	if err == nil {
		obj, _ := c.Db.Get(models.User{}, claims.UserID)
		if obj != nil {
			c.Args[currentUserKey] = obj.(*models.User)
		}
	}
	return nil
}

func (c ApplicationController) FindUserByUsername(username string) *models.User {
	return c.findUserByCondition("Username=?", username)
}

func (c ApplicationController) FindUserByEmail(email string) *models.User {
	return c.findUserByCondition("Email=?", email)
}

func (c ApplicationController) findUserByCondition(pred interface{}, args ...interface{}) *models.User {
	user := &models.User{}
	err := c.Txn.SelectOne(user, c.Db.SqlStatementBuilder.Select("*").From("User").Where(pred, args...))
	if err != nil {
		if err != sql.ErrNoRows {
			c.Log.Fatal("Unexpected error get article", "error", err)
		}
		return nil
	}
	return user
}

func (c ApplicationController) FindArticleBySlug(slug string) *models.Article {
	article := &models.Article{}
	err := c.Txn.SelectOne(article, c.Db.SqlStatementBuilder.Select("*").From("Article").Where("Slug=?", slug))
	if err != nil {
		if err != sql.ErrNoRows {
			c.Log.Fatal("Unexpected error get article", "error", err)
		}
		return nil
	}
	return article
}
