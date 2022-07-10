# Access Control REST API

An authorization API built to utilize Attribute-based access control (ABAC).

# Getting started

In most cases this API is going to be used to provide a means for authorization and management of the policy database. This interface is designed to be used and pointed to a MongoDB instance. The access control enforcement mechanism is designed to use ABAC (attribute based access control) but could be modified to used otherwise. For more information see the access control library we've utilized found here: https://casbin.org/.

At its core the authorization functionality will need a small amount of information to check for access. The information given to the API will be checked against the policy rule database, and a response will be returned. There are two methods provided for checking if an entity has access, JSON and URL parameters can be used as requests.

JSON will be expected in this format:

    {
        "Label" : "admin",
        "Object" : "Registar Office",
        "Action" : "write"
    }

The Label (or subject) field is the entity requesting access. The Object is what the entity is requesting access to. The API needs to know who wants access and to what, then checks this against the policy-rule database to verify if the operation can be authorized. The Action field will be whether the entity wants to perform a read or write. Please note that for some users a policy rule for write and read may be needed individually (two policy entries in the database).

If you are using the access control manager functions you may notice different fields being used when inserting a new policy. The JSON structure used to posted to the API is meant to be more easily human readable. This gets converted to match the keys the access control library uses to perform enforcement. What's important to note here is the evaluator. When adding a new policy only the value in the single quotes need to be updated. For example in "r.sub.Label == 'corrections'", the only change necessary is "r.sub.Label == 'audit'". If an entity (subject) needs access to more than one Object then a policy rule can be created for each relationship with respect to each Action required.

"PType" : "p" where "p" stands for type policy. The rest of the information necessary to post a new policy rule is fairly self explanatory. You can find more about policy definitions here: https://casbin.org/docs/en/syntax-for-models.

# MongoDB

The API is coded to look for a MongoDB Database named "PolicyDatabase" with a Collection called "casbin_rule". Currently the policies are stored in the "casbin_rule" collection.

## Build

go build

## Run 

$ ./go-abac

# REST API

The Access Control REST API is described below.

## Get list of policies

### Request

`GET /auth/manager/queryAll`

### Successful Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    [
        {
            "_id": "5f73311abde4575be3feba0e",
            "ptype": "p",
            "v0": "r.sub.Label == 'admin'",
            "v1": "/corrections",
            "v2": "write"
        }
    ]

### Empty Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    "Policy rule database is empty."

## Check if has access (Enforce)
### Request

`POST /auth/hasAccess`

{
    "Label": "corrections",
    "Object": "corrections",
    "Action": "write"
}

### Successful Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    {
        "code": "00",
        "result": "granted",
        "log": "2020-09-29 08:08:42 [Request: corrections corrections write] true"
    }


## Check if has access (Enforce)
This is a supplemental endpoint which can accept parameters to enforce policy.
### Request

`POST /auth/hasAccess/{sub}/{obj}/{act}`

### Successful Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    {
        "code": "00",
        "result": "granted",
        "log": "2020-09-29 08:09:37 [Request: admin corrections write] true"
    }

## Create a new policy

### Request

`POST /auth/manager/new`

{
    "PType" : "p", 
    "Evaluator" : "r.sub.Label == 'corrections'", 
    "Object" : "/corrections", 
    "Action" : "read" 
}

### Successful Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    {
        "code": "00",
        "result": "successful",
        "log": "2020-09-29 08:06:15"
    }

## Remove a policy

### Request

`DELETE /auth/manager/remove`

    {
        "_ID": "5f71f5131cc6c8aae16d70e2"
    }

### Successful Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    {
        "code": "00",
        "result": "success",
        "log": "2020-09-29 07:51:52"
    }

## Query policy(ies) by label
This endpoint will return one specific policy. 
### Request

`POST /auth/manager/querySingle`

    {
        "Label" : "admin",
        "Object" : "corrections",
        "Action" : "write"
    }

### Successful Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    [
        {
            "_id": "5f73311abde4575be3feba0e",
            "ptype": "p",
            "v0": "r.sub.Label == 'admin'",
            "v1": "/corrections",
            "v2": "write"
        }
    ]

## Query policy(ies) by all attributes
This endpoint will return list containing any policy with matching label (subject).
### Request

`POST /auth/manager/queryMatching`

    {
        "Label" : "corrections"
    }

### Successful Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json; charset=UTF-8

    [
        {
            "_id": "5f733125bde4575be3feba0f",
            "ptype": "p",
            "v0": "r.sub.Label == 'corrections'",
            "v1": "/corrections",
            "v2": "write"
        },
        {
            "_id": "5f733147bde4575be3feba10",
            "ptype": "p",
            "v0": "r.sub.Label == 'corrections'",
            "v1": "/corrections",
            "v2": "read"
        }
    ]
