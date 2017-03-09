# Authorization in Eve
Maintained by Concur Platform R&D <platform-engineering@concur.com> 

Users require to be authorized when communicating to Eve's APIs and operating project resources via Eve. The authorization are placed at both [API](#apis) level and [Resource](#resources) level. We will explain both levels in more detail at respective sections on this document.

User can be a real end user (a human), or a chat-bot (walle), or eve's system user (terraform). User has to be authenticated before obtaining authorization permissions.

Eve's APIs are reachable REST endpoints which perform certain operation on resource. Without permission, user's request to the API will be denied.

Resources in Eve are Quoin, Quoin Archive, Infrastructure, Provider as of now. Every resource on Eve has its owner. By default, owner is the creator of the resource. Owner can only be modified (ownership transfer) by eve admin users' request.

*We might expend "ownership" to a team/organization in the future design.

[Requirements for Authorization]
---
1. When user is authenticated, authentication provider (OKTA, or OAuth service) will provide user's basic information including unique identifier, organization, teams to Eve

2. When resource is being requested to be accessed, Eve will check authorization based on:
	- If user has the proper permission (Read/Write/Execute) on this resource
    - If user belongs to an organization which has the proper permission (Read/Write/Execute) on this resource
    - If user belongs to a team which has the proper permission (Read/Write/Execute) on this resource

3. Resource on Eve are compossible. Permission to access a specific resource doesn't guarantee the same level of permission to its referenced resources. An example will be user's access to read an infrastructure doesn't automatically applied to user's access to that infrastructure's quoin or quoin archive objects. When user doesn't have permission to the referenced resource, the whole request might be denied

4. User will require to have permission to access Eve API endpoint with only exception of Health endpoint

5. System user (Terraform) will have full access on all of APIs, and all of resources

## Group Level

### User (Individaul)
 
### Organization (Group)

### Team (Group)

### Public (Everyone)

## Resources

### Quoin

1. Read Quoin
2. Write/Modify Quoin
3. Execute/Use Quoin

### Quoin Archive 

1. Read Quoin Archive
2. Write/Modify Quoin Archive
3. Execute/Use Quoin Archive

### Infrastructure

1. Read Infrastructure
2. Write/Modify Infrastructure
3. Execute/Use Infrastructure

### Provider

1. Read Provider
2. Write/Modify Provider
3. Execute/Use Provider

## APIs

### Health API
Everyone can access `GET /health`

### Quoin APIs
- Access to `GET /quoin/:name`
- Access to `POST /quoin`
- Access to `POST /quoin/:name/upload`

### Infrastructure APIs
- Access to `GET /infrastructure/:name`
- Access to `GET /infrastructure/:name/state`
- Access to `POST /infrastructure`
- Access to `POST /infrastructure/:name/state`
- Access to `DELETE /infrastructure/:name`
- Access to `DELETE /infrastructure/:name/state`

### Provider APIs
- Access to `GET /provider/:name`

