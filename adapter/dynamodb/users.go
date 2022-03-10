package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/mkaiho/go-lambda-api-sample/usecase"
)

var _ entity.UsersReader = (*usersReader)(nil)
var _ entity.UsersWriter = (*usersWriter)(nil)

type userAttributeValue struct {
	ID    string `dynamodbav:"id"`
	Name  string `dynamodbav:"name"`
	Email string `dynamodbav:"email"`
}

func (a *userAttributeValue) toEntity(idm entity.IDManager) (entity.User, error) {
	id, err := idm.From(a.ID)
	if err != nil {
		return nil, err
	}
	email, err := entity.NewEmail(a.Email)
	if err != nil {
		return nil, err
	}
	return entity.NewUser(id, a.Name, email), nil
}

func convertToUserAttributeValue(user entity.User) *userAttributeValue {
	return &userAttributeValue{
		ID:    user.UserID().Value(),
		Name:  user.Name(),
		Email: user.Email().Value(),
	}
}

/** UsersReader **/
type usersReader struct {
	idm    entity.IDManager
	client DynamoDBClient
	mapper AttributeValueMapper
}

func NewUsersReader(idm entity.IDManager, client DynamoDBClient, mapper AttributeValueMapper) entity.UsersReader {
	return &usersReader{
		idm:    idm,
		client: client,
		mapper: mapper,
	}
}

func (r *usersReader) FindAll() ([]entity.User, error) {
	resp, err := r.client.Scan(dynamodb.ScanInput{
		TableName: aws.String("users"),
	})
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(resp.Items))
	for i, item := range resp.Items {
		attributeValue := userAttributeValue{}
		err := r.mapper.UnmarshalMap(item, &attributeValue)
		if err != nil {
			return nil, err
		}
		user, err := attributeValue.toEntity(r.idm)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}
	return users, nil
}

func (r *usersReader) FindByID(id entity.UserID) (entity.User, error) {
	resp, err := r.client.GetItem(dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id.Value()),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		return nil, usecase.NewErrNotFoundEntity("users", map[string]string{
			"id": id.Value(),
		})
	}

	attributeValue := userAttributeValue{}
	if err := r.mapper.UnmarshalMap(resp.Item, &attributeValue); err != nil {
		return nil, err
	}
	user, err := attributeValue.toEntity(r.idm)
	if err != nil {
		return nil, err
	}

	return user, nil
}

/** UsersWriter **/
type usersWriter struct {
	idm    entity.IDManager
	client DynamoDBClient
	mapper AttributeValueMapper
}

func NewUsersWriter(idm entity.IDManager, client DynamoDBClient, mapper AttributeValueMapper) entity.UsersWriter {
	return &usersWriter{
		idm:    idm,
		client: client,
		mapper: mapper,
	}
}

func (w *usersWriter) Insert(user entity.User) (entity.User, error) {
	item, err := w.mapper.MarshalMap(convertToUserAttributeValue(user))
	if err != nil {
		return nil, err
	}
	if _, ok := item["id"]; !ok {
		item["id"] = &dynamodb.AttributeValue{
			S: aws.String(w.idm.Generate().Value()),
		}
	}
	if err != nil {
		return nil, err
	}
	_, err = w.client.PutItem(dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      item,
	})
	if err != nil {
		return nil, err
	}

	attributeValue := userAttributeValue{}
	if err := w.mapper.UnmarshalMap(item, &attributeValue); err != nil {
		return nil, err
	}
	user, err = attributeValue.toEntity(w.idm)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (w *usersWriter) Delete(id entity.UserID) error {
	_, err := w.client.DeleteItem(dynamodb.DeleteItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id.Value()),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
