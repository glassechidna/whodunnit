package dynamo

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/glassechidna/whodunnit/cloudtrail"
	"strings"
)

type AssumeRoleItem struct {
	AccessKeyId string
	AccountId   string
	Event       json.RawMessage
}

func (a *AssumeRoleItem) Detail() cloudtrail.EventDetail {
	detail := cloudtrail.EventDetail{}
	err := json.Unmarshal(a.Event, &detail)
	if err != nil {
	    panic(err)
	}
	return detail
}

func (a *AssumeRoleItem) UnmarshalDynamoItem(m map[string]*dynamodb.AttributeValue) error {
	a.AccessKeyId = strings.TrimPrefix(*m["pk"].S, "key#")
	a.Event = json.RawMessage(*m["event"].S)
	a.AccountId = *m["accountId"].S
	return nil
}

func (a *AssumeRoleItem) MarshalDynamoItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"pk":        {S: aws.String("key#" + a.AccessKeyId)},
		"event":     {S: aws.String(string(a.Event))},
		"accountId": {S: &a.AccountId},
	}
}

func AssumeRoleKey(accessKeyId string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"pk": {S: aws.String("key#" + accessKeyId)},
	}
}
