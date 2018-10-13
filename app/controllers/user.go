package controllers

import (
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/revel/revel"
	"net/http"
)

type UserController struct {
	ApplicationController
}

type UserJSON struct {
	models.User `json:"user"`
}

func (c UserController) Create() revel.Result {
	body := UserJSON{}

	var err error

	err = c.Params.BindJSON(&body)

	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}

	bodyUser := body.User
	revel.TRACE.Println("User Binding:", bodyUser)

	user := models.NewUser(bodyUser.Username, bodyUser.Email, bodyUser.Password)
	revel.TRACE.Println("User Entity:", user)
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		errs := &errorJSON{}
		errs = errs.Build(c.Validation.ErrorMap())
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errs)
	}

	err = c.Txn.Insert(user)
	if err != nil {
		panic(err)
	}

	res := &UserJSON{
		models.User{
			Username: user.Username,
			Email:    user.Email,
			Token:    c.JWT.NewToken(user.Username),
		},
	}
	return c.RenderJSON(res)
}
func (c UserController) Read() revel.Result {
	return c.RenderJSON(UserJSON{*c.Args[currentUserKey].(*models.User)})
}
func (c UserController) Update() revel.Result {
	body := UserJSON{}
	var err error

	err = c.Params.BindJSON(&body)

	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}

	bodyUser := body.User
	revel.TRACE.Println("User Binding:", bodyUser)

	user := c.Args[currentUserKey].(*models.User)
	user.Fill(bodyUser)
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		errs := &errorJSON{}
		errs = errs.Build(c.Validation.ErrorMap())
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errs)
	}

	_, err = c.Txn.Update(user)
	if err != nil {
		panic(err)
	}

	res := &UserJSON{
		models.User{
			Username: user.Username,
			Email:    user.Email,
			Token:    c.JWT.NewToken(user.Username),
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}

	return c.RenderJSON(res)
}
func (c UserController) Login() revel.Result {
	body := UserJSON{}
	var err error

	err = c.Params.BindJSON(&body)

	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}

	bodyUser := body.User
	revel.TRACE.Println("Body:", bodyUser)

	user := c.FindUserByEmail(bodyUser.Email)
	if user == nil {
		c.Response.Status = http.StatusNotFound
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"User": {http.StatusText(http.StatusNotFound)}}})
	}
	revel.TRACE.Println(user)
	ok := user.MatchPassword(bodyUser.Password)
	if !ok {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"User": {http.StatusText(http.StatusUnprocessableEntity)}}})
	}
	res := &UserJSON{
		models.User{
			Username: user.Username,
			Email:    user.Email,
			Token:    c.JWT.NewToken(user.Username),
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}

	return c.RenderJSON(res)
}
