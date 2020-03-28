package jobs

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/publisher"
	"borges.ai/utils"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func SyncGoodreadsBook(user models.User, bookID uint64, goodreadsBookEditionID string) error {
	defer utils.Duration(utils.Track("SyncGoodreadsBook"))

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return err
	}
	defer db.Close()

	job, err := data.CreateJob(db, user, bookID, "book", "sync-book-from-goodreads")
	if err != nil {
		log.WithError(err).Error("failed to create job")
		return err
	}
	err = publisher.SyncGoodreadsBook(user.ID, job.ID, bookID, goodreadsBookEditionID)
	if err != nil {
		log.WithError(err).Error("failed to publish message")
		return err
	}
	return nil
}
func SyncGoodreadsUserAndWait(user models.User, isInitiatedByUser bool) (models.Job, error) {
	defer utils.Duration(utils.Track("SyncGoodreadsUserAndWait"))
	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return models.Job{}, err
	}
	defer db.Close()

	job, err := SyncGoodreadsUser(user, isInitiatedByUser)
	if err != nil {
		if err != nil {
			log.WithError(err).Error("failed to create sync job")
			return models.Job{}, err
		}
	}

	// we wait up to 15 seconds

	//time.Sleep(10 * time.Second)
	deadline := time.Now().Add(15 * time.Second)
	for {
		updatedJob, _ := data.FindJobByID(db, job.ID)
		if (updatedJob.ID > 0 && !updatedJob.FinishDate.IsZero()) || time.Now().After(deadline) {
			return updatedJob, nil
		}
		// sleep for one second
		time.Sleep(1 * time.Second)
	}
	return job, nil
}

func SyncGoodreadsUser(user models.User, isInitiatedByUser bool) (models.Job, error) {
	defer utils.Duration(utils.Track("SyncGoodreadsUser"))

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return models.Job{}, err
	}
	defer db.Close()

	job, err := data.CreateJob(db, user, user.ID, "user", "sync-user-from-goodreads")
	if err != nil {
		log.WithError(err).Error("failed to create job")
		return job, err
	}
	err = publisher.SyncGoodreadsUser(user.ID, job.ID, isInitiatedByUser)
	if err != nil {
		log.WithError(err).Error("failed to publish message")
		return job, err
	}
	return job, err
}

func SyncUserReviewToExternalSystems(user models.User, reviewID, bookID, editionID uint64) error {
	var wg sync.WaitGroup
	var err error
	wg.Add(2)
	go (func() {
		err = SyncUserReviewToGoodreads(user, reviewID, bookID, editionID)
		wg.Done()
	})()
	go (func() {
		err = SyncUserReviewToDropbox(user, reviewID, bookID, editionID)
		wg.Done()
	})()
	wg.Wait()
	return err
}

func SyncUserReviewToGoodreads(user models.User, reviewID, bookID, editionID uint64) error {
	defer utils.Duration(utils.Track("SyncUserReviewToGoodreads"))

	if user.GoodreadsAccessToken == "" || !user.GoodreadsSyncEnabled {
		// do nothing
		log.WithField("userID", user.ID).WithField("reviewID", reviewID).Info("sync to goodreads disabled for this user")
		return nil
	}

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return err
	}
	defer db.Close()
	job, err := data.CreateJob(db, user, user.ID, "user", "sync-user-review-to-goodreads")
	if err != nil {
		log.WithError(err).Error("failed to create job")
		return err
	}
	err = publisher.SyncUserReviewToGoodreads(user.ID, job.ID, reviewID, bookID, editionID)
	if err != nil {
		log.WithError(err).Error("failed to publish message")
		return err
	}
	return nil
}

func SyncUserReviewToDropbox(user models.User, reviewID, bookID, editionID uint64) error {
	defer utils.Duration(utils.Track("SyncUserReviewToDropbox"))

	if user.DropboxID == "" || user.DropboxAccessToken == "" {
		// do nothing
		log.WithField("userID", user.ID).WithField("reviewID", reviewID).Info("sync to dropbox disabled for this user")
		return nil
	}

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return err
	}
	defer db.Close()
	job, err := data.CreateJob(db, user, user.ID, "user", "sync-user-review-to-dropbox")
	if err != nil {
		log.WithError(err).Error("failed to create job")
		return err
	}
	err = publisher.SyncUserReviewToDropbox(user.ID, job.ID, reviewID, bookID, editionID)
	if err != nil {
		log.WithError(err).Error("failed to publish message")
		return err
	}
	return nil
}
