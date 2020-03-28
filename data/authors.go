package data

import (
	"borges.ai/models"
	"borges.ai/text"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func FindOrCreateAuthor(db *gorm.DB, user models.User, goodreadsID, name, websiteURL string) (models.Author, error) {
	record := models.Author{}
	sanitizedName := strings.TrimSpace(strings.Split(name, " (")[0])
	record.Name = sanitizedName
	slug := text.Slug(sanitizedName)
	record.Slug = slug
	record.GoodreadsID = goodreadsID
	// we can only do latin
	record.AbbreviatedName = text.Abbreviate(slug)
	record.WebsiteURL = websiteURL

	record.UpdatedAt = time.Now()
	record.CreatorID = user.ID
	err := db.Debug().Where(models.Author{GoodreadsID: goodreadsID}).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			log.WithField("goodreadsID", goodreadsID).WithField("slug", slug).WithField("record", record).WithError(err).Info("failed to find or create an author (duplicate)")
			found, foundErr := FindAuthorBySlug(db, slug)
			if foundErr !=nil {
				if foundErr.Error() == "record not found" {
					// this is here for the case like The School of Life that has multiple goodreads ids
					return FindAuthorByGoodreadsID(db, goodreadsID)
				}
			}
			return found, foundErr
		}
		log.WithField("record", record).WithError(err).Error("failed to find or create an author")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created author")
	return record, nil
}

func FindAuthorsByIDs(db *gorm.DB, user models.User, authorIDs []uint64) ([]models.Author, error) {
	authors := []models.Author{}
	err := db.Where("id IN (?)", authorIDs).Find(&authors).Error
	return authors, err
}

func FindAuthorsByIDsAsMap(db *gorm.DB, user models.User, authorIDs []uint64) (map[uint64]models.Author, error) {
	authors, err := FindAuthorsByIDs(db, user, authorIDs)
	authorsMap := make(map[uint64]models.Author, 0)
	if err != nil {
		return authorsMap, err
	}
	for _, author := range authors {
		authorsMap[author.ID] = author
	}
	return authorsMap, err
}

func FindAuthorBySlug(db *gorm.DB, slug string) (models.Author, error) {
	record := models.Author{}
	err := db.Where("slug=?", slug).First(&record).Error
	log.WithField("user", record).Debug("found author by slug")
	return record, err
}

func FindAuthorByGoodreadsID(db *gorm.DB, goodreadsID string) (models.Author, error) {
	record := models.Author{}
	err := db.Where("goodreads_id=?", goodreadsID).First(&record).Error
	log.WithField("user", record).Debug("found author by goodreads_id")
	return record, err
}
