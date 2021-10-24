package dynamodbrepo

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/volmedo/almendruco.git/internal/repo"
)

const tableName = "almendruco-target-chats"

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

func (dr *dynamoDBRepo) UpdateLastNotifiedMessage(chatID repo.ChatID, lastNotifiedMessage uint64) error {
	return errors.New("not implemented yet")
}
