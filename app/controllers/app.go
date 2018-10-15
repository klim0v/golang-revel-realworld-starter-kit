package controllers

import (
	"database/sql"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/lib/auth"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/revel/modules/orm/gorp/app/controllers"
	"github.com/revel/revel"
)

const (
	currentUserKey = "current_user"
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
func (c *ApplicationController) AddUser() revel.Result {
	c.Args[currentUserKey] = c.currentUser()
	return nil
}

func (c *ApplicationController) currentUser() *models.User {
	claims, err := c.JWT.CheckRequest(c.Request)
	if err != nil {
		return nil
	}
	obj, err := c.Db.Get(models.User{}, claims.UserID)
	if obj == nil {
		return nil
	}
	user := obj.(*models.User)
	user.Token = c.JWT.NewToken(user.ID, user.Username)
	return user
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
			revel.ERROR.Fatal(err)
		}
		return nil
	}
	return user
}
