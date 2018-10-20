package models

import (
	"database/sql"
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
	TagList            []string `db:"-" json:"tagList"`
	User               *User    `db:"-" json:"author"`
	Favorited          bool     `db:"-" json:"favorited"`
	CreatedAtFormatted string   `db:"-" json:"createdAt"`
	UpdatedAtFormatted string   `db:"-" json:"updatedAt"`
	IsChangeTitle      bool     `db:"-" json:"-"`
}

func NewArticle(title, description, body string, tagList []string, user *User) *Article {
	article := &Article{Title: title, Description: description, Body: body}
	article.TagList = tagList
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
	var tagIds []string
	_, err := s.Select(&tagIds, "select TagID from ArticleTag where ArticleID=?", article.ID)
	if err != nil {
		panic(err)
	}
	if len(tagIds) > 0 {
		var tagListId []string
		for _, tagId := range tagIds {
			tagListId = append(tagListId, tagId)
		}
		for _, v := range tagListId {
			tag, err := s.SelectStr("select Name from Tag where ID=?", v)
			if err != nil {
				if err != sql.ErrNoRows {
					panic(err)
				}
				revel.TRACE.Println(err)
			}
			article.TagList = append(article.TagList, tag)
		}

	}
	return nil
}

func (article *Article) setSlug(s gorp.SqlExecutor) {
	if !article.IsChangeTitle && !article.CreatedAt.IsZero() {
		return
	}
	slugFromTitle := slugify.Slugify(article.Title)
	article.Slug = slugFromTitle
	var slugs []string
	_, err := s.Select(&slugs, "select `Slug` from `Article` where `Slug` LIKE ?", article.Slug+"%")
	if err != nil {
		panic(err)
	}
	revel.TRACE.Println(slugs)
	changed := true
	for changed {
		changed = false
		for k, slug := range slugs {
			if article.Slug == slug {
				article.Slug = slugFromTitle + "-" + strconv.Itoa(k)
				changed = true
				break
			}
		}
	}
}

func (article *Article) setTags(s gorp.SqlExecutor) {
	if len(article.TagList) == 0 {
		return
	}

	var tags []Tag
	var tag Tag
	var TagList []string
	for _, v := range article.TagList { // not supported 'in' https://github.com/go-gorp/gorp/issues/85
		err := s.SelectOne(&tag, "select * from Tag where Name=?", v)
		if err != nil {
			if err != sql.ErrNoRows {
				panic(err)
			}
			TagList = append(TagList, v)
			continue
		}
		tags = append(tags, tag)
	}
	revel.TRACE.Println("already exists tags", tags)
	revel.TRACE.Println("for create tags", TagList)
	var articleTags []*ArticleTag
	if len(tags) > 0 {
		for _, tag := range tags {
			articleTag := &ArticleTag{ArticleID: article.ID, TagID: tag.ID}
			articleTags = append(articleTags, articleTag)
		}
	}
	if len(TagList) > 0 {
		for _, name := range TagList {
			tag := &Tag{Name: name}
			if err := s.Insert(tag); err != nil {
				panic(err)
			}
			articleTag := &ArticleTag{ArticleID: article.ID, TagID: tag.ID}
			articleTags = append(articleTags, articleTag)
		}
	}
	for _, v := range articleTags {
		if err := s.Insert(v); err != nil {
			panic(err)
		}
	}

}

func (article *Article) PreInsert(s gorp.SqlExecutor) error {
	article.setSlug(s)
	article.CreatedAt = time.Now()
	article.UpdatedAt = article.CreatedAt
	article.CreatedAtFormatted = article.CreatedAt.UTC().Format(TimeFormat)
	article.UpdatedAtFormatted = article.CreatedAt.UTC().Format(TimeFormat)
	return nil
}

func (article *Article) PostInsert(s gorp.SqlExecutor) error {
	article.setTags(s)
	return nil
}

func (article *Article) PreUpdate(s gorp.SqlExecutor) error {
	article.setSlug(s)
	article.UpdatedAt = time.Now()
	article.UpdatedAtFormatted = article.UpdatedAt.UTC().Format(TimeFormat)
	return nil
}

func (article *Article) setUser(user *User) {
	article.UserID = user.ID
	article.User = user
}

func (article *Article) Fill(userJson *Article) {
	if userJson.Body != "" || userJson.Body != article.Body {
		article.Body = userJson.Body
	}
	if userJson.Description != "" || userJson.Description != article.Description {
		article.Description = userJson.Description
	}
	if userJson.Title != "" || userJson.Title != article.Title {
		article.Title = userJson.Title
		article.IsChangeTitle = true
	}
}
