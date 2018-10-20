package models

type Tag struct {
	ID   int
	Name string
}

type ArticleTag struct {
	ID        int
	ArticleID int
	TagID     int
}
