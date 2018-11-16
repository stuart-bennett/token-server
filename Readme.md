# Token Server

A simple server written in Go that handles authentication and responds to requests for authentication tokens.

## Endpoints

### POST /login

Returns an login token that can be used across other endpoints on this service

#### Request
```
{
	"username": "admin",
	"password": "admin1000"
}
```

#### Response
```
{
	"token": "login-token-will-appear-here"
}
```

### GET /username

#### Request
*Requires Authentication* Be sure to include the following header, substituting in your login token:

`X-Auth-Token: <your-login-token>`

#### Response
```
{
	"username": "The username belonging to the token will appear here"
}
```
