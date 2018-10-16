package controllers

import (
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/revel/revel"
	"net/http"
)

type ArticleController struct {
	ApplicationController
}

type ArticleJSON struct {
	*models.Article `json:"article"`
}

type ArticlesJSON struct {
	Articles      []*models.Article `json:"articles"`
	ArticlesCount int               `json:"articlesCount"`
}

func (c ArticleController) Index() revel.Result {
	return c.Todo()
}
func (c ArticleController) Create() revel.Result {
	article, err := c.getBodyArticle()
	if err != nil {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}
	user := c.Args[currentUserKey].(*models.User)
	newArticle := models.NewArticle(article.Title, article.Description, article.Body, article.TagList, user)
	err = c.Txn.Insert(newArticle)
	if err != nil {
		revel.ERROR.Println(err)
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJSON(http.StatusText(c.Response.Status))
	}

	res := &ArticleJSON{
		newArticle,
	}
	c.Response.Status = http.StatusCreated
	return c.RenderJSON(res)
}
func (c ArticleController) Read() revel.Result {
	return c.Todo()
}
func (c ArticleController) Update() revel.Result {
	return c.Todo()
}
func (c ArticleController) Delete() revel.Result {
	return c.Todo()
}

func (c ArticleController) getBodyArticle() (*models.Article, error) {
	body := ArticleJSON{}
	err := c.Params.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return body.Article, nil
}
