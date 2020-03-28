package data

import (
	"borges.ai/text"
	"strconv"
	"time"

	"borges.ai/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func FindOrCreateUserFromTwitter(db *gorm.DB, twitterID int64, username, email, twitterAccessToken, name, twitterLocation,
	profileTextColor, profileBackgroundColor, profileLinkColor, profileSidebarBorderColor, profileSidebarFillColor string) (models.User, error) {
	record := models.User{}
	record.TwitterID = twitterID
	record.Email = email
	record.TwitterAccessToken = twitterAccessToken
	record.Username = username
	record.Name = name
	record.TwitterLocation = twitterLocation
	if name != "" {
		// we can only do latin
		slug := text.Slug(name)
		record.AbbreviatedName = text.Abbreviate(slug)
	} else {
		record.AbbreviatedName = ":)"
	}
	record.ProfileTextColor = profileTextColor
	record.ProfileBackgroundColor = profileBackgroundColor
	record.ProfileLinkColor = profileLinkColor
	record.ProfileSidebarBorderColor = profileSidebarBorderColor
	record.ProfileSidebarFillColor = profileSidebarFillColor

	record.UpdatedAt = time.Now()
	err := db.Where(models.User{TwitterID: twitterID}).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a user")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created user")
	return record, nil
}

func FindOrCreateUserFromGoodreads(db *gorm.DB, goodreadsID int, goodreadsAccessToken, goodreadsAccessSecret, name, link, website string, reviewCount int) (models.User, error) {
	record := models.User{}
	record.GoodreadsID = goodreadsID
	record.GoodreadsAccessToken = goodreadsAccessToken
	record.GoodreadsAccessSecret = goodreadsAccessSecret
	record.GoodreadsURL = link
	record.WebsiteURL = website
	record.GoodreadsBooksCount = reviewCount
	//
	if name != "" {
		// we can only do latin
		record.Name = name
		slug := text.Slug(name)
		record.AbbreviatedName = text.Abbreviate(slug)
	} else {
		record.AbbreviatedName = ":)"
	}
	record.UpdatedAt = time.Now()
	err := db.Where(models.User{GoodreadsID: goodreadsID}).Assign(record).FirstOrCreate(&record).Error
	if err != nil {
		log.WithField("record", record).WithError(err).Error("failed to find or create a user")
		return record, err
	}
	log.WithField("record", record.ID).Info("found or created user")
	return record, nil
}

func FindUserByID(db *gorm.DB, id uint64) (models.User, error) {
	record := models.User{}
	err := db.First(&record, id).Error
	log.WithField("record", record).Debug("found by id")
	return record, err
}

func FindUserByUsername(db *gorm.DB, username string) (models.User, error) {
	record := models.User{}
	err := db.Where("username=?", username).First(&record).Error
	log.WithField("user", record).Debug("found user by username")
	return record, err
}

func FindUserByTwitterID(db *gorm.DB, twitterID string) (models.User, error) {
	record := models.User{}
	err := db.Where("twitter_id=?", twitterID).First(&record).Error
	log.WithField("user", record).Debug("found user by twitter id")
	return record, err
}

func FindUserByGoodreadsID(db *gorm.DB, goodreadsID int) (models.User, error) {
	record := models.User{}
	err := db.Where("goodreads_id=?", goodreadsID).First(&record).Error
	log.WithField("user", record).Debug("found user by goodreads id")
	return record, err
}

func FindUserByCustomDomain(db *gorm.DB, domain string) (models.User, error) {
	record := models.User{}
	err := db.Where("custom_domain=?", domain).First(&record).Error
	log.WithField("user", record).Debug("found user by custom domain")
	return record, err
}

func UpdateUserWithGoodreadsInfo(db *gorm.DB, user models.User, goodreadsShelvesIDs []string, goodreadsURL, websiteURL string, goodreadsBooksCount int) (models.User, error) {
	update := models.User{
		GoodreadsShelvesIDs: goodreadsShelvesIDs,
		GoodreadsURL:        goodreadsURL,
		GoodreadsBooksCount: goodreadsBooksCount,
		WebsiteURL:          websiteURL,
		GoodreadsSyncEnabled:  true,
	}
	err := db.Debug().Model(&user).Update(update).Error
	return user, err
}

func UpdateUserWithDropboxInfo(db *gorm.DB, user models.User, dropboxID, dropboxAccessToken string) (models.User, error) {
	update := models.User{
		DropboxID: dropboxID,
		DropboxAccessToken: dropboxAccessToken,
	}
	err := db.Debug().Model(&user).Update(update).Error
	return user, err
}

func UpdateUserWithDropboxInfoAndEmail(db *gorm.DB, user models.User, dropboxID, dropboxAccessToken, email string) (models.User, error) {
	update := models.User{
		DropboxID: dropboxID,
		DropboxAccessToken: dropboxAccessToken,
		Email: email,
	}
	err := db.Debug().Model(&user).Update(update).Error
	return user, err
}

func UpdateUserUsername(db *gorm.DB, user models.User) (models.User, error) {
	idStr := strconv.FormatUint(user.ID, 10)
	idInt, _ := strconv.ParseInt(idStr, 10, 64)
	username :=  text.Slug(user.Name) + "-" + text.NumberToHash(idInt)
	update := models.User{
		Username: username,
	}
	err := db.Debug().Model(&user).Update(update).Error
	return user, err
}

func UpdateUserGoodreadsAccountInfo(db *gorm.DB, user models.User, goodreadsID int, accessToken, accessSecret string) (models.User, error) {
	update := models.User{
		GoodreadsID:           goodreadsID,
		GoodreadsAccessToken:  accessToken,
		GoodreadsAccessSecret: accessSecret,
		GoodreadsSyncEnabled:  true,
	}
	err := db.Debug().Model(&user).Update(update).Error
	return user, err
}

func UpdateUserGoodreadsSyncFlag(db *gorm.DB, user models.User, enabled bool) (models.User, error) {
	err := db.Debug().Model(&user).Update("goodreads_sync_enabled", enabled).Error
	return user, err
}
