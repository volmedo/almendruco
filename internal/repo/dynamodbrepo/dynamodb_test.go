package dynamodbrepo

import (
	"fmt"
	"strconv"
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

func (m *dynamoDBClientMock) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	lastStr := input.ExpressionAttributeValues[":last"].N
	last, err := strconv.ParseUint(*lastStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad lastNotifiedMessage value: %s", *lastStr)
	}

	expectedLast := uint64(11)
	if last != expectedLast {
		return nil, fmt.Errorf("expected last to be %d, but got %d", expectedLast, last)
	}

	key := input.Key["id"].S
	expectedKey := "some_chat"
	if *key != expectedKey {
		return nil, fmt.Errorf("expected key to be %s but got %s", expectedKey, *key)
	}

	updateExp := input.UpdateExpression
	expectedExp := "SET lastNotifiedMessage = :last"
	if *updateExp != expectedExp {
		return nil, fmt.Errorf("expected update exp to be \"%s\" but got \"%s\"", expectedExp, *updateExp)
	}

	return &dynamodb.UpdateItemOutput{}, nil
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

	assert.NoError(t, err)
}
