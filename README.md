# Echo Service (Thrift)

This is a simple POC for an echo service using Thrift. The goal is to validate the implementation of a thrift Go service mounted within an HTTP server.

## Running the service

Using this code requires Go 1.17+ and thrift 1.15 installed in the development machine.

### Generating Thrift code
Before attempting to run the code you will need to generate the thrift layers by running the following command within the project root folder:

```
thrift --gen go echoservice.thrift
```

### Running the service
Once those have been generated you can run the service locally with

```
go run ./cmd/service
```

That will start the service on the port `9090`.

## Service endpoints

The services responds to 4 endpoints, each of which is documented below.

#### `GET /thrift`
Responds with `polo` in plain text.

#### `POST /thrift`
Processes a Thrift request for the echo service in the binary protocol.
#### GET /thrift/json
Responds with `polo` in plain text.

#### POST /thrift/json
Processes a Thrift request for the echo service in the JSON protocol.