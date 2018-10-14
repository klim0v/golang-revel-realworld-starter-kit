
create table Article
(
	ID int auto_increment
		primary key,
	Slug varchar(255) null,
	Title varchar(255) null,
	Description varchar(255) null,
	Body varchar(255) null,
	UserID int null,
	FavoritesCount int null,
	CreatedAt datetime null,
	UpdatedAt datetime null,
	constraint Slug
		unique (Slug)
)
charset=utf8
;

create table Comment
(
	ID int auto_increment
		primary key,
	Body varchar(255) null,
	ArticleID int null,
	UserID int null,
	CreatedAt datetime null,
	UpdatedAt datetime null
)
charset=utf8
;

create table Favorite
(
	ID int auto_increment
		primary key,
	UserID int null,
	ArticleID int null
)
charset=utf8
;

create table Tag
(
	ID int unsigned auto_increment
		primary key,
	Name varchar(255) null,
	TaggingsCount int unsigned null,
	constraint Name
		unique (Name)
)
charset=utf8
;

create table User
(
	ID int auto_increment
		primary key,
	CreatedAt datetime null,
	UpdatedAt datetime null,
	Username varchar(255) null,
	Email varchar(255) null,
	Bio varchar(255) null,
	Image varchar(255) null,
	HashedPassword mediumblob null,
	constraint Email
		unique (Email),
	constraint Username
		unique (Username)
)
charset=utf8
;

