package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/glassechidna/whodunnit/cloudtrail"
	"github.com/glassechidna/whodunnit/dynamo"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func main() {
	sess := session.Must(session.NewSession())
	api := dynamodb.New(sess)
	table := os.Getenv("TABLE_NAME")

	w := newWriter(api, table)
	lambda.Start(w.handle)
}

type writer struct {
	api   dynamodbiface.DynamoDBAPI
	table string
}

func newWriter(api dynamodbiface.DynamoDBAPI, table string) *writer {
	return &writer{api: api, table: table}
}

func (w *writer) handle(ctx context.Context, input *events.CloudWatchEvent) error {
	detail := cloudtrail.EventDetail{}
	err := json.Unmarshal(input.Detail, &detail)
	if err != nil {
		return errors.WithStack(err)
	}

	request := cloudtrail.AssumeRoleRequestParameters{}
	err = json.Unmarshal(detail.RequestParameters, &request)
	if err != nil {
		return errors.WithStack(err)
	}

	response := cloudtrail.AssumeRoleResponseElements{}
	err = json.Unmarshal(detail.ResponseElements, &response)
	if err != nil {
		return errors.WithStack(err)
	}

	roleArn := request.RoleArn
	roleAccountId := strings.Split(roleArn, ":")[4]

	item := dynamo.AssumeRoleItem{
		AccessKeyId: response.Credentials.AccessKeyId,
		AccountId:   input.AccountID,
		Event:       input.Detail,
	}

	_, err = w.api.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName:           &w.table,
		Item:                item.MarshalDynamoItem(),
		ConditionExpression: aws.String("attribute_not_exists(pk) OR accountId = :accountId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":accountId": {S: &roleAccountId},
		},
	})
	if err != nil {
		fmt.Printf("%+v\n", err)
		return errors.WithStack(err)
	}

	return nil
}
