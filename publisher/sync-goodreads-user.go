package publisher

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type SyncGoodreadsUserMsg struct {
	UserID uint64
	JobID  uint64
	IsInitiatedByUser bool
}

func SyncGoodreadsUser(userID, jobID uint64, isInitiatedByUser bool) error {
	data := SyncGoodreadsUserMsg{
		UserID: userID,
		JobID:  jobID,
		IsInitiatedByUser: isInitiatedByUser,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("failed to create data")
		return err
	}
	queueURL := "https://sqs.us-west-2.amazonaws.com/411700958227/borges-sync-goodreads-user-queue.fifo"
	userIDStr := strconv.FormatUint(userID, 10)
	body := string(dataJSON)
	return sendMessage(body, userIDStr, queueURL)
}
