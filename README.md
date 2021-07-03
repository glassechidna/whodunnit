
## Dynamo schema



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
