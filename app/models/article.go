package models

import (
	"github.com/Machiel/slugify"
	"github.com/revel/revel"
	"gopkg.in/gorp.v2"
	"strconv"
	"time"
)

type Article struct {
	ID             int       `json:"-"`
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	UserID         int       `json:"-"`
	FavoritesCount int       `json:"favoritesCount"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`

	// Transient
	TagList            []string `json:"tagList"`
	User               *User    `json:"author"`
	Favorited          bool     `json:"favorited"`
	CreatedAtFormatted string   `json:"createdAt"`
	UpdatedAtFormatted string   `json:"updatedAt"`
}

func NewArticle(title, description, body string, tagList []string, user *User) *Article {
	article := &Article{Title: title, Description: description, Body: body}
	article.setTagList(tagList)
	article.setUser(user)
	return article
}

func (article *Article) Validate(v *revel.Validation) {
	v.Required(article.Title).Key("title").Message(EmptyMsg)
	v.Required(article.Description).Key("description").Message(EmptyMsg)
	v.Required(article.Body).Key("body").Message(EmptyMsg)
}

func (article *Article) PostGet(s gorp.SqlExecutor) error {
	article.CreatedAtFormatted = article.CreatedAt.UTC().Format(TimeFormat)
	article.UpdatedAtFormatted = article.UpdatedAt.UTC().Format(TimeFormat)
	user, _ := s.Get(User{}, article.UserID)
	article.User = user.(*User)
	return nil
}

func (article *Article) setSlug(s gorp.SqlExecutor, slugFromTitle string) {
	article.Slug = slugFromTitle
	var slugs []string
	_, err := s.Select(&slugs, "select `Slug` from `Article` where `Slug` LIKE ?", article.Slug+"%")
	if err != nil {
		panic(err)
	}
	revel.TRACE.Println(slugs)
	for k, slug := range slugs {
		if article.Slug == slug {
			article.Slug = slugFromTitle + "-" + strconv.Itoa(k)
		}
	}
}
func (article *Article) PreInsert(s gorp.SqlExecutor) error {
	article.CreatedAt = time.Now()
	article.UpdatedAt = article.CreatedAt
	article.CreatedAtFormatted = article.CreatedAt.UTC().Format(TimeFormat)
	article.UpdatedAtFormatted = article.CreatedAt.UTC().Format(TimeFormat)
	slugFromTitle := slugify.Slugify(article.Title)
	article.setSlug(s, slugFromTitle)
	return nil
}

func (article *Article) PreUpdate(s gorp.SqlExecutor) error {
	article.UpdatedAt = time.Now()
	article.UpdatedAtFormatted = article.UpdatedAt.UTC().Format(TimeFormat)
	slugFromTitle := slugify.Slugify(article.Title)
	article.setSlug(s, slugFromTitle)
	return nil
}
func (article *Article) setTagList(tagList []string) {
	article.TagList = tagList
}

func (article *Article) setUser(user *User) {
	article.UserID = user.ID
	article.User = user
}

func (article *Article) Fill(userJson *Article) {
	if userJson.Body != "" {
		article.Body = userJson.Body
	}
	if userJson.Description != "" {
		article.Description = userJson.Description
	}
	if userJson.Title != "" {
		article.Title = userJson.Title
	}
}
