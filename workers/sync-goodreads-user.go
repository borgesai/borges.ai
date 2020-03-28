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

func ProcessSyncGoodreadsUser(msg *sqs.Message) error {
	defer utils.Duration(utils.Track("ProcessSyncGoodreadsUser"))

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return err
	}
	defer db.Close()
	msgBody := aws.StringValue(msg.Body)
	var msgData publisher.SyncGoodreadsUserMsg
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
	err = goodreads_import.ImportUserDataFromGoodreads(db, user, msgData.IsInitiatedByUser)
	if err != nil {
		log.WithError(err).Error("failed to sync user")
		data.UpdateJobAsFailed(db, user, job, err.Error())
	} else {
		data.UpdateJobAsCompleted(db, user, job)
	}
	// we are going to retry this one
	return err
}
