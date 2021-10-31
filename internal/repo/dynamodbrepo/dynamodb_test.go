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
		ID: "chat1",
		Credentials: repo.Credentials{
			User: "user1",
			Pass: "pass1",
		},
		LastNotifiedMessage: 1,
	}

	chat2 = repo.Chat{
		ID: "chat2",
		Credentials: repo.Credentials{
			User: "user2",
			Pass: "pass2",
		},
		LastNotifiedMessage: 2,
	}
)

type dynamoDBClientMock struct {
	dynamodbiface.DynamoDBAPI
}

func (m *dynamoDBClientMock) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	item1, _ := dynamodbattribute.MarshalMap(chat1)
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

	err := dynamoRepo.UpdateLastNotifiedMessage("some_chat", 11)

	assert.Error(t, err)
}
