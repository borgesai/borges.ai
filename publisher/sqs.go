package publisher

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/segmentio/ksuid"
	logger "github.com/sirupsen/logrus"
)

var sqsLogger = logger.WithField("module", "sqs")

func makeSQSConn() (*sqs.SQS, error) {
	sess, err := session.NewSession()
	if err != nil {
		logger.WithError(err).Error("error creating session")
		return nil, err
	}
	conn := sqs.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	return conn, nil
}

// queueURL format is this "https://sqs.us-west-2.amazonaws.com/411700958227/borges-sync-goodreads-user.fifo"
func sendMessage(body, groupID, queueURL string) error {
	updateID := ksuid.New().String()
	log := sqsLogger.
		WithField("queue_url", queueURL).
		WithField("group_id", groupID).
		WithField("method", "sendMessage")
	log.Info("start")
	sqsConn, connErr := makeSQSConn()
	if connErr != nil {
		return connErr
	}

	result, err := sqsConn.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         aws.String(groupID),
		MessageDeduplicationId: aws.String(updateID),
		MessageBody:            aws.String(body),
		QueueUrl:               &queueURL,
	})

	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}

	log.WithField("message_id", *result.MessageId).Info("message was sent")
	return nil
}
