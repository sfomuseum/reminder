# go-messenger

Opinionated Go package providing interfaces and implementations for delivering to-from-subject-body style messages.

## Documentation

Documentation is incomplete at this time.

## Tools

```
$> make cli 
go build -mod vendor -ldflags="-s -w" -o bin/message cmd/message/main.go
```

### message

Command-line tool for delivering messages using or more delivery agents.

```
$> ./bin/message -h
Command-line tool for delivering messages using or more delivery agents.
Usage:
	 ./bin/message [options] message
Valid options are:
  -agent-uri value
    	One or more known sfomuseum/go-messenger.DeliveryAgent URIs. Valid options are: beeep://, email-null://, email-smtp://, email-stdout://, null://, slack://, stdout://
  -from string
    	The name or address of the person or process sending the message.
  -subject string
    	The subject of the message.
  -to value
    	One or more addresses where messages should be delivered.

If the only message input is "-" then data will be read from STDIN.
```

For example:

```
$> echo "HELLO WORLD" | ./bin/message -from aaron -to aaron -agent-uri stdout:// -subject testing -
HELLO WORLD

2024/10/06 17:09:19 INFO Message delivered to=aaron
```