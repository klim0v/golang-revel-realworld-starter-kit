package tests

import (
	"github.com/klim0v/golang-revel-realworld-starter-kit/app"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/services/auth"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel/testing"
	"golang.org/x/crypto/bcrypt"
	"io"
)

type AppTest struct {
	testing.TestSuite
}
type ErrorJSON struct {
	Errors map[string][]string `json:"errors"`
}

var (
	JWT   auth.Tokener
	users []*models.User
)

const (
	demoUsername    = "DemoUser"
	demoEmail       = "user@demo.ru"
	demoPassword    = "DemoPassword"
	demoRegUsername = "DemoUserReg"
	demoRegEmail    = "user-reg@demo.ru"
	demoRegPassword = "DemoRegPassword"
)

func (t *AppTest) Before() {
	println("Set up")

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(demoPassword), bcrypt.DefaultCost)

	users = []*models.User{
		{Username: demoUsername, Email: demoEmail, Password: demoPassword, HashedPassword: bcryptPassword},
		{Username: "foo" + demoUsername, Email: "foo" + demoEmail, Password: "foo" + demoPassword, HashedPassword: bcryptPassword},
	}
	for _, user := range users {
		if err := app.Dbm.Insert(user); err != nil {
			panic(err)
		}
	}
	JWT = auth.NewJWT()
}

func (t *AppTest) After() {
	println("Tear down")
	for _, user := range users {
		app.Dbm.Delete(user)
	}
	app.Dbm.Exec("delete from User where Username=? and Email=?", demoRegUsername, demoRegEmail)
}

func (t *AppTest) TestConnection() {
	t.Assert(app.Dbm.Db.Ping() == nil)
	t.AssertEqual(2, len(users))
	for _, user := range users {
		count, err := app.Dbm.SelectInt("select count(*) from User where Username=? and Email=?", user.Username, user.Email)
		t.Assert(err == nil)
		t.Assert(count == 1)
	}
}

func (t *AppTest) MakePostRequest(url string, body io.Reader, header interface{}) {
	request := t.PostCustom(t.BaseUrl()+url, "application/json", body)
	if header != nil {
		request.Header.Set("Authorization", header.(string))
	}

	request.Send()
}

func (t *AppTest) MakePutRequest(url string, body io.Reader, header interface{}) {
	request := t.PutCustom(t.BaseUrl()+url, "application/json", body)
	if header != nil {
		request.Header.Set("Authorization", header.(string))
	}
	request.Send()
}
