package workers

import (
	"borges.ai/data"
	"github.com/repetitive/dropbox"
	"borges.ai/services"
	"borges.ai/utils"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"

	"borges.ai/models"

	"borges.ai/publisher"
)

func ProcessSyncUserReviewToDropbox(msg *sqs.Message) error {
	defer utils.Duration(utils.Track("ProcessSyncUserReviewToDropbox"))

	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		return err
	}
	defer db.Close()

	msgBody := aws.StringValue(msg.Body)
	var msgData publisher.SyncUserReviewToDropboxMsg
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
	book, err := services.FindUserBookByID(db, user, msgData.BookID)
	if err != nil {
		log.WithError(err).Error("failed to find book")
		return nil
	}

	content := "+++\n"
	content += "title = \"" + book.Title + "\"\n"
	if book.Subtitle != ""{
		content += "subtitle = \"" + book.Subtitle + "\"\n"
	}
	content += "authors = ["
	for idx, author := range book.Authors {
		content += "\"" + author.Name + "\""
		if idx < len(book.Authors) - 1 {
			content += ", "
		}
	}
	content += "]\n"
	if book.Edition.ISBN13 != "" {
		content += "isbn = \"" + book.Edition.ISBN13 + "\"\n"
	} else if book.Edition.ISBN != "" {
		content += "isbn = \"" + book.Edition.ISBN + "\"\n"
	}
	if book.OriginalYear > 0 {
		content += "year = " + strconv.Itoa(book.OriginalYear) + "\n"
	}
	if book.NumPages > 0 {
		content += "pages = " + strconv.Itoa(book.NumPages) + "\n"
	}
	if book.Review.Rating > 0 {
		content += "rating = " + strconv.Itoa(book.Review.Rating) + "\n"
	}
	canonicalURL := ""
	if user.CustomDomain != "" {
		canonicalURL = "https://" + user.CustomDomain + "/b/" + book.Slug + "-" + book.ShortID
	} else {
		canonicalURL = "https://borges.ai/@" + user.Username + "/b/" + book.Slug + "-" + book.ShortID
	}
	content += "url = \"" + canonicalURL + "\"\n"

	if book.Review.Status == 1 && len(book.Readings) > 0 {
		for _, reading := range book.Readings {
			content += "[[readings]]\n"
			content += "started = \"" + utils.FormatDate(reading.StartDate) + "\"\n"
			content += "finished = \"" + utils.FormatDate(reading.FinishDate) + "\"\n"
			content += "notes = \"" + reading.Note + "\"\n"
		}
	}

	content += "+++\n"
	content += book.Review.Content

	r := strings.NewReader(content)
	dir := ""
	if book.Review.Status == 1 {
		dir = "/read/"
	}
	if book.Review.Status == 2 {
		dir = "/reading/"
	}
	if book.Review.Status == 3 {
		dir = "/to-read/"
	}
	// we have short id to prevent names clash
	filePath := dir + book.Slug + "-" + book.ShortID + ".md"
	err = dropbox.ReplaceFile(user.DropboxAccessToken, filePath, r)

	if err != nil {
		log.WithError(err).Error("failed to create/update review in dropbox")
		data.UpdateJobAsFailed(db, user, job, err.Error())
		// we are going to retry this one
		return err
	}
	data.UpdateJobAsCompleted(db, user, job)

	return nil
}
