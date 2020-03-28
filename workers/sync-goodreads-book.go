package workers

import (
	"borges.ai/data"
	"borges.ai/goodreads_import"
	"borges.ai/utils"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"

	"borges.ai/models"

	"borges.ai/publisher"
)

func ProcessSyncGoodreadsBook(msg *sqs.Message) error {
	defer utils.Duration(utils.Track("ProcessSyncGoodreadsBook"))

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return err
	}
	defer db.Close()

	msgBody := aws.StringValue(msg.Body)
	var msgData publisher.SyncGoodreadsBookMsg
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

	err = goodreads_import.SyncBookFromGoodreads(db, user, book, msgData.GoodreadsBookEditionID)
	if err != nil {
		log.WithError(err).Error("failed to sync book")
		data.UpdateJobAsFailed(db, user, job, err.Error())
	} else {
		data.UpdateJobAsCompleted(db, user, job)
	}

	// we are going to retry this one
	return err
}
