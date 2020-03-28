package oauth

import (
	"borges.ai/goodreads"
	"encoding/xml"
	"github.com/repetitive/oauth1"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

func GetUserID(config *oauth1.Config, accessToken, accessSecret string) (int, error) {
	token := oauth1.Token{
		Token:       accessToken,
		TokenSecret: accessSecret,
	}
	cli := oauth1.NewClient(oauth1.NoContext, config, &token)
	resp, err := cli.Get(goodreads.API_ROOT + "/api/auth_user")
	if err != nil {
		log.WithError(err).Error("failed to get good reads user")
		return 0, err
	}
	var r struct {
		User goodreads.AuthUser `xml:"user"`
	}
	err = xml.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.WithError(err).Error("failed to get good reads user")
		return 0, err
	}
	return r.User.ID, nil
}
func CreateReview(config *oauth1.Config, accessToken, accessSecret string, goodreadsBookID, review, finishedAt string, rating, status int) (string, error) {
	token := oauth1.Token{
		Token:       accessToken,
		TokenSecret: accessSecret,
	}
	cli := oauth1.NewClient(oauth1.NoContext, config, &token)
	endpointUrl := "https://www.goodreads.com/review.xml"

	data := url.Values{}
	data.Set("book_id", goodreadsBookID)
	if review != "" {
		data.Set("review[review]", review)
	}
	if finishedAt != "" {
		data.Set("review[read_at]", finishedAt)
	}
	if rating > 0 {
		data.Set("review[rating]", strconv.Itoa(rating))
	}
	if status == 1 {
		data.Set("shelf", goodreads.READ_SHELF)
	}
	if status == 2 {
		data.Set("shelf", goodreads.CURRENTLY_READING_SELF)
	}
	if status == 3 {
		data.Set("shelf", goodreads.TO_READ_SHELF)
	}

	resp, err := cli.PostForm(endpointUrl, data)
	if err != nil {
		return "", err
	}
	log.WithField("url", endpointUrl).WithField("status", resp.Status).Info("got response from goodreads on edit review")

	var r goodreads.Review
	err = xml.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.WithError(err).Info("failed to decode resp")
		return "", err
	}
	return r.ID, err
}

func EditReview(config *oauth1.Config, accessToken, accessSecret string, reviewID string, review, finishedAt string, rating, status int) (string, error) {
	token := oauth1.Token{
		Token:       accessToken,
		TokenSecret: accessSecret,
	}
	cli := oauth1.NewClient(oauth1.NoContext, config, &token)
	endpointUrl := "https://www.goodreads.com/review/" + reviewID + ".xml"

	data := url.Values{}
	data.Set("id", reviewID)
	if review != "" {
		data.Set("review[review]", review)
	}
	if finishedAt != "" {
		data.Set("review[read_at]", finishedAt)
	}
	if rating > 0 {
		data.Set("review[rating]", strconv.Itoa(rating))
	}
	if status == 1 {
		data.Set("finished", "true")
		data.Set("shelf", goodreads.READ_SHELF)
	}
	if status == 2 {
		data.Set("shelf", goodreads.CURRENTLY_READING_SELF)
	}
	if status == 3 {
		data.Set("shelf", goodreads.TO_READ_SHELF)
	}

	resp, err := cli.PostForm(endpointUrl, data)
	if err != nil {
		return "", err
	}
	log.WithField("url", endpointUrl).WithField("status", resp.Status).Info("got response from goodreads on edit review")
	var r goodreads.Review
	err = xml.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.WithError(err).Info("failed to decode resp")
		return "", err
	}
	return r.ID, err
}
