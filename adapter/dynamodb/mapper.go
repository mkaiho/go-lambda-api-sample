package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/** AttributeValueMapper **/
type AttributeValueMapper interface {
	MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
	UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error
}
