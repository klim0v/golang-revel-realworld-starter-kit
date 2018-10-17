package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app"
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

type UserUpdate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Password string `json:"password"`
}

type UserRegistrationBody struct {
	User UserRegister `json:"user"`
}

type UserUpdateBody struct {
	User UserUpdate `json:"user"`
}

type UserLoginBody struct {
	User UserLogin `json:"user"`
}

type testRegistration struct {
	errorKey string
	message  string
	body     UserRegistrationBody
}

type testUpdate struct {
	errorKey string
	message  string
	body     UserUpdateBody
}

type testLogin struct {
	errorKey string
	message  string
	body     UserLoginBody
}

func (t *UserControllerTest) TestLoginSuccessFully() {
	bodyUser := UserLoginBody{
		UserLogin{
			Email:    users[0].Email,
			Password: users[0].Password,
		},
	}

	jsonBody, _ := json.Marshal(bodyUser)

	t.MakePostRequest(routes.UserController.Login(), bytes.NewBuffer(jsonBody), nil)
	t.AssertOk()

	var UserJSON = controllers.UserJSON{}
	json.Unmarshal(t.ResponseBody, &UserJSON)

	claims, err := JWT.GetClaims(UserJSON.User.Token)
	t.Assert(err == nil)
	t.AssertEqual(users[0].Username, claims.Username)
	t.AssertEqual(users[0].Username, UserJSON.User.Username)
	t.AssertEqual(bodyUser.User.Email, UserJSON.User.Email)
}

func (t *UserControllerTest) TestLoginFail() {
	tests := []testLogin{
		{
			errorKey: "email",
			message:  models.EMPTY_MSG,
			body: UserLoginBody{
				UserLogin{
					Email:    "",
					Password: demoRegPassword,
				},
			},
		},
		{
			errorKey: "password",
			message:  models.EMPTY_MSG,
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
		t.AssertEqual(models.EMPTY_MSG, msg[0])
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
	t.AssertStatus(http.StatusCreated)

	var UserJSON = controllers.UserJSON{}
	json.Unmarshal(t.ResponseBody, &UserJSON)
	userId, err := app.Dbm.SelectInt("select ID from User where Username=? and Email=?", demoRegUsername, demoRegEmail)
	t.Assert(err == nil)
	t.AssertEqual(JWT.NewToken(int(userId), demoRegUsername), UserJSON.User.Token)
	t.AssertEqual(bodyUser.User.Username, UserJSON.User.Username)
	t.AssertEqual(bodyUser.User.Email, UserJSON.User.Email)
}

func (t *UserControllerTest) TestRegistrationFail() {
	tests := []testRegistration{
		{
			errorKey: "username",
			message:  models.EMPTY_MSG,
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
			message:  models.EMPTY_MSG,
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
			message:  models.EMPTY_MSG,
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
		t.AssertStatus(http.StatusUnprocessableEntity)

		var ErrorJSON = ErrorJSON{}
		json.Unmarshal(t.ResponseBody, &ErrorJSON)

		msg, ok := ErrorJSON.Errors[test.errorKey]
		t.Assert(ok)
		t.AssertEqual(test.message, msg[0])
	}

	jsonBody, _ := json.Marshal(UserRegistrationBody{})

	t.MakePostRequest(routes.UserController.Create(), bytes.NewBuffer(jsonBody), nil)
	t.AssertStatus(http.StatusUnprocessableEntity)

	var ErrorJSON = ErrorJSON{}
	json.Unmarshal(t.ResponseBody, &ErrorJSON)

	var errorKeys = []string{"username", "email", "password"}
	for _, errorKey := range errorKeys {
		msg, ok := ErrorJSON.Errors[errorKey]
		t.Assert(ok)
		t.AssertEqual(models.EMPTY_MSG, msg[0])
	}
}

func (t *UserControllerTest) TestGetCurrentUserUnauthorized() {
	t.Get(routes.UserController.Read())
	t.AssertStatus(http.StatusUnauthorized)
}

func (t *UserControllerTest) TestGetCurrentUserSuccess() {
	request := t.GetCustom(t.BaseUrl() + routes.UserController.Read())
	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Token %v", JWT.NewToken(users[0].ID, users[0].Username))},
	}
	request.Send()
	t.AssertOk()
}

func (t *UserControllerTest) TestGetCurrentUserNotFound() {
	request := t.GetCustom(t.BaseUrl() + routes.UserController.Read())

	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Token %v", JWT.NewToken(users[0].ID+999, ""))},
	}
	request.Send()
	t.AssertStatus(http.StatusUnauthorized)
}

func (t *UserControllerTest) TestGetCurrentUserInvalidToken() {
	request := t.GetCustom(t.BaseUrl() + routes.UserController.Read())

	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Token %v", "invalid-token")},
	}
	request.Send()
	t.AssertStatus(http.StatusUnauthorized)
}

func (t *UserControllerTest) TestUpdateUserFail() {
	testUpdate := []testUpdate{
		{
			errorKey: "email",
			message:  models.EMPTY_MSG,
			body: UserUpdateBody{
				UserUpdate{
					Username: demoUsername,
					Password: demoRegPassword,
				},
			},
		},
		{
			errorKey: "username",
			message:  models.EMPTY_MSG,
			body: UserUpdateBody{
				UserUpdate{
					Email:    demoRegEmail,
					Password: demoRegPassword,
				},
			},
		},
		{
			errorKey: "email",
			message:  models.TAKEN_MSG,
			body: UserUpdateBody{
				UserUpdate{
					Username: demoUsername,
					Email:    users[1].Email,
					Password: demoRegPassword,
				},
			},
		},
		{
			errorKey: "username",
			message:  models.TAKEN_MSG,
			body: UserUpdateBody{
				UserUpdate{
					Username: users[1].Username,
					Email:    demoEmail,
					Password: demoRegPassword,
				},
			},
		},
	}

	header := fmt.Sprintf("Token %v", JWT.NewToken(users[0].ID, users[0].Username))

	for _, test := range testUpdate {
		jsonBody, _ := json.Marshal(test.body)

		t.MakePutRequest(routes.UserController.Update(), bytes.NewBuffer(jsonBody), header)
		var ErrorJSON = ErrorJSON{}
		json.Unmarshal(t.ResponseBody, &ErrorJSON)
		msg, ok := ErrorJSON.Errors[test.errorKey]

		t.Assert(ok)
		t.AssertEqual(test.message, msg[0])
	}

	jsonBody, _ := json.Marshal(UserUpdateBody{})

	t.MakePutRequest(routes.UserController.Update(), bytes.NewBuffer(jsonBody), header)
	t.AssertStatus(http.StatusUnprocessableEntity)

	var ErrorJSON = ErrorJSON{}
	json.Unmarshal(t.ResponseBody, &ErrorJSON)

	var errorKeys = []string{"username", "email"}
	for _, errorKey := range errorKeys {
		msg, ok := ErrorJSON.Errors[errorKey]
		t.Assert(ok)
		t.AssertEqual(models.EMPTY_MSG, msg[0])
	}
}

func (t *UserControllerTest) TestUpdateSuccessFully() {
	bodyUser := UserUpdateBody{
		UserUpdate{
			Username: demoRegUsername,
			Email:    demoRegEmail,
			Image:    "newImage",
			Bio:      "newBio",
		},
	}
	header := fmt.Sprintf("Token %v", JWT.NewToken(users[0].ID, users[0].Username))

	jsonBody, err := json.Marshal(bodyUser)
	t.Assert(err == nil)
	t.MakePutRequest(routes.UserController.Update(), bytes.NewBuffer(jsonBody), header)
	t.AssertOk()

	var UserJSON = controllers.UserJSON{}
	json.Unmarshal(t.ResponseBody, &UserJSON)
	userId, err := app.Dbm.SelectInt("select ID from User where Username=? and Email=?", demoRegUsername, demoRegEmail)
	t.Assert(err == nil)
	t.AssertEqual(JWT.NewToken(int(userId), demoRegUsername), UserJSON.User.Token)
	t.AssertEqual(bodyUser.User.Username, UserJSON.User.Username)
	t.AssertEqual(bodyUser.User.Email, UserJSON.User.Email)
	t.AssertEqual(bodyUser.User.Bio, UserJSON.User.Bio)
	t.AssertEqual(bodyUser.User.Image, UserJSON.User.Image)
}

func (t *UserControllerTest) TestUpdatePasswordSuccessFully() {
	bodyUser := UserUpdateBody{
		UserUpdate{
			Username: users[0].Username,
			Email:    users[0].Email,
			Password: demoRegPassword,
		},
	}
	header := fmt.Sprintf("Token %v", JWT.NewToken(users[0].ID, users[0].Username))

	jsonBody, err := json.Marshal(bodyUser)
	t.Assert(err == nil)
	t.MakePutRequest(routes.UserController.Update(), bytes.NewBuffer(jsonBody), header)
	t.AssertOk()

	var UserJSON = controllers.UserJSON{}
	json.Unmarshal(t.ResponseBody, &UserJSON)
	t.AssertEqual(JWT.NewToken(users[0].ID, users[0].Username), UserJSON.User.Token)
	t.AssertEqual(bodyUser.User.Username, UserJSON.User.Username)
	t.AssertEqual(bodyUser.User.Email, UserJSON.User.Email)
}
