package goodreads_import

import (
	"borges.ai/data"
	"borges.ai/goodreads"
	"borges.ai/jobs"
	"borges.ai/models"
	"borges.ai/services"
	"borges.ai/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

// TODO when importing we need to save friends list, followers and following
// I can import even profile for the user who is not here. but there is not twitter id. when they signup how do I merge.
// too difficult:(

const GOODREADS_LIMIT = 200

func ImportUserDataFromGoodreads(db *gorm.DB, user models.User, isInitiatedByUser bool) error {
	defer utils.Duration(utils.Track("ImportUserDataFromGoodreads"))
	// TODO this should come as a token from the connected user
	grc := goodreads.NewClient(os.Getenv("GOODREADS_KEY"))
	grUser, err := grc.UserShow(user.GoodreadsID)
	if err != nil {
		log.WithError(err).Error("what happened")
		return err
	}

	shelvesIDs := make([]string, 0)
	shelvesNameCountMap := make(map[string]int, 0)
	for _, shelf := range grUser.UserShelves {
		shelvesIDs = append(shelvesIDs, shelf.ID)
		shelvesNameCountMap[strings.ToLower(shelf.Name)] = shelf.BookCount
	}

	log.WithField("user", grUser.Name).WithField("location", grUser.Location).WithField("user shelves", shelvesIDs).Info("user")
	var wg sync.WaitGroup

	wg.Add(1)
	go (func() {
		_, err = data.UpdateUserWithGoodreadsInfo(db, user, shelvesIDs, grUser.Link, grUser.Website, grUser.ReviewCount)
		if err != nil {
			log.WithError(err).Error("failed to save shelves")
		}
		wg.Done()
	})()
	// TODO later. fetch friends and following

	for _, shelfName := range goodreads.ShelvesNames {
		wg.Add(1)
		bookCount := shelvesNameCountMap[shelfName]
		ImportGoodreadsShelf(db, user, grc, shelfName, bookCount, isInitiatedByUser, &wg)
	}
	wg.Wait()
	return err
}

func calcPages(count, limit int) int {
	return int(math.Ceil(float64(count) / float64(limit)))
}

func ImportGoodreadsShelf(db *gorm.DB, user models.User, grc *goodreads.Client, goodreadsShelfName string, bookCount int, isInitiatedByUser bool, wg *sync.WaitGroup) {
	defer wg.Done()
	pagesNum := calcPages(bookCount, GOODREADS_LIMIT)

	var fetchWg sync.WaitGroup
	for i := 0; i <= pagesNum; i++ {
		fetchWg.Add(1)
		pageNum := i + 1 // 1 based
		ImportGoodreadsShelfPaged(db, user, grc, goodreadsShelfName, pageNum, isInitiatedByUser, &fetchWg)
	}
	fetchWg.Wait()
}

func ImportGoodreadsShelfPaged(db *gorm.DB, user models.User, grc *goodreads.Client, goodreadsShelfName string, pageNum int,isInitiatedByUser bool, wg *sync.WaitGroup) {
	defer wg.Done()
	// TODO this we can use access token!
	// we have total number of pages now
	reviews, err := grc.ReviewList(user.GoodreadsID, goodreadsShelfName, "", "", "", pageNum, GOODREADS_LIMIT)
	if err != nil {
		log.WithError(err).Error("what happened")
		return
	}
	log.WithField("reviews", len(reviews)).WithField("shelf", goodreadsShelfName).Info("books to be imported")
	var booksWg sync.WaitGroup
	for _, goodreadsBookReview := range reviews {
		booksWg.Add(1)
		go importGoodreadsBookAndStatus(db, user, goodreadsBookReview, goodreadsShelfName, isInitiatedByUser, &booksWg)
	}
	booksWg.Wait()
}
func makeTitleAndSubtitle(goodreadsTitle string) (string, string) {
	titleTokens := strings.Split(goodreadsTitle, ": ")
	bookTitle := titleTokens[0]
	var bookSubtitle string
	if len(titleTokens) > 1 {
		bookSubtitle = titleTokens[1]
	}
	titleTokensByBraces := strings.Split(bookTitle, " (")
	bookTitle = titleTokensByBraces[0]
	if len(titleTokensByBraces) > 1 {
		bookSubtitle = strings.ReplaceAll(titleTokensByBraces[1], ")", "")
	}
	// for this case Surely You're Joking, Mr. Feynman! Adventures of a Curious Character. when subtitle is also extracted
	bookTitle = strings.TrimSuffix(bookTitle, bookSubtitle)

	titleTokensByExclamation := strings.Split(bookTitle, "! ")
	bookTitle = titleTokensByExclamation[0]
	if len(titleTokensByExclamation) > 1 {
		bookTitle = bookTitle + "!" // add exclamation back if we did split by it
		bookSubtitle = titleTokensByExclamation[1]
	}
	return bookTitle, bookSubtitle
}

func ImportGoodreadsBookAndEdition(db *gorm.DB, user models.User, goodreadsBook goodreads.Book) (models.Book, error) {
	defer utils.Duration(utils.Track("ImportGoodreadsBookAndEdition"))
	book := models.Book{}
	authors := make([]models.Author, 0)
	for _, goodreadsAuthor := range goodreadsBook.Authors {
		author, err := data.FindOrCreateAuthor(db, user, goodreadsAuthor.ID, goodreadsAuthor.Name, goodreadsAuthor.Link)
		if err != nil {
			log.WithError(err).Error("failed to create an author")
			return book, err
		}
		author.Role = goodreadsAuthor.Role
		authors = append(authors, author)
	}

	bookTitle, bookSubtitle := makeTitleAndSubtitle(goodreadsBook.Title)


	isbn13 := goodreadsBook.ISBN13
	isbn := goodreadsBook.ISBN
	asin := goodreadsBook.ASIN
	publishedYear := goodreadsBook.Published
	publisherName := goodreadsBook.Publisher
	numPages := goodreadsBook.NumPages
	isEbook := goodreadsBook.IsEbook
	goodreadsBookID := goodreadsBook.Work.ID
	goodreadsBestEditionID := goodreadsBook.Work.BestBookID
	goodreadsID := goodreadsBook.ID
	titleWithoutSeries := goodreadsBook.TitleWithoutSeries
	description := goodreadsBook.Description
	book, err := services.ImportBook(db, user, goodreadsBookID, goodreadsBestEditionID, bookTitle, bookSubtitle, titleWithoutSeries, description, authors)
	if err != nil {
		log.WithError(err).Error("failed to create a book")
		return book, err
	}
	bookEdition, err := services.ImportBookEdition(db, user, book.ID, goodreadsBookID, goodreadsID, isbn, isbn13, asin, publisherName, numPages, publishedYear, isEbook, goodreadsBook.Format)
	if err != nil {
		log.WithError(err).Error("failed to create book edition")
		// NOTE do not return. Not essential. Just log
	}

	err = jobs.SyncGoodreadsBook(user, book.ID, bookEdition.GoodreadsID)
	if err != nil {
		log.WithError(err).Error("failed to publish message to sync book")
	}
	if bookEdition.ID > 0 {
		book.Edition = bookEdition
	} else {
		book.Edition = models.BookEdition{}
	}
	if len(authors) > 0 {
		book.Authors = authors
	}
	// this will trigger short id re-population
	book.AfterFind()
	return book, err
}

func importGoodreadsBookAndStatus(db *gorm.DB, user models.User, goodreadsReview goodreads.Review, goodreadsShelfName string,isInitiatedByUser bool, wg *sync.WaitGroup) {
	defer wg.Done()
	book, err := ImportGoodreadsBookAndEdition(db, user, goodreadsReview.Book)
	if err != nil || book.ID == 0 {
		log.WithError(err).WithField("goodreadsBook", goodreadsReview.Book.Title).Error("failed to import book or book edition")
		return
	}

	addedAtStr := strings.TrimSpace(goodreadsReview.DateAdded)
	startedAtStrOriginal := strings.TrimSpace(goodreadsReview.StartedAt)
	finishedAtStrOriginal := strings.TrimSpace(goodreadsReview.ReadAt)

	startedAtStr := startedAtStrOriginal
	finishedAtStr := finishedAtStrOriginal

	addedAt := goodreads.ParseDate(addedAtStr)
	startedAt := goodreads.ParseDate(startedAtStr)
	finishedAt := goodreads.ParseDate(finishedAtStr)

	// reformat dates in "2006-02-01" format
	if !startedAt.IsZero() {
		startedAtStr = utils.FormatDate(startedAt)
		// fallback
		if startedAtStr == "" {
			startedAtStr = startedAtStrOriginal
		}
	}

	if !finishedAt.IsZero() {
		finishedAtStr = utils.FormatDate(finishedAt)
		// fallback
		if finishedAtStr == "" {
			finishedAtStr = finishedAtStrOriginal
		}
	}

	// TODO check what to do if dates is only in the year.
	// also it's possible to have end date
	if goodreadsShelfName == goodreads.READ_SHELF {
		_, err = services.ImportReading(db, user, book.ID, book.Edition.ID, goodreadsReview.ID, startedAtStr, finishedAtStr, startedAt, finishedAt)
		if err != nil {
			log.WithError(err).Error("failed to create a book reading")
		}
	}

	reviewContent := strings.TrimSpace(goodreadsReview.Body)
	rating := goodreadsReview.Rating

	var bookStatusDate time.Time
	var bookStatus int
	if goodreadsShelfName == goodreads.READ_SHELF {
		readStatusDate := finishedAt
		if readStatusDate.IsZero() {
			readStatusDate = addedAt
		}
		bookStatus = data.READ_STATUS
		bookStatusDate = readStatusDate
	} else if goodreadsShelfName == goodreads.CURRENTLY_READING_SELF {
		bookStatus = data.STARTED_STATUS
		bookStatusDate = startedAt
	} else if goodreadsShelfName == goodreads.TO_READ_SHELF {
		bookStatus = data.WANT_TO_READ_STATUS
		bookStatusDate = addedAt
	}
	_, err = services.ImportReview(db, user, book.ID, book.Edition.ID, goodreadsReview.ID, reviewContent, rating, bookStatus, bookStatusDate, isInitiatedByUser)
}
