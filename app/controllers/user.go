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
	*models.User `json:"user"`
}

func (c UserController) Create() revel.Result {
	bodyUser, err := c.getBodyUser()
	if err != nil {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}

	user := models.NewUser(bodyUser.Username, bodyUser.Email, bodyUser.Password)

	user.Validate(c.Validation)

	usernameUnique := c.FindUserByUsername(bodyUser.Username) == nil
	emailUnique := c.FindUserByEmail(bodyUser.Email) == nil
	c.Validation.Required(usernameUnique).Key("username").Message(models.TAKEN_MSG)
	c.Validation.Required(emailUnique).Key("email").Message(models.TAKEN_MSG)
	if c.Validation.HasErrors() {
		errs := &errorJSON{}
		errs = errs.Build(c.Validation.ErrorMap())
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errs)
	}

	err = c.Txn.Insert(user)
	if err != nil {
		revel.ERROR.Println(err)
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJSON(http.StatusText(c.Response.Status))
	}

	res := &UserJSON{
		&models.User{
			Username: user.Username,
			Email:    user.Email,
			Token:    c.JWT.NewToken(user.ID, user.Username),
		},
	}
	return c.RenderJSON(res)
}
func (c UserController) Read() revel.Result {
	return c.RenderJSON(UserJSON{c.Args[currentUserKey].(*models.User)})
}
func (c UserController) Update() revel.Result {
	bodyUser, err := c.getBodyUser()
	if err != nil {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}

	user := c.Args[currentUserKey].(*models.User)
	user.Fill(bodyUser)
	user.Validate(c.Validation)
	c.checkAlreadyTaken(bodyUser, user)
	if c.Validation.HasErrors() {
		errs := &errorJSON{}
		errs = errs.Build(c.Validation.ErrorMap())
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errs)
	}

	_, err = c.Txn.Update(user)
	if err != nil {
		revel.ERROR.Fatal(err)
	}

	res := &UserJSON{
		&models.User{
			Username: user.Username,
			Email:    user.Email,
			Token:    c.JWT.NewToken(user.ID, user.Username),
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}

	return c.RenderJSON(res)
}

func (c UserController) getBodyUser() (*models.User, error) {
	body := UserJSON{}
	err := c.Params.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return body.User, nil
}

func (c *UserController) checkAlreadyTaken(bodyUser *models.User, user *models.User) {
	c.checkAlreadyTakenUsername(bodyUser, user)
	c.checkAlreadyTakenEmail(bodyUser, user)
}

func (c *UserController) checkAlreadyTakenEmail(bodyUser *models.User, user *models.User) {
	userByEmail := c.FindUserByEmail(bodyUser.Email)
	emailUnique := userByEmail == nil || userByEmail.ID == user.ID
	c.Validation.Required(emailUnique).Key("email").Message(models.TAKEN_MSG)
}

func (c *UserController) checkAlreadyTakenUsername(bodyUser *models.User, user *models.User) {
	userByUsername := c.FindUserByUsername(bodyUser.Username)
	usernameUnique := userByUsername == nil || userByUsername.ID == user.ID
	c.Validation.Required(usernameUnique).Key("username").Message(models.TAKEN_MSG)
}

func (c UserController) Login() revel.Result {
	body := UserJSON{}
	var err error

	err = c.Params.BindJSON(&body)

	if err != nil {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errorJSON{Errors: ValidationErrors{"BindJSON": {err.Error()}}})
	}

	user, errs := c.checkValidate(body.User)
	if errs != nil {
		c.Response.Status = http.StatusUnprocessableEntity
		return c.RenderJSON(errs)
	}
	res := &UserJSON{
		&models.User{
			Username: user.Username,
			Email:    user.Email,
			Token:    c.JWT.NewToken(user.ID, user.Username),
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}

	return c.RenderJSON(res)
}

func (c UserController) checkValidate(bodyUser *models.User) (user *models.User, errs *errorJSON) {
	errs = &errorJSON{}
	c.Validation.Required(bodyUser.Email).Key("email").Message(models.EMPTY_MSG)
	c.Validation.Required(bodyUser.Password).Key("password").Message(models.EMPTY_MSG)
	if c.Validation.HasErrors() {
		errs = errs.Build(c.Validation.ErrorMap())
		return nil, errs
	}

	user = c.FindUserByEmail(bodyUser.Email)
	c.Validation.Required(user != nil).Key("email").Message(models.INVALID_MSG)
	if c.Validation.HasErrors() {
		errs = errs.Build(c.Validation.ErrorMap())
		return nil, errs
	}
	match := user.MatchPassword(bodyUser.Password)
	c.Validation.Required(match).Key("password").Message(models.INVALID_MSG)
	if c.Validation.HasErrors() {
		errs = errs.Build(c.Validation.ErrorMap())
		return nil, errs
	}
	return user, nil
}
