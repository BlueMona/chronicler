# Chronicler

A lightweight distributed micro-service that offers a fast and flexible solution for logging user actions in high-load environments with realtime data querying.

##Install and run

- Install Go lang pack
- make directory like "go_src" - root for all go code
- make directory go_src/src
- `export GOPATH=./go_src`
- `go get github.com/basho/riak-go-client`
- `go get github.com/sdming/gosnow`
- `go get github.com/gin-gonic/gin`
- `go get github.com/PeerioTechnologies/chronicler`
- `go install github.com/PeerioTechnologies/chronicler/chronicler-rest-server`

## Run service
Service binary will appear as the `$GOPATH/bin/chronicler-rest-server`

Service accepts configuration file in json format such as following example with default config values:
```
{
  "days-to-keep": 30,
  "enable-debug": true,
  "rest-port": 8080,
  "storage": {
    "nodes": ["127.0.0.1:11087"],
    "node-template": {
      "min-connections": 10,
      "max-connections": 30
    },
    "log-bucket": "log-entries",
    "index-bucket": "log-indexes"
  }
}
```

In the `nodes` array riak node ports must point to protobuf riak node endpoint.

## Run service with default settings

- `$GOPATH/bin/chronicler-rest-server`

## Run service with custom settings

- `$GOPATH/bin/chronicler-rest-server -c config.json`

## REST interface

###Store new log record

**Request**
- **Method**: POST
- **URL**: /v1/log/append
- **Form fields**:
	- `id` - [string] logs collection identifier (for example user ID) , 
 	- `level` - [string] severity of the log,
 	- `type` - [string] context data of log record, 
 	- `log` - [string] body of the log record

```
curl --data "id=zipp&level=ERROR&type=LOGIN&log=Login+error+some+beautiful+thursday+morning" http://localhost:8080/v1/log/append

```

**Response**
Status **200** if record is saved successfully, status **400** if error is occured and `error` field in resul object contains description of error
```
{
	"error": null 
}
```

###Fetch logs collection
Log collection for particular identifier are sorted on chronological order and keeping records not older than `days-to-keep` days.

**Request**
- **Method**: POST
- **URL**: /v1/log/get/{id}

**Response**
```
{
	"error": null
	"payload": [
		...
		{
			key: "513946109254131712",
			time: "2015-11-18T07:18:21.103801314+02:00",
			level: "ERROR",
			type: "LOGIN",
			caption: "Login error some beautiful thursday morning"
		},
	]
}
```
Example

```
$ curl http://localhost:8080/v1/log/get/zipp

$ curl http://localhost:8080/v1/log/get/zipp
{
  "error": null,
  "payload": [
    {
      "key": "513946069915754496",
      "time": "2015-11-18T07:18:11.724264733+02:00",
      "level": "ERROR",
      "type": "LOGIN",
      "caption": "Login error some beautiful thursday morning 2"
    },
    {
      "key": "513946109254131712",
      "time": "2015-11-18T07:18:21.103801314+02:00",
      "level": "ERROR",
      "type": "LOGIN",
      "caption": "Login error some beautiful thursday morning 3"
    },
    {
      "key": "513946339051659264",
      "time": "2015-11-18T07:19:15.891890789+02:00",
      "level": "ERROR",
      "type": "LOGIN",
      "caption": "Login error some beautiful thursday morning 4"
    },
    {
      "key": "513946456286650368",
      "time": "2015-11-18T07:19:43.842014144+02:00",
      "level": "ERROR",
      "type": "LOGIN",
      "caption": "Login error some beautiful thursday morning 5"
    }
  ]
}
```


###Check service status
Reports current service status

**Request**
- **Method**: GET
- **URL**: /v1/report

**Response**

If everything is fine:

```
{
  state: "ok",
  description: ""
}
```
Example for service in situation where Riak is down

```
$ curl http://localhost:8080/v1/log/get/zipp

$ curl http://localhost:8080/v1/log/get/zipp
{
  state: "fail",
  description: "Timeline: no ping to Riak instance. LogRecords: no ping to Riak instance."
}
```