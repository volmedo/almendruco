package dynamodbrepo

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

type dynamoDBClientMock struct {
	dynamodbiface.DynamoDBAPI
}

func (m *dynamoDBClientMock) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"username":            {S: aws.String("Some User")},
			"password":            {S: aws.String("s0m3p4ss")},
			"lastNotifiedMessage": {N: aws.String(fmt.Sprint(123456))},
		},
	}, nil
}

func TestGetPassword(t *testing.T) {
	mockClient := &dynamoDBClientMock{}
	dynamoRepo := NewRepoWithClient(mockClient)

	pass, err := dynamoRepo.GetPassword("Some User")

	assert.NoError(t, err)
	assert.Equal(t, "s0m3p4ss", pass)
}

func TestGetLastNotifiedMessage(t *testing.T) {
	mockClient := &dynamoDBClientMock{}
	dynamoRepo := NewRepoWithClient(mockClient)

	last, err := dynamoRepo.GetLastNotifiedMessage("Some User")

	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), last)
}
