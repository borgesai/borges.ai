package workers

import (
	"borges.ai/data"
	goodreadsOAuth "borges.ai/goodreads/oauth"
	"borges.ai/services"
	"borges.ai/utils"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"

	"borges.ai/models"

	"borges.ai/publisher"
)

func ProcessSyncUserReviewToGoodreads(msg *sqs.Message) error {
	defer utils.Duration(utils.Track("ProcessSyncUserReviewToGoodreads"))

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return err
	}
	defer db.Close()

	msgBody := aws.StringValue(msg.Body)
	var msgData publisher.SyncUserReviewToGoodreadsMsg
	err = json.Unmarshal([]byte(msgBody), &msgData)
	if err != nil {
		log.WithError(err).Error("failed to parse message")
		return err
	}

	job, err := data.FindJobByID(db, msgData.JobID)
	if err != nil {
		log.WithError(err).Error("failed to find job")
		return nil
	}
	user, err := data.FindUserByID(db, msgData.UserID)
	if err != nil {
		log.WithError(err).Error("failed to find user")
		return nil
	}
	book, err := data.FindBookByID(db, msgData.BookID)
	if err != nil {
		log.WithError(err).Error("failed to find book")
		return nil
	}
	review, err := data.FindUserReviewByID(db, user, msgData.ReviewID)
	if err != nil {
		log.WithError(err).Error("failed to find review")
		return nil
	}
	readings, err := data.FindUserBookReadings(db, user, msgData.BookID)
	if err != nil {
		log.WithError(err).Error("failed to find readings")
		return nil
	}

	finishedAt := ""
	if len(readings) > 0 && !readings[0].FinishDate.IsZero() {
		finishedAt = utils.FormatDate(readings[0].FinishDate)
	}

	var goodreadsReviewID string
	if review.GoodreadsID == "" && book.BestEditionGoodreadsID != "" {
		goodreadsReviewID, err = goodreadsOAuth.CreateReview(goodreadsOAuth.GetGoodreadsOAuth1Config(), user.GoodreadsAccessToken, user.GoodreadsAccessSecret,
			book.BestEditionGoodreadsID, review.Content, finishedAt, review.Rating, review.Status)
	} else {
		goodreadsReviewID, err = goodreadsOAuth.EditReview(goodreadsOAuth.GetGoodreadsOAuth1Config(), user.GoodreadsAccessToken, user.GoodreadsAccessSecret,
			review.GoodreadsID, review.Content, finishedAt, review.Rating, review.Status)
	}
	if err != nil {
		log.WithError(err).Error("failed to create/update review in goodreads")
		data.UpdateJobAsFailed(db, user, job, err.Error())
		// we are going to retry this one
		return err
	}
	err = services.UpdateReviewAndReadingsWithGoodreadsReviewID(db, user, goodreadsReviewID, review, readings)
	if err != nil {
		log.WithError(err).WithField("goodreads_id", goodreadsReviewID).Error("failed to update data with goodreads is")
		data.UpdateJobAsFailed(db, user, job, err.Error())
		// we are going to retry this one
		return err
	}
	log.WithField("goodreads_id", goodreadsReviewID).Error("updated review with goodreads is")
	data.UpdateJobAsCompleted(db, user, job)

	return nil
}
