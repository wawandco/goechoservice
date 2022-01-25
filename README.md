# Goechoservice

This is a simple POC for an echo service using Thrift. The goal is to validate the implementation of thift in Go and write a simplem client that can consume the echo service.

## Running locally

Using this code requires Go 1.17+ and thrift 1.15 installed in the development machine.

### Generating Thrift code
Before attempting to run the code you will need to generate the thrift layers with:

```
thrift --gen go echoservice.thrift
```

### Running the binaries
Once those have been generated you can run the service locally with
```
go run ./cmd/eghoservice/main.go
```

And the client with

```
go run ./cmd/client/main.go
```



