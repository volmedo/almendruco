package dynamodbrepo

import (
	"errors"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/volmedo/almendruco.git/repo"
)

const tableName = "almendruco-user-data"

type userRecord struct {
	UserName            string
	Password            string
	LastNotifiedMessage uint64
}

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

func (dr *dynamoDBRepo) GetPassword(userName string) (string, error) {
	ur, err := dr.fetchUserRecord(userName)
	if err != nil {
		return "", err
	}

	return ur.Password, nil
}

func (dr *dynamoDBRepo) fetchUserRecord(userName string) (userRecord, error) {
	res, err := dr.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(userName),
			},
		},
	})
	if err != nil {
		return userRecord{}, fmt.Errorf("error calling GetItem: %s", err)
	}

	if res.Item == nil {
		return userRecord{}, fmt.Errorf("user %s not found", userName)
	}

	ur := userRecord{}

	if err := dynamodbattribute.UnmarshalMap(res.Item, &ur); err != nil {
		return userRecord{}, fmt.Errorf("failed to unmarshal record: %s", err)
	}

	return ur, nil
}

func (dr *dynamoDBRepo) GetLastNotifiedMessage(userName string) (uint64, error) {
	ur, err := dr.fetchUserRecord(userName)
	if err != nil {
		return 0, err
	}

	return ur.LastNotifiedMessage, nil
}

func (dr *dynamoDBRepo) SetLastNotifiedMessage(userName string, id uint64) error {
	return errors.New("not implemented yet")
}
