package cloudtrail

import (
	"encoding/json"
	"time"
)

const IamDateFormat = "Jan 2, 2006 3:04:05 PM"

type EventDetail struct {
	EventVersion       string          `json:"eventVersion"`
	UserIdentity       UserIdentity    `json:"userIdentity,omitempty"`
	EventTime          time.Time       `json:"eventTime"`
	EventSource        string          `json:"eventSource"`
	EventName          string          `json:"eventName"`
	AwsRegion          string          `json:"awsRegion"`
	SourceIPAddress    string          `json:"sourceIPAddress"`
	UserAgent          string          `json:"userAgent"`
	RequestParameters  json.RawMessage `json:"requestParameters,omitempty"`
	ResponseElements   json.RawMessage `json:"responseElements,omitempty"`
	RequestID          string          `json:"requestID"`
	EventID            string          `json:"eventID"`
	Resources          []Resource      `json:"resources"`
	EventType          string          `json:"eventType"`
	RecipientAccountID string          `json:"recipientAccountId"`
	SharedEventID      string          `json:"sharedEventID"`
}

type UserIdentity struct {
	Type           string         `json:"type"`
	PrincipalID    string         `json:"principalId"`
	AccountID      string         `json:"accountId"`
	InvokedBy      string         `json:"invokedBy"`
	Arn            string         `json:"arn"`
	AccessKeyID    string         `json:"accessKeyId"`
	SessionContext SessionContext `json:"sessionContext"`
}

type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	Arn         string `json:"arn"`
	AccountID   string `json:"accountId"`
	UserName    string `json:"userName"`
}

type WebIDFederationData struct {
}

type Attributes struct {
	MfaAuthenticated string    `json:"mfaAuthenticated"`
	CreationDate     time.Time `json:"creationDate"`
}

type SessionContext struct {
	SessionIssuer       SessionIssuer       `json:"sessionIssuer"`
	WebIDFederationData WebIDFederationData `json:"webIdFederationData"`
	Attributes          Attributes          `json:"attributes"`
	Ec2RoleDelivery     string              `json:"ec2RoleDelivery"`
}

type Resource struct {
	AccountID string `json:"accountId"`
	Type      string `json:"type"`
	ARN       string `json:"ARN"`
}

type AssumeRoleResponseElements struct {
	AssumedRoleUser struct {
		Arn           string `json:"arn"`
		AssumedRoleId string `json:"assumedRoleId"`
	} `json:"assumedRoleUser"`
	Credentials struct {
		AccessKeyId  string `json:"accessKeyId"`
		Expiration   string `json:"expiration"`
		SessionToken string `json:"sessionToken"`
	} `json:"credentials"`
}

type AssumeRoleRequestParameters struct {
	RoleArn           string   `json:"roleArn"`
	RoleSessionName   string   `json:"roleSessionName"`
	Tags              []Tag    `json:"tags"`
	TransitiveTagKeys []string `json:"transitiveTagKeys"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
