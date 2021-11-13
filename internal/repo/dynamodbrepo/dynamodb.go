package dynamodbrepo

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/volmedo/almendruco.git/internal/repo"
)

const tableName = "almendruco-chats"

type dynamoDBRepo struct {
	db dynamodbiface.DynamoDBAPI
}

func NewRepo() (repo.Repo, error) {
	s, err := session.NewSession()
	if err != nil {
		return &dynamoDBRepo{}, fmt.Errorf("session creation failed: %s", err)
	}

	db := dynamodb.New(s)

	return &dynamoDBRepo{db: db}, nil
}

func NewRepoWithClient(client dynamodbiface.DynamoDBAPI) repo.Repo {
	return &dynamoDBRepo{db: client}
}

func (dr *dynamoDBRepo) GetChats() ([]repo.Chat, error) {
	out, err := dr.db.Scan(&dynamodb.ScanInput{TableName: aws.String(tableName)})
	if err != nil {
		return []repo.Chat{}, fmt.Errorf("unable to fetch chats from DB: %w", err)
	}

	chats := make([]repo.Chat, 0, *out.Count)
	for _, item := range out.Items {
		chat := repo.Chat{}
		if err := dynamodbattribute.UnmarshalMap(item, &chat); err != nil {
			return []repo.Chat{}, fmt.Errorf("failed to unmarshal record: %w", err)
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (dr *dynamoDBRepo) UpdateLastNotifiedMessage(chatID string, lastNotifiedMessage uint64) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":last": {
				N: aws.String(strconv.FormatUint(lastNotifiedMessage, 10)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(chatID),
			},
		},
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String("SET lastNotifiedMessage = :last"),
	}

	_, err := dr.db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("update last notified message failed: %s", err)
	}

	return nil
}
