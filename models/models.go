package models

import (
	"borges.ai/text"
	"github.com/jinzhu/gorm"
	"html/template"
	"strconv"

	//"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func NewDB() (*gorm.DB, error) {
	stage := os.Getenv("ENV")
	var db *gorm.DB
	var err error
	if stage == "prod" {
		log.Debug("trying to connect to production")
		db.DB().SetMaxOpenConns(100)
		db.DB().SetMaxIdleConns(0)
	} else {
		connectionString := "host=" + os.Getenv("DB_HOST") + " port=5432 dbname=" + os.Getenv("DB_NAME") + " user=" + os.Getenv("DB_USERNAME") + " sslmode=disable password=" + os.Getenv("DB_PASSWORD")
		db, err = gorm.Open("postgres", connectionString)
		db.DB().SetMaxOpenConns(300)
		db.DB().SetMaxIdleConns(0)
	}
	if err != nil {
		log.WithError(err).Error("failed to connect")
	}

	db.LogMode(true)
	return db, err
}

type Model struct {
	ID        uint64     `gorm:"primary_key"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `sql:"index"`
	CreatorID uint64     `gorm:"index:creator_id;not null"`
}

// Authorable embedded model
type Authorable struct {
	AuthorsIDs      pq.StringArray `gorm:"type:varchar(100)[]"`
	AuthorsCache    pq.StringArray `gorm:"type:varchar(100)[]"`
	AuthorsCacheStr string
	AuthorsIDsUint  []uint64 `gorm:"-"`
	Authors         []Author `gorm:"-"`

	IllustratorsIDs      pq.StringArray `gorm:"type:varchar(100)[]"`
	IllustratorsCache    pq.StringArray `gorm:"type:varchar(100)[]"`
	IllustratorsCacheStr string

	TranslatorsIDs      pq.StringArray `gorm:"type:varchar(100)[]"`
	TranslatorsCache    pq.StringArray `gorm:"type:varchar(100)[]"`
	TranslatorsCacheStr string
}

type User struct {
	Model
	TwitterID             int64
	GoodreadsID           int
	Username              string         `gorm:"type:varchar(100)"`
	Email                 string         `gorm:"type:varchar(100)"`
	Name                  string         `gorm:"type:varchar(100)"`
	TwitterAccessToken    string         `gorm:"type:varchar(200)"`
	TwitterLocation       string         `gorm:"type:varchar(100)"`
	GoodreadsAccessToken  string         `gorm:"type:varchar(200)"`
	GoodreadsAccessSecret string         `gorm:"type:varchar(200)"`
	AbbreviatedName       string         `gorm:"type:varchar(20)"`
	GoodreadsShelvesIDs   pq.StringArray `gorm:"type:varchar(100)[]"`
	GoodreadsURL          string         `gorm:"type:varchar(300)"`
	GoodreadsBooksCount   int
	GoodreadsSyncEnabled  bool

	DropboxID                 string `gorm:"type:varchar(200)"`
	DropboxAccessToken        string `gorm:"type:varchar(200)"`
	WebsiteURL                string `gorm:"type:varchar(300)"`
	ProfileTextColor          string
	ProfileBackgroundColor    string
	ProfileLinkColor          string
	ProfileSidebarBorderColor string
	ProfileSidebarFillColor   string
	CustomDomain              string `gorm:"index"`
}

type Book struct {
	Model
	GoodreadsID            string `gorm:"type:varchar(100)"`
	BestEditionGoodreadsID string `gorm:"type:varchar(100)"`
	Authorable
	Title              string `gorm:"type:varchar(500);index"`
	Slug               string `gorm:"type:varchar(500);index"`
	ShortID            string `gorm:"-"`
	Subtitle           string `gorm:"type:varchar(500)"`
	SearchableText     string `gorm:"type:tsvector"`
	TitleWithoutSeries string `gorm:"type:varchar(500)"`

	Description string

	OriginalTitle string `gorm:"type:varchar(500);index"`
	OriginalYear  int
	// approximate
	NumPages         int
	GoodreadsShelves pq.StringArray `gorm:"type:varchar(100)[]"`
	// holder for data user model
	Readings []Reading     `gorm:"-"`
	Edition  BookEdition   `gorm:"-"`
	Editions []BookEdition `gorm:"-"`
	Review   Review        `gorm:"-"`
	Reading  Reading       `gorm:"-"`
}

func (b *Book) GetNumPages() int {
	var numPages int
	if b.Edition.ID > 0 && b.Edition.NumPages > 0 {
		numPages = b.Edition.NumPages
	} else {
		numPages = b.NumPages
	}
	log.WithField("num pages", numPages).Info("got num of pages")
	return numPages
}

func (b *Book) AfterFind() (err error) {
	if b.ID > 0 {
		idStr := strconv.FormatUint(b.ID, 10)
		idInt, _ := strconv.ParseInt(idStr, 10, 64)
		b.ShortID = text.NumberToHash(idInt)
	}
	if len(b.AuthorsIDs) > 0 {
		ids := make([]uint64, 0)
		for _, idStr := range b.AuthorsIDs {
			id, _ := strconv.ParseUint(idStr, 10, 64)
			if id > 0 {
				ids = append(ids, id)
			}
		}
		b.AuthorsIDsUint = ids
	}
	return
}

type Author struct {
	Model
	GoodreadsID     string `gorm:"type:varchar(100)"`
	Name            string
	Slug            string `gorm:"type:varchar(100);unique_index"`
	AbbreviatedName string
	ImageURL        string
	WebsiteURL      string
	// non persistent
	Role string `gorm:"-"`
}

type BookEdition struct {
	Model
	GoodreadsID     string `gorm:"type:varchar(100);index"`
	GoodreadsBookID string `gorm:"type:varchar(100);index"`
	BookID          uint64 `gorm:"index"`
	ISBN            string `gorm:"type:varchar(20);index"`
	ISBN13          string `gorm:"type:varchar(20);index"`
	ASIN            string `gorm:"type:varchar(20);index"`
	Publisher       string `gorm:"type:varchar(200);"`
	PublicationYear int
	NumPages        int
	IsEbook         bool
	Format          string `gorm:"type:varchar(50)"`
}

type Review struct {
	Model
	Content            string
	ContentHash        string `gorm:"type:varchar(1000)"`
	ContentHTML        string
	ContentHTMLEscaped template.HTML `gorm:"-"`

	Rating     int
	Status     int
	StatusDate time.Time

	BookID uint64
	// optional
	BookEditionID uint64
	// optional
	GoodreadsID string `gorm:"type:varchar(100)"`
}

func (m *Review) AfterFind() (err error) {
	if m.ContentHTML != "" {
		m.ContentHTMLEscaped = template.HTML(m.ContentHTML)
	}
	return
}

type Reading struct {
	Model
	BookID uint64
	// optional
	BookEditionID uint64
	// optional
	GoodreadsID string `gorm:"type:varchar(100);index'"`

	StartDateRaw  string
	FinishDateRaw string

	StartDate  time.Time
	FinishDate time.Time

	Duration int    `gorm:"-"`
	Note     string `gorm:"type:varchar(1000)"`
}

type Job struct {
	Model
	Type       string
	StartDate  time.Time
	FinishDate time.Time
	Success    bool
	ModelType  string `gorm:"type:varchar(20)"`
	ModelID    uint64
	Error      string
}

func (m *Reading) AfterFind() (err error) {
	if !m.StartDate.IsZero() && !m.FinishDate.IsZero() {
		// TODO ensure finish date is larger than start date!
		days := m.FinishDate.Sub(m.StartDate).Hours() / 24
		m.Duration = int(days) + 1
	}
	return
}
