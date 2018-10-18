package controllers

import (
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/revel/revel"
	"net/http"
	"strconv"
	"strings"
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

func (c ArticleController) Index(tag, favorited, author string, offset, limit uint64) revel.Result {
	var articles []*models.Article
	var user *models.User
	builder := c.Db.SqlStatementBuilder.Select("*").
		From("Article").
		Offset(offset).
		Limit(limit)
	if author != "" {
		user = c.FindUserByUsername(author)
		if user != nil {
			builder.Where("UserID=?", user.ID)
		}
	}
	if _, err := c.Txn.Select(&articles, builder); err != nil {
		c.Log.Fatal("Unexpected error loading articles", "error", err)
	}
	if user != nil {
		for _, article := range articles {
			article.User = user
		}
	} else {
		var users []*models.User
		var userIds []string
		for _, article := range articles {
			userIds = append(userIds, strconv.Itoa(article.UserID))
		}
		selectBuilder := c.Db.SqlStatementBuilder.
			Select("*").
			From("User").
			Where(
				"ID in (?)",
				strings.Join(userIds, ","),
			)
		c.Txn.Select(&users, selectBuilder)
		for _, article := range articles {
			for _, user := range users {
				if user.ID == article.UserID {
					article.User = user
					break
				}
			}
		}
	}
	return c.RenderJSON(&ArticlesJSON{articles, len(articles)})
}

func (c ArticleController) Create() revel.Result {
	article, err := c.getBodyArticle()
	if err != nil {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}
	article.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(BuildErrors(c.Validation.ErrorMap()))
	}
	articleObj := models.NewArticle(
		article.Title,
		article.Description,
		article.Body,
		article.TagList,
		c.Args[currentUserKey].(*models.User),
	)
	if err = c.Txn.Insert(articleObj); err != nil {
		c.Log.Fatal("Unexpected error insert article", "error", err)
	}
	c.Response.Status = http.StatusCreated
	return c.RenderJSON(&ArticleJSON{articleObj})
}

func (c ArticleController) Read() revel.Result {
	article := c.Args[fetchedArticleKey].(*models.Article)
	return c.RenderJSON(&ArticleJSON{article})
}

func (c ArticleController) Update() revel.Result {
	article, err := c.getBodyArticle()
	if err != nil {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}

	articleObj := c.Args[fetchedArticleKey].(*models.Article)
	articleObj.Fill(article)
	if _, err = c.Txn.Update(articleObj); err != nil {
		c.Log.Fatal("Unexpected error update article", "error", err)
	}
	return c.RenderJSON(&ArticleJSON{articleObj})
}

func (c ArticleController) Delete() revel.Result {
	article := c.Args[fetchedArticleKey].(*models.Article)
	if _, err := c.Txn.Delete(article); err != nil {
		c.Log.Fatal("Unexpected error delete article", "error", err)
	}
	c.Response.Status = http.StatusOK
	return c.RenderJSON(http.StatusText(c.Response.Status))
}

func (c ArticleController) getBodyArticle() (*models.Article, error) {
	body := ArticleJSON{}
	err := c.Params.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return body.Article, nil
}
