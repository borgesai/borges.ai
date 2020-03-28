package goodreads

type BestBook struct {
	ID     string `xml:"id"`
	Title  string `xml:"title"`
	Author Author `xml:"author"`
}

type Work struct {
	ID                       string   `xml:"id"`
	Book                     BestBook `xml:"best_book"`
	OriginalPublicationYear  int      `xml:"original_publication_year"`
	OriginalPublicationMonth int      `xml:"original_publication_month"`
	OriginalPublicationDay   int      `xml:"original_publication_day"`
	OriginalTitle            string   `xml:"original_title"`
	BestBookID               string   `xml:"best_book_id"`
}

type Author struct {
	ID               string  `xml:"id"`
	Name             string  `xml:"name"`
	ImageURL         string  `xml:"image_url"`
	SmallImageURL    string  `xml:"small_image_url"`
	LargeImageURL    string  `xml:"large_image_url"`
	Link             string  `xml:"link"`
	Role             string  `xml:"role"`
	AverageRating    float32 `xml:"average_rating"`
	RatingsCount     int     `xml:"ratings_count"`
	TextReviewsCount int     `xml:"text_reviews_count"`
	FansCount        int     `xml:"fans_count"`
	AuthorFollowers  int     `xml:"author_followers"`
	About            string  `xml:"about"`
	WorksCount       int     `xml:"works_count"`
	Gender           string  `xml:"gender"`
	Hometown         string  `xml:"hometown"`
	BornAt           string  `xml:"born_at"`
	DiedAt           string  `xml:"died_at"`
	GoodreadsAuthor  bool    `xml:"goodreads_author"`
	UserID           string  `xml:"user>user_id"`
	Books            []Book  `xml:"books>book"`
}

type Book struct {
	ID                 string   `xml:"id"`
	ISBN               string   `xml:"isbn"`
	ISBN13             string   `xml:"isbn13"`
	ASIN               string   `xml:"asin"`
	TextReviewsCount   int      `xml:"text_reviews_count"`
	URI                string   `xml:"uri"`
	Title              string   `xml:"title"`
	TitleWithoutSeries string   `xml:"title_without_series"`
	ImageURL           string   `xml:"image_url"`
	SmallImageURL      string   `xml:"small_image_url"`
	LargeImageURL      string   `xml:"large_image_url"`
	Link               string   `xml:"link"`
	NumPages           int      `xml:"num_pages"`
	Format             string   `xml:"format"`
	EditionInformation string   `xml:"edition_information"`
	Publisher          string   `xml:"publisher"`
	Published          int      `xml:"published"`
	PublicationDay     int      `xml:"publication_day"`
	PublicationYear    int      `xml:"publication_year"`
	PublicationMonth   int      `xml:"publication_month"`
	AverageRating      float32  `xml:"average_rating"`
	RatingsCount       int      `xml:"ratings_count"`
	Description        string   `xml:"description"`
	Authors            []Author `xml:"authors>author"`
	Work               Work     `xml:"work"`
	IsEbook            bool     `xml:"is_ebook"`
	PopularShelves     []Shelf  `xml:"popular_shelves>shelf"`
}
type Shelf struct {
	Name  string `xml:"name,attr"`
	Count string  `xml:"count,attr"`
}

type Review struct {
	ID          string `xml:"id"`
	Book        Book   `xml:"book"`
	Rating      int    `xml:"rating"`
	StartedAt   string `xml:"started_at"`
	ReadAt      string `xml:"read_at"`
	DateAdded   string `xml:"date_added"`
	DateUpdated string `xml:"date_updated"`
	ReadCount   int    `xml:"read_count"`
	Body        string `xml:"body"`
}
type ReviewContent struct {
	Review string `xml:"review"`
}

type UserReview struct {
	ID     string `xml:"id"`
	Review string `xml:"review[review]"`
	//Review ReviewContent `xml:"review"`
}

type AuthUser struct {
	ID int `xml:"id,attr"`
}
type User struct {
	ID            string      `xml:"id"`
	Name          string      `xml:"name"`
	Link          string      `xml:"link"`
	ImageURL      string      `xml:"image_url"`
	SmallImageURL string      `xml:"small_image_url"`
	About         string      `xml:"about"`
	Gender        string      `xml:"gender"`
	Location      string      `xml:"location"`
	Website       string      `xml:"website"`
	Joined        string      `xml:"joined"`
	LastActive    string      `xml:"last_active"`
	FriendsCount  int         `xml:"friends_count"`
	GroupsCount   int         `xml:"groups_count"`
	ReviewCount   int         `xml:"reviews_count"`
	UserShelves   []UserShelf `xml:"user_shelves>user_shelf"`
}

type UserShelf struct {
	ID            string `xml:"id"`
	Name          string `xml:"name"`
	BookCount     int    `xml:"book_count"`
	ExclusiveFlag bool   `xml:"exclusive_flag"`
	Description   string `xml:"description"`
}
