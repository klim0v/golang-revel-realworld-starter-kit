package controllers

import (
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/revel/revel"
)

type CommentController struct {
	ApplicationController
}

type Comment struct {
	ID        int         `json:"id"`
	Body      string      `json:"body"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	Author    models.User `json:"author"`
}

type CommentJSON struct {
	Comment Comment `json:"comment"`
}

type CommentsJSON struct {
	Comments []Comment `json:"comments"`
}

func (c CommentController) Index() revel.Result {
	return c.Todo()
}

func (c CommentController) Create() revel.Result {
	return c.Todo()
}

func (c CommentController) Delete() revel.Result {
	return c.Todo()
}
