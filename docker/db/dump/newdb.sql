-- we don't know how to generate schema database (class Schema) :(
create table Article
(
	ID int auto_increment
		primary key,
	Slug varchar(255) not null,
	Title varchar(255) not null,
	Description varchar(255) not null,
	Body varchar(255) not null,
	UserID int not null,
	FavoritesCount int null,
	CreatedAt datetime not null,
	UpdatedAt datetime not null,
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
	ID int auto_increment
		primary key,
	Name varchar(255) null,
	constraint Name
		unique (Name)
)
charset=utf8
;

create table ArticleTag
(
	ID int auto_increment
		primary key,
	ArticleID int null,
	TagID int null,
	constraint ArticleIDTagID
		unique (ArticleID, TagID),
	constraint ArticleID
		foreign key (ArticleID) references Article (ID)
			on delete cascade,
	constraint TagID
		foreign key (TagID) references Tag (ID)
			on delete cascade
)
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

