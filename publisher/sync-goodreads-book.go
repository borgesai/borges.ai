package publisher

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type SyncGoodreadsBookMsg struct {
	UserID                 uint64
	JobID                  uint64
	BookID                 uint64
	GoodreadsBookEditionID string
}

func SyncGoodreadsBook(userID, jobID, bookID uint64, goodreadsBookEditionID string) error {
	data := SyncGoodreadsBookMsg{
		UserID:                 userID,
		JobID:                  jobID,
		BookID:                 bookID,
		GoodreadsBookEditionID: goodreadsBookEditionID,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("failed to create data")
		return err
	}
	queueURL := "https://sqs.us-west-2.amazonaws.com/411700958227/borges-sync-goodreads-book-queue.fifo"
	userIDStr := strconv.FormatUint(userID, 10)
	body := string(dataJSON)
	return sendMessage(body, userIDStr, queueURL)
}
