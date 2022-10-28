# OpenSlides Autoupdate Service

The Autoupdate Service is part of the OpenSlides environment. It is a http
endpoint where the clients can connect to get the actual data and also get
updates, when the requested data changes.

IMPORTANT: The data are sent via an open http-connection. All browsers limit the
amount of open http1.1 connections to a domain. For this service to work, the
browser has to connect to the service with http2 and therefore needs https.


## Start

The service needs some secrets to run. You can create them with:

```
mkdir secrets
printf "password" > secrets/postgres_password
printf "my_token_key" > secrets/auth_token_key 
printf "my_cookie_key" > secrets/auth_cookie_key
```

It also needs a running postgres and redis instance. You can start one with:

```
docker run  --network host -e POSTGRES_PASSWORD=password -e POSTGRES_USER=openslides -e POSTGRES_DB=openslides postgres:11
```

and

```
docker run --network host redis
```


### With Golang

```
export SECRETS_PATH=secrets
go build ./cmd/autoupdate
./autoupdate
```


### With Docker

The docker build uses the auth token as default. Either configure it to use the
auth-fake services (see environment variables below) or make sure the service
inside the docker container can connect to the auth service. For example with
the docker argument --network host. The auth-secrets have to given as a file.

```
docker build . --tag openslides-autoupdate
docker run --network host -v $PWD/secrets:/run/secrets openslides-autoupdate
```

It uses the host network to connect to redis and postgres.


### With Auto Restart

To restart the service when ever a source file has shanged, the tool
[CompileDaemon](https://github.com/githubnemo/CompileDaemon) can help.

```
go install github.com/githubnemo/CompileDaemon@latest
CompileDaemon -log-prefix=false -build "go build ./cmd/autoupdate" -command "./autoupdate"
```

The make target `build-dev` creates a docker image that uses this tool. The
environment varialbe `OPENSLIDES_DEVELOPMENT` is used to use default auth keys.

```
make build-dev
docker run --network host --env OPENSLIDES_DEVELOPMENT=true openslides-autoupdate-dev
```


## Test

### With Golang

```
go test ./...
```


### With Make

There is a make target, that creates and runs the docker-test-container:

```
make run-tests
```


## Examples

Curl needs the flag `-N / --no-buffer` or it can happen, that the output is not
printed immediately.


### HTTP requests

When the server is started, clients can listen for keys to do so, they have to
send a keyrequest in the body of the request. An example request is:

`curl -N localhost:9012/system/autoupdate -d '[{"ids": [1], "collection": "user", "fields": {"username": null}}]'`

To see a list of possible json-strings see the file
internal/autoupdate/keysbuilder/keysbuilder_test.go

Keys can also defined with the query parameter `k`:

`curl -N localhost:9012/system/autoupdate?k=user/1/username,user/2/username`

With this query method, it is not possible to request related keys.

A request can have a body and the `k`-query parameter.

After the request is send, the values to the keys are returned as a json-object
without a newline:
```
{"user/1/name":"value","user/2/name":"value"}
```

With the query parameter `single` the server writes the first response and
closes the request immediately. So there are not autoupdates:

`curl -N localhost:9012/system/autoupdate?k=user/1/username&single=1`

With the query parameter `position=XX` it is possible to request the data at a
specific position from the datastore. This implieds `single`:

`curl -N localhost:9012/system/autoupdate?k=user/1/username&position=42`


### Updates via redis

Keys are updated via redis:

`xadd ModifiedFields * user/1/username newName user/1/password newPassword`


### Projector

The data for a projector can be accessed with autoupdate requests. For example use:


```
curl -N localhost:9012/system/autoupdate -d '
[
  {
    "ids": [1],
    "collection": "projector",
    "fields": {
      "current_projection_ids": {
        "type": "relation-list",
        "collection": "projection",
        "fields": {
          "content": null,
          "content_object_id": null,
          "stable": null,
          "type": null,
          "options": null
        }
      }
    }
  }
]'
```

### History Information

To get all history information for an fqid call:

`curl localhost:9012/system/autoupdate/history_information?fqid=motion/42`

It returns a list of all changes to the requested fqid. Each element in the list
is an object like this:

```
{
  "position": 23,
  "user_id": 5,
  "information": "motion was created",
  "timestamp: 1234567
}
```

To get the data at a position, use the normal autoupdate request with the
attribute `position`. See above.


### Internal Restrict FQIDs

The autoupdate service provides an internal route to restrict a list of fqids.

`curl localhost:9012/internal/autoupdate -d '{"user_id":23,"fqids":["user/1","motion/6"]}'`

It returns all fields for the given objects restricted for the given user id.


## Configuration

### Environment variables

The Service uses the following environment variables:

* `AUTOUPDATE_PORT`: Lets the service listen on port 9012. The default is
  `9012`.
* `AUTOUPDATE_HOST`: The device where the service starts. The default is am
  empty string which starts the service on any device.
* `DATASTORE_READER_HOST`: Host of the datastore reader. The default is
  `localhost`.
* `DATASTORE_READER_PORT`: Port of the datastore reader. The default is `9010`.
* `DATASTORE_READER_PROTOCOL`: Protocol of the datastore reader. The default is
  `http`.
* `MESSAGE_BUS_HOST`: Host of the redis server. The default is `localhost`.
* `MESSAGE_BUS_PORT`: Port of the redis server. The default is `6379`.
* `REDIS_TEST_CONN`: Test the redis connection on startup. Disable on the cloud
  if redis needs more time to start then this service. The default is `true`.
* `DATASTORE_DATABASE_HOST`: Postgres Host. The default is `localhost`.
* `DATASTORE_DATABASE_PORT`: Postgres Port. The default is `5432`.
* `DATASTORE_DATABASE_USER`: Postgres User. The default is `openslides`.
* `DATASTORE_DATABASE_NAME`: Postgres Database. The default is `openslides`.
* `VOTE_HOST`: Host of the vote-service. The default is `localhost`.
* `VOTE_PORT`: Port of the vote-service. The default is `9013`.
* `VOTE_PROTOCOL`: Protocol of the vote-service. The default is `http`.
* `AUTH`: Sets the type of the auth service. `fake` (default) or `ticket`.
* `AUTH_HOST`: Host of the auth service. The default is `localhost`.
* `AUTH_PORT`: Port of the auth service. The default is `9004`.
* `AUTH_PROTOCOL`: Protocol of the auth servicer. The default is `http`.
* `OPENSLIDES_DEVELOPMENT`: If set, the service uses the default secrets. The
  default is `false`.
* `METRIC_INTERVAL`: Time in how often the metrics are gathered. Zero disables
  the metrics. The default is `5m`.
* `MAX_PARALLEL_KEYS`: Max keys that are send in one request to the datastore.
  The default is `1000`.
* `DATASTORE_TIMEOUT`: Time until a request to the datastore times out. The
  default is `3s`.
* `SECRETS_PATH`: Path where the secrets are stored. The default is
  `/run/secrets/`.

Valid units for duration values are "ns", "us" (or "µs"), "ms", "s", "m", "h".
One number without a unit is interpreted as seconds. So `3` is the same as `3s`.

### Secrets

Secrets are filenames in the directory `SECRETS_PATH` (default:
`/run/secrets/`). The service only starts if it can find each secret file and
read its content. The default values are only used, if the environment variable
`OPENSLIDES_DEVELOPMENT` is set.

* `auth_token_key`: Key to sign the JWT auth tocken.
* `auth_cookie_key`: Key to sign the JWT auth cookie.
* `postgres_password`: Postgres password.


## Update models.yml

To use a new models.yml update the value in the file `models-version`.
Afterwards call `go generate ./...` to update the generated files.
