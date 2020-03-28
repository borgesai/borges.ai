package routes

import "borges.ai/models"

type IndexPageData struct {
	Page           string
	Title          string
	User           models.User
	Books          []models.Book
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}

type ErrorPageData struct {
	Page           string
	Title          string
	User           models.User
	Message        string
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}

type UserPageData struct {
	Page   string
	Title  string
	IsMine bool

	Readonly       bool
	User           models.User
	Profile        models.User
	Books          []models.Book
	ShowActions    bool
	AverageRating  float32
	ReviewsCount   int
	Submenu        string
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}
type UserChartsPageData struct {
	Page   string
	Title  string
	IsMine bool

	Readonly       bool
	User           models.User
	Profile        models.User
	Books          map[int]models.IndexedData
	MaxBooks       int
	MaxPages       int
	Submenu        string
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}

type UserSettingsPageData struct {
	Page           string
	Title          string
	IsMine         bool
	Readonly       bool
	User           models.User
	Profile        models.User
	Submenu        string
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}

type SearchPageData struct {
	Page           string
	Title          string
	User           models.User
	Books          []models.Book
	Query          string
	SearchedByISBN bool
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}

type UserBookPageData struct {
	Page           string
	Title          string
	IsMine         bool
	Readonly       bool
	User           models.User
	Profile        models.User
	Book           models.Book
	SessionUserReview models.Review
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}

type BookPageData struct {
	Page           string
	Title          string
	User           models.User
	Book           models.Book
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}

type AuthorPageData struct {
	Page           string
	Title          string
	User           models.User
	Author         models.Author
	Books          []models.Book
	ShowFooter     bool
	IsCustomDomain bool
	CanonicalURL   string
}
