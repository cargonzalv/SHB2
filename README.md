## sps-header-bidder
Header Bidder service written in Golang. This service provides protected POST API for accepting open-RTB request from integrations like SpringServe and FreeWheel

## Version Requirements:

Go 1.18

Docker 20.10.16+

## To build and start

1. `go build -o sps-hb cmd/app/main.go`
2. `./sps-hb`


## Running Tests

1. make eunit

## Connection pooling for kafka 

```
Add these to the env file
  # kafka.pool_enabled: true
  # kafka.pool_size: "5"
  NOTE: if kafka.pool_enabled is true and no kafka.pool_size is provided, pool_size is considered "1" as default.
```

## tlsDialer
```
kafkaParams.UseTlsDialer is true for dev and test(local) env and false for other envs
```

## I want to...
- [... explore and run Header Bidder service http apis](doc/header-bidder-apis.md)
- [... know more about Header Bidder service](doc/header-bidder-context.md)