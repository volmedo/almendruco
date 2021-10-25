package dynamodbrepo

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"

	"github.com/volmedo/almendruco.git/internal/repo"
)

var (
	chat1 = repo.Chat{
		ID: repo.ChatID(1),
		Credentials: repo.Credentials{
			UserName: "user1",
			Password: "pass1",
		},
		LastNotifiedMessage: 1,
	}

	chat2 = repo.Chat{
		ID: repo.ChatID(2),
		Credentials: repo.Credentials{
			UserName: "user2",
			Password: "pass2",
		},
		LastNotifiedMessage: 2,
	}
)

type dynamoDBClientMock struct {
	dynamodbiface.DynamoDBAPI
}

func (m *dynamoDBClientMock) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	item1, err := dynamodbattribute.MarshalMap(chat1)
	if err != nil {
		return &dynamodb.ScanOutput{}, err
	}
	item2, _ := dynamodbattribute.MarshalMap(chat2)

	items := []map[string]*dynamodb.AttributeValue{item1, item2}
	count := int64(len(items))

	return &dynamodb.ScanOutput{
		Count: &count,
		Items: items,
	}, nil
}

func TestGetChats(t *testing.T) {
	mockClient := &dynamoDBClientMock{}
	dynamoRepo := NewRepoWithClient(mockClient)

	chats, err := dynamoRepo.GetChats()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(chats))
	assert.Equal(t, chat1, chats[0])
	assert.Equal(t, chat2, chats[1])
}

func TestUpdateLastNotifiedMessage(t *testing.T) {
	mockClient := &dynamoDBClientMock{}
	dynamoRepo := NewRepoWithClient(mockClient)

	err := dynamoRepo.UpdateLastNotifiedMessage(repo.ChatID(12345678), 11)

	assert.Error(t, err)
}
