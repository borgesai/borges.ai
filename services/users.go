package services

import (
	"borges.ai/data"
	"borges.ai/goodreads"
	"borges.ai/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"os"
)

func FindOrCreateGoodreadsUser(db *gorm.DB, user models.User, goodreadsUserID int, accessToken, accessSecret string) (models.User, bool, error) {
	existed := false
	if user.ID > 0 && user.GoodreadsID > 0 {
		existed = true

	}
	if user.ID == 0 {
		user, _ = data.FindUserByGoodreadsID(db, goodreadsUserID)
		existed = false
	}
	// user exist. we just need to update
	if user.ID > 0 {
		user, err := data.UpdateUserGoodreadsAccountInfo(db, user, goodreadsUserID, accessToken, accessSecret)
		return user, existed, err
	} else { // use doesn't exit -> create
		grc := goodreads.NewClient(os.Getenv("GOODREADS_KEY"))
		grUser, err := grc.UserShow(goodreadsUserID)
		if err != nil {
			log.WithError(err).Error("couldn't fetch user from goodreads")
			return models.User{}, false, err
		}
		user, err := data.FindOrCreateUserFromGoodreads(db, goodreadsUserID, accessToken, accessSecret, grUser.Name, grUser.Link, grUser.Website, grUser.ReviewCount)
		if err != nil {
			log.WithError(err).Error("failed to create new user")
			return models.User{}, false, err
		}
		if user.Username == "" {
			user, err = data.UpdateUserUsername(db, user)
			if err != nil {
				log.WithError(err).Error("failed to update username")
				return models.User{}, false, err
			}
		}
		return user, false, err
	}
}
