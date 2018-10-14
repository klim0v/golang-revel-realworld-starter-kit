package models

import (
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v2"
	"regexp"
	"time"
)

type User struct {
	ID             int       `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Bio            string    `json:"bio"`
	Image          string    `json:"image"`
	HashedPassword []byte    `json:"-"`

	// Transient
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty"`
}

func (user *User) String() string {
	return user.Username
}

func NewUser(username, email, password string) *User {
	user := &User{Email: email, Username: username}
	user.setPassword(password)
	return user
}

func (user *User) setPassword(password string) {
	user.Password = password
	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{Max: 15},
		revel.MinSize{Min: 4},
		revel.Match{Regexp: userRegex},
	).Key("user.username")

	v.Check(user.Email,
		revel.Required{},
		revel.ValidEmail(),
	).Key("user.email")

	if user.CreatedAt.IsZero() || user.Password != "" {
		ValidatePassword(v, user.Password).
			Key("user.password")
	}
}

func (user *User) MatchPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	return err == nil
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{Max: 15},
		revel.MinSize{Min: 5},
	)
}

func (user *User) PreInsert(s gorp.SqlExecutor) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	return nil
}

func (user *User) PreUpdate(s gorp.SqlExecutor) error {
	user.UpdatedAt = time.Now()
	return nil
}
func (user *User) Fill(userJson *User) {
	user.Email = userJson.Email
	user.Username = userJson.Username
	user.Bio = userJson.Bio
	user.Image = userJson.Image
	if userJson.Password != "" {
		user.setPassword(userJson.Password)
	}
}
