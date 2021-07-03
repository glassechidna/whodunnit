# whodunnit

Working towards this: https://twitter.com/__steele/status/1410437278489477120. Dumping code now to
validate if it's useful or not before investing any more time in it.

The idea consists of three parts:

1. Recording all AWS API calls relevant to creating access credentials (e.g. AssumeRole, CreateAccessKey, 
   and so on) across an organization into a DynamoDB table. Uses EventBridge rules deployed to every 
   account and region.

2. An API that reads from the aforementioned DynamoDB table. Given an access key ID, it should yield as
   much useful info as possible, e.g. the CloudTrail event for the creation of that key, the CT event for
   the creation of the role that assumed _that_ role (i.e. role chaining) and so on. It should also be
   able to reconstruct the principal's "effective" tags by combining tags from users, roles and role sessions.
   
3. A stream processor that reads raw CloudTrail files from an S3 bucket and writes enriched trails to a 
   different bucket. It could be enriched with the role chain and tag data from the API. As a bonus, the
   enriched trails could be consolidated into a format that actually plays nicely with AWS Athena. 

Bonus fourth part: give me your feedback and ideas. Still deciding if this would be useful.

![diagram](/diagram.png)

## Example output

Input: 

```json
{"AccessKeyId": "ASIATJVZAAAARRT5X7PR"}
```

Output:

```json
{
  "Principal": "arn:aws:iam::226947510000:role/whodunnit-e2e-RoleEverywhere",
  "PrincipalChain": [
    "arn:aws:iam::226947510000:role/whodunnit-e2e-RoleEverywhere",
    "arn:aws:iam::607481580000:role/whodunnit-e2e-Secondary",
    "arn:aws:iam::607481580000:role/whodunnit-e2e-Initial",
    "arn:aws:iam::607481580000:user/asteele"
  ],
  "Tags": {
    "a": "a",
    "b": "c",
    "c": "c",
    "d": "b"
  },
  "TransitiveTags": [],
  "Start": "2021-07-03T23:12:09Z",
  "End": "2021-07-03T23:27:09Z",
  "Trails": [
    {
      "eventVersion": "1.08",
      "userIdentity": {
        "type": "AssumedRole",
        "principalId": "AROAY24FZZZZDAZYLZ5XI:1625353927984705000",
        "arn": "arn:aws:sts::607481580000:assumed-role/whodunnit-e2e-Secondary/1625353927984705000",
        "accountId": "607481580000",
        "accessKeyId": "ASIAY24FZZZZIJJL6PUI",
        "sessionContext": {
          "sessionIssuer": {
            "type": "Role",
            "principalId": "AROAY24FZZZZDAZYLZ5XI",
            "arn": "arn:aws:iam::607481580000:role/whodunnit-e2e-Secondary",
            "accountId": "607481580000",
            "userName": "whodunnit-e2e-Secondary"
          },
          "webIdFederationData": {},
          "attributes": {
            "creationDate": "2021-07-03T23:12:09Z",
            "mfaAuthenticated": "false"
          }
        }
      },
      "eventTime": "2021-07-03T23:12:09Z",
      "eventSource": "sts.amazonaws.com",
      "eventName": "AssumeRole",
      "awsRegion": "global",
      "sourceIPAddress": "101.176.82.246",
      "userAgent": "aws-sdk-go/1.38.69 (go1.16.5; darwin; arm64)",
      "requestParameters": {
        "tags": [
          {
            "value": "c",
            "key": "c"
          },
          {
            "value": "c",
            "key": "b"
          }
        ],
        "roleArn": "arn:aws:iam::226947510000:role/whodunnit-e2e-RoleEverywhere",
        "roleSessionName": "1625353927984487000",
        "durationSeconds": 900
      },
      "responseElements": {
        "packedPolicySize": 2,
        "credentials": {
          "accessKeyId": "ASIATJVZAAAARRT5X7PR",
          "expiration": "Jul 3, 2021 11:27:09 PM",
          "sessionToken": "FwoGZXIvYXd...."
        },
        "assumedRoleUser": {
          "assumedRoleId": "AROATJVZAAAA6Q53XX4JW:1625353927984487000",
          "arn": "arn:aws:sts::226947510000:assumed-role/whodunnit-e2e-RoleEverywhere/1625353927984487000"
        }
      },
      "requestID": "b3862845-019f-4a5b-a94d-db4ffc380f5a",
      "eventID": "b4066adf-d9af-4ad7-ae62-eaaac953e304",
      "readOnly": true,
      "resources": [
        {
          "accountId": "226947510000",
          "type": "AWS::IAM::Role",
          "ARN": "arn:aws:iam::226947510000:role/whodunnit-e2e-RoleEverywhere"
        }
      ],
      "eventType": "AwsApiCall",
      "managementEvent": true,
      "recipientAccountId": "607481580000",
      "sharedEventID": "8dcbd001-f8ef-42a9-8fcf-e2d9c86afcdf",
      "eventCategory": "Management",
      "tlsDetails": {
        "tlsVersion": "TLSv1.2",
        "cipherSuite": "ECDHE-RSA-AES128-SHA",
        "clientProvidedHostHeader": "sts.amazonaws.com"
      }
    },
    {
      "eventVersion": "1.08",
      "userIdentity": {
        "type": "AssumedRole",
        "principalId": "AROAY24FZZZZI4PQIZZHL:1625353927984733000",
        "arn": "arn:aws:sts::607481580000:assumed-role/whodunnit-e2e-Initial/1625353927984733000",
        "accountId": "607481580000",
        "accessKeyId": "ASIAY24FZZZZJJS7IPOU",
        "sessionContext": {
          "sessionIssuer": {
            "type": "Role",
            "principalId": "AROAY24FZZZZI4PQIZZHL",
            "arn": "arn:aws:iam::607481580000:role/whodunnit-e2e-Initial",
            "accountId": "607481580000",
            "userName": "whodunnit-e2e-Initial"
          },
          "webIdFederationData": {},
          "attributes": {
            "creationDate": "2021-07-03T23:12:09Z",
            "mfaAuthenticated": "false"
          }
        }
      },
      "eventTime": "2021-07-03T23:12:09Z",
      "eventSource": "sts.amazonaws.com",
      "eventName": "AssumeRole",
      "awsRegion": "global",
      "sourceIPAddress": "101.176.82.246",
      "userAgent": "aws-sdk-go/1.38.69 (go1.16.5; darwin; arm64)",
      "requestParameters": {
        "tags": [
          {
            "value": "b",
            "key": "b"
          },
          {
            "value": "b",
            "key": "c"
          },
          {
            "value": "b",
            "key": "d"
          }
        ],
        "roleArn": "arn:aws:iam::607481580000:role/whodunnit-e2e-Secondary",
        "roleSessionName": "1625353927984705000",
        "durationSeconds": 900
      },
      "responseElements": {
        "packedPolicySize": 3,
        "credentials": {
          "accessKeyId": "ASIAY24FZZZZIJJL6PUI",
          "expiration": "Jul 3, 2021 11:27:09 PM",
          "sessionToken": "FwoGZXIvYX...."
        },
        "assumedRoleUser": {
          "assumedRoleId": "AROAY24FZZZZDAZYLZ5XI:1625353927984705000",
          "arn": "arn:aws:sts::607481580000:assumed-role/whodunnit-e2e-Secondary/1625353927984705000"
        }
      },
      "requestID": "848bcd52-bc71-401c-b6fd-497c9a992005",
      "eventID": "bd031a4b-eeb3-47ec-bb35-c8372a556c6d",
      "readOnly": true,
      "resources": [
        {
          "accountId": "607481580000",
          "type": "AWS::IAM::Role",
          "ARN": "arn:aws:iam::607481580000:role/whodunnit-e2e-Secondary"
        }
      ],
      "eventType": "AwsApiCall",
      "managementEvent": true,
      "recipientAccountId": "607481580000",
      "eventCategory": "Management",
      "tlsDetails": {
        "tlsVersion": "TLSv1.2",
        "cipherSuite": "ECDHE-RSA-AES128-SHA",
        "clientProvidedHostHeader": "sts.amazonaws.com"
      }
    },
    {
      "eventVersion": "1.08",
      "userIdentity": {
        "type": "IAMUser",
        "principalId": "AIDAJCXW2GXJQGTRBT3DC",
        "arn": "arn:aws:iam::607481580000:user/asteele",
        "accountId": "607481580000",
        "accessKeyId": "AKIAJWMRCSBBKQO4HQCQ",
        "userName": "asteele"
      },
      "eventTime": "2021-07-03T23:12:09Z",
      "eventSource": "sts.amazonaws.com",
      "eventName": "AssumeRole",
      "awsRegion": "global",
      "sourceIPAddress": "101.176.82.246",
      "userAgent": "aws-sdk-go/1.38.69 (go1.16.5; darwin; arm64)",
      "requestParameters": {
        "tags": [
          {
            "value": "a",
            "key": "a"
          },
          {
            "value": "a",
            "key": "b"
          }
        ],
        "roleArn": "arn:aws:iam::607481580000:role/whodunnit-e2e-Initial",
        "roleSessionName": "1625353927984733000",
        "durationSeconds": 900
      },
      "responseElements": {
        "packedPolicySize": 2,
        "credentials": {
          "accessKeyId": "ASIAY24FZZZZJJS7IPOU",
          "expiration": "Jul 3, 2021 11:27:09 PM",
          "sessionToken": "FwoGZXIvYXdzE...."
        },
        "assumedRoleUser": {
          "assumedRoleId": "AROAY24FZZZZI4PQIZZHL:1625353927984733000",
          "arn": "arn:aws:sts::607481580000:assumed-role/whodunnit-e2e-Initial/1625353927984733000"
        }
      },
      "requestID": "5ccbcf5b-a1fa-4e95-9b84-8012732ce3aa",
      "eventID": "bd1f4773-8432-446a-b492-61f4137c4486",
      "readOnly": true,
      "resources": [
        {
          "accountId": "607481580000",
          "type": "AWS::IAM::Role",
          "ARN": "arn:aws:iam::607481580000:role/whodunnit-e2e-Initial"
        }
      ],
      "eventType": "AwsApiCall",
      "managementEvent": true,
      "recipientAccountId": "607481580000",
      "eventCategory": "Management",
      "tlsDetails": {
        "tlsVersion": "TLSv1.2",
        "cipherSuite": "ECDHE-RSA-AES128-SHA",
        "clientProvidedHostHeader": "sts.amazonaws.com"
      }
    }
  ]
}
```

## TODO

### Types of CloudTrail events to handle:

- `sts:GetSessionToken`
- `sts:GetFederationToken`
- `sts:AssumeRoleWithSAML`
- `sts:AssumeRoleWithWebIdentity`
- `sts:AssumeRole` aws service
- `sts:AssumeRole` aws service: service-linked role (is this different?)
- `sts:AssumeRole` same-account
- `sts:AssumeRole` cross-account
- `cloudformation:CreateStack` w/ and w/o execution role
- `cloudformation:UpdateStack` w/ and w/o execution role
- `cloudformation:DeleteStack` w/ and w/o execution role
- stack set apis  
- `iam:CreateAccessKey`
- `iam:UpdateAccessKey`
- `iam:DeleteAccessKey`
- `iam:UpdateUser` an IAM user's username and path can change
- `iam:TagUser`
- `iam:UntagUser`
- `iam:TagRole`
- `iam:UntagRole`
- `iam:TagOpenIDConnectProvider`
- `iam:UntagOpenIDConnectProvider`
- `iam:TagSAMLProvider`
- `iam:UntagSAMLProvider`
