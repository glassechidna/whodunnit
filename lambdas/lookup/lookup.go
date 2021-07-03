package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/glassechidna/whodunnit/cloudtrail"
	"github.com/glassechidna/whodunnit/dynamo"
	"github.com/pkg/errors"
	"os"
	"time"
)

func main() {
	sess := session.Must(session.NewSession())
	api := dynamodb.New(sess)
	table := os.Getenv("TABLE_NAME")

	l := newLookup(api, table)
	lambda.Start(l.handle)
}

type lookup struct {
	api   dynamodbiface.DynamoDBAPI
	table string
}

func newLookup(api dynamodbiface.DynamoDBAPI, table string) *lookup {
	return &lookup{api: api, table: table}
}

func (l *lookup) handle(ctx context.Context, input *Input) (*Output, error) {
	output := &Output{
		Principal:      "",
		PrincipalChain: nil,
		Tags:           nil,
		TransitiveTags: nil,
		Start:          time.Time{},
		End:            time.Time{},
		Trails:         nil,
	}

	get, err := l.api.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: &l.table,
		Key:       dynamo.AssumeRoleKey(input.AccessKeyId),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := dynamo.AssumeRoleItem{}
	err = item.UnmarshalDynamoItem(get.Item)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	detail := item.Detail()
	req := cloudtrail.AssumeRoleRequestParameters{}
	err = json.Unmarshal(detail.RequestParameters, &req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp := cloudtrail.AssumeRoleResponseElements{}
	err = json.Unmarshal(detail.ResponseElements, &resp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	output.Principal = req.RoleArn
	output.Start = detail.EventTime
	output.End, _ = time.Parse(cloudtrail.IamDateFormat, resp.Credentials.Expiration)

	itemStack := []dynamo.AssumeRoleItem{item}
	nextKeyId := detail.UserIdentity.AccessKeyID

	for nextKeyId != "" {
		get, err = l.api.GetItemWithContext(ctx, &dynamodb.GetItemInput{
			TableName: &l.table,
			Key:       dynamo.AssumeRoleKey(nextKeyId),
		})
		if err != nil {
		    return nil, errors.WithStack(err)
		}

		if get.Item == nil {
			break
		}

		item := dynamo.AssumeRoleItem{}
		err = item.UnmarshalDynamoItem(get.Item)
		if err != nil {
		    return nil, errors.WithStack(err)
		}

		itemStack = append(itemStack, item)
		nextKeyId = item.Detail().UserIdentity.AccessKeyID
	}

	for idx := range itemStack {
		detail := itemStack[idx].Detail()
		req := cloudtrail.AssumeRoleRequestParameters{}
		err = json.Unmarshal(detail.RequestParameters, &req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		output.PrincipalChain = append(output.PrincipalChain, req.RoleArn)
		output.Trails = append(output.Trails, itemStack[idx].Event)
	}

	// TODO: the iam user is missing from the principal chain - what to do about that?
	veryFirstPrincipal := itemStack[len(itemStack)-1].Detail().UserIdentity.Arn
	output.PrincipalChain = append(output.PrincipalChain, veryFirstPrincipal)

	transitive := map[string]struct{}{}
	output.Tags = map[string]string{}

	for idx := len(itemStack)-1; idx >= 0; idx-- {
		detail := itemStack[idx].Detail()
		req := cloudtrail.AssumeRoleRequestParameters{}
		err = json.Unmarshal(detail.RequestParameters, &req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, key := range req.TransitiveTagKeys {
			transitive[key] = struct{}{}
		}

		for _, tag := range req.Tags {
			output.Tags[tag.Key] = tag.Value
		}
	}

	output.TransitiveTags = []string{}
	for key := range transitive {
		output.TransitiveTags = append(output.TransitiveTags, key)
	}

	return output, nil
}

type Input struct {
	AccessKeyId string
}

type Output struct {
	Principal      string
	PrincipalChain []string
	Tags           map[string]string
	TransitiveTags []string
	Start          time.Time
	End            time.Time
	Trails         []json.RawMessage
}
