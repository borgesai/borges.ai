package publisher

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type SyncUserReviewToGoodreadsMsg struct {
	UserID    uint64
	JobID     uint64
	ReviewID  uint64
	BookID    uint64
	EditionID uint64
}

func SyncUserReviewToGoodreads(userID, jobID, reviewID, bookID, editionID uint64) error {
	data := SyncUserReviewToGoodreadsMsg{
		UserID:    userID,
		JobID:     jobID,
		ReviewID:  reviewID,
		BookID:    bookID,
		EditionID: editionID,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("failed to create data")
		return err
	}
	queueURL := "https://sqs.us-west-2.amazonaws.com/411700958227/borges-sync-user-review-to-goodreads-queue.fifo"
	userIDStr := strconv.FormatUint(userID, 10)
	body := string(dataJSON)
	return sendMessage(body, userIDStr, queueURL)
}
