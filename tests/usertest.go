package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/controllers"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/routes"
	"net/http"
)

type UserControllerTest struct {
	AppTest
}

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistrationBody struct {
	User UserRegister `json:"user"`
}

type UserLoginBody struct {
	User UserLogin `json:"user"`
}

type testRegistration struct {
	errorKey string
	message  string
	body     UserRegistrationBody
}

type testLogin struct {
	errorKey string
	message  string
	body     UserLoginBody
}

func (t *UserControllerTest) TestLoginSuccessFully() {
	bodyUser := UserLoginBody{
		UserLogin{
			Email:    demoEmail,
			Password: demoPassword,
		},
	}

	jsonBody, _ := json.Marshal(bodyUser)

	t.MakePostRequest(routes.UserController.Login(), bytes.NewBuffer(jsonBody), nil)
	t.AssertOk()

	var UserJSON = controllers.UserJSON{}
	json.Unmarshal(t.ResponseBody, &UserJSON)

	t.AssertEqual(JWT.NewToken(demoUsername), UserJSON.User.Token)
	t.AssertEqual(demoUsername, UserJSON.User.Username)
	t.AssertEqual(bodyUser.User.Email, UserJSON.User.Email)
}

func (t *UserControllerTest) TestLoginFail() {
	tests := []testLogin{
		{
			errorKey: "email",
			message:  models.CORRECT_MSG,
			body: UserLoginBody{
				UserLogin{
					Email:    "",
					Password: demoRegPassword,
				},
			},
		},
		{
			errorKey: "password",
			message:  models.CORRECT_MSG,
			body: UserLoginBody{
				UserLogin{
					Email:    demoRegEmail,
					Password: "",
				},
			},
		},
	}

	for _, test := range tests {
		jsonBody, _ := json.Marshal(test.body)

		t.MakePostRequest(routes.UserController.Login(), bytes.NewBuffer(jsonBody), nil)
		t.AssertStatus(422)

		var ErrorJSON = ErrorJSON{}
		json.Unmarshal(t.ResponseBody, &ErrorJSON)

		msg, ok := ErrorJSON.Errors[test.errorKey]
		t.Assert(ok)
		t.AssertEqual(test.message, msg[0])
	}

	jsonBody, _ := json.Marshal(UserLoginBody{})

	t.MakePostRequest(routes.UserController.Login(), bytes.NewBuffer(jsonBody), nil)
	t.AssertStatus(422)

	var ErrorJSON = ErrorJSON{}
	json.Unmarshal(t.ResponseBody, &ErrorJSON)

	var errorKeys = []string{"email", "password"}
	for _, errorKey := range errorKeys {
		msg, ok := ErrorJSON.Errors[errorKey]
		t.Assert(ok)
		t.AssertEqual(models.CORRECT_MSG, msg[0])
	}
}

func (t *UserControllerTest) TestRegistrationSuccessFully() {
	bodyUser := UserRegistrationBody{
		UserRegister{
			Username: demoRegUsername,
			Email:    demoRegEmail,
			Password: demoRegPassword,
		},
	}

	jsonBody, _ := json.Marshal(bodyUser)

	t.MakePostRequest(routes.UserController.Create(), bytes.NewBuffer(jsonBody), nil)
	t.AssertOk()

	var UserJSON = controllers.UserJSON{}
	json.Unmarshal(t.ResponseBody, &UserJSON)

	t.AssertEqual(JWT.NewToken(demoRegUsername), UserJSON.User.Token)
	t.AssertEqual(bodyUser.User.Username, UserJSON.User.Username)
	t.AssertEqual(bodyUser.User.Email, UserJSON.User.Email)
}

func (t *UserControllerTest) TestRegistrationFail() {
	tests := []testRegistration{
		{
			errorKey: "username",
			message:  models.CORRECT_MSG,
			body: UserRegistrationBody{
				UserRegister{
					Username: "",
					Email:    demoRegEmail,
					Password: demoRegPassword,
				},
			},
		},
		{
			errorKey: "email",
			message:  models.CORRECT_MSG,
			body: UserRegistrationBody{
				UserRegister{
					Username: demoRegUsername,
					Email:    "",
					Password: demoRegPassword,
				},
			},
		},
		{
			errorKey: "password",
			message:  models.CORRECT_MSG,
			body: UserRegistrationBody{
				UserRegister{
					Username: demoRegUsername,
					Email:    demoRegEmail,
					Password: "",
				},
			},
		},
	}

	for _, test := range tests {
		jsonBody, _ := json.Marshal(test.body)

		t.MakePostRequest(routes.UserController.Create(), bytes.NewBuffer(jsonBody), nil)
		t.AssertStatus(422)

		var ErrorJSON = ErrorJSON{}
		json.Unmarshal(t.ResponseBody, &ErrorJSON)

		msg, ok := ErrorJSON.Errors[test.errorKey]
		t.Assert(ok)
		t.AssertEqual(test.message, msg[0])
	}

	jsonBody, _ := json.Marshal(UserRegistrationBody{})

	t.MakePostRequest(routes.UserController.Create(), bytes.NewBuffer(jsonBody), nil)
	t.AssertStatus(422)

	var ErrorJSON = ErrorJSON{}
	json.Unmarshal(t.ResponseBody, &ErrorJSON)

	var errorKeys = []string{"username", "email", "password"}
	for _, errorKey := range errorKeys {
		msg, ok := ErrorJSON.Errors[errorKey]
		t.Assert(ok)
		t.AssertEqual(models.CORRECT_MSG, msg[0])
	}
}

func (t *UserControllerTest) TestGetCurrentUserUnauthorized() {
	t.Get(routes.UserController.Read())
	t.AssertStatus(401)
}

func (t *UserControllerTest) TestGetCurrentUserSuccess() {
	request := t.GetCustom(t.BaseUrl() + routes.UserController.Read())

	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Token %v", JWT.NewToken(demoUsername))},
	}
	request.Send()
	t.AssertOk()
}

func (t *UserControllerTest) TestGetCurrentUserNotFound() {
	request := t.GetCustom(t.BaseUrl() + routes.UserController.Read())

	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Token %v", JWT.NewToken("not-found-user"))},
	}
	request.Send()
	t.AssertStatus(401)
}
