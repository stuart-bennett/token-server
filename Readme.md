# Token Server

A simple server written in Go that handles authentication and responds to requests for authentication tokens.

## Running the app
From within your Go workspace:
`go get github.com/stuart-bennett/token-server/server`
`/bin/server`
or
`/bin/server -r redis-host@6379` to use Redis token store

### Running the tests
`go test github.com/stuart-bennett/token-server/server`

### Docker
Docker / Docker Compose support is included with the intention that you don't have to worry about finding a suitable instance of Redis to use with the app.

You'll need to install `docker` and `docker-compose`. Either use your distros package manager to install Docker and Docker Compose or follow the instructions below:

	- `docker`
		- (Windows)[https://docs.docker.com/docker-for-windows/install/]
		- (Mac)[https://docs.docker.com/docker-for-mac/install/]
		- Linux - Either use your distros package manager or get a (binary)[https://docs.docker.com/install/linux/docker-ce/binaries/]
    - `docker-compose` [https://docs.docker.com/compose/install/]

`docker-compose up --build`

...and don't forget `sudo` if you don't run Docker as non-root

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
