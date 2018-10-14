package tests

import (
	"github.com/klim0v/golang-revel-realworld-starter-kit/app"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/lib/auth"
	"github.com/klim0v/golang-revel-realworld-starter-kit/app/models"
	"github.com/revel/revel/testing"
	"io"
	"net/http"
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

func (t *AppTest) Before() {
	println("Set up")
	query := "select * " + "from User "
	app.Dbm.Select(&users, query)
	JWT = auth.NewJWT()
}

func (t *AppTest) After() {
	println("Tear down")
}

func (t *AppTest) MakePostRequest(url string, body io.Reader, header http.Header) {
	request := t.PostCustom(t.BaseUrl()+url, "application/json", body)
	if header != nil {
		request.Header = header
	}

	request.Send()
}
