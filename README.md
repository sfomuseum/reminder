# reminder

An simple and opinionate tool for managing and dispatching reminders.

## Reminders

```
type Reminder struct {
	// A unique identifier for the reminder.
	Id int64 `json:"id"`
	// A valid cron expression (that can be parsed by adhocore/gronx) for the scheduled event.
	Schedule string `json:"schedule"`
	// An ISO8601 duration string indicating the amount of time before the scheduled event is due to start sending reminders.
	NotifyBefore string `json:"notify_before"`
	// The message body of the reminder.
	Message string `json:"message"`
	// The address where the reminder should be delivered to.
	DeliverTo string `json:deliver_to"`
	// The address where the reminder should be delivered from.	
	DeliverFrom string `json:"deliver_from"`	
}
```

Reminders are defined as a cron expression for when a reminder is due and an ISO8601 duration string indicating how soon _before_ that due date reminders should be dispatched (for example using the `cmd/process-reminders` tool descibed below).

## Databases

Reminders are stored in any database that implements the `database.RemindersDatabase` interface:

```
type RemindersDatabase interface {
	AddReminder(context.Context, *reminder.Reminder) error
	RemoveReminder(context.Context, *reminder.Reminder) error
	Reminders(context.Context) iter.Seq2[*reminder.Reminder, error]
	Close() error
}
```

The following database implementations are provided by default:

### CSV

Read reminders from a CSV on the local disk. CSV database URIs take the form of:

```
csv:///path/to/file.csv
```

_As of this writing the CSV database implementation does not support adding or removing reminders._

### Docstore

Read and write reminders using anything that supports the [gocloud.dev/docstore](https://pkg.go.dev/gocloud.dev/docstore) interfaces, for example DynamoDB. Docstore database URIs take the form of: 

```
{DOCSTORE}://
```

For example:

```
awsdynamodb://{TABLE_NAME}?partition_key=Id&allow_scans=true&region={REGION}&credentials={CREDENTIALS}
```

Docstore tables needs to be set up manually. See [schema/dynamodb](schema/dynamodb) and [cmd/create-dynamodb-tables](cmd/create-dynamodb-tables) for details.

### SQL

Read and write reminders using anything that supports the [database/sql.DB](#) interface, for example SQLite. SQL database URIs take the form of:

```
sql://{ENGINE}?dsn={DSN}
```

For example:

```
sql://sqlite3?dsn=reminders.db
```

Support for the [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) package is enabled by default. Database tables need to be set up manually. See [schema/sqlite](schema/sqlite) for details.

## Dispatching reminders

Reminders are dispatched (delivered) using the [sfomuseum/go-messenger](https://github.com/sfomuseum/go-messenger) package. The following messenger agents are available by default:

### Beeep

Dispatch reminders to a desktop notification. Beeep messenger agent URIs take the form of:

```
beeep://
```

### Email

Dispatch reminders to one or more email providers. The following email providers are available by default:

#### SES

Dispatch email reminders using the AWS Simple Email Service (SES). SES email messenger agent URIs take the form of:

```
email-ses://?region={REGION}&credentials={CREDENTIALS}
```

_Note: `{CREDENTIALS}` is expected to be a valid [aaronland/go-aws-auth](https://github.com/aaronland/go-aws-auth?tab=readme-ov-file#credentials) credentials string._

#### SMTP

Dispatch email reminders using a SMTP server. SMTP email messenger agent URIs take the form of:

```
email-smtp://?host={HOST}&port={PORT}&username={USERNAME]&password={PASSWORD}
```

### Slack

Dispatch reminders to a Slack channel. Slack messenger agent URIs take the form of:

```
slack://?webhook={SLACK_WEBHOOK_URL}
```

The Slack channel the reminder should be sent to is expected to be defined in the `Reminder` instance's `To` property.

### Stdout

Dispatch reminders to `STDOUT`. Stdout messenger agent URIs take the form of:

```
stdout://
```

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/add-reminder cmd/add-reminder/main.go
go build -mod vendor -ldflags="-s -w" -o bin/remove-reminders cmd/remove-reminders/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-reminders cmd/list-reminders/main.go
go build -mod vendor -ldflags="-s -w" -o bin/process-reminders cmd/process-reminders/main.go
```

### add-reminder

```
> ./bin/add-reminder -h
  -deliver-from string
    	The address where the reminder should be delivered from.
  -deliver-to string
    	The address where the reminder should be delivered to.
  -message string
    	The message body of the reminder.
  -notify-before string
    	An ISO8601 duration string indicating the amount of time before the scheduled event is due to start sending reminders.
  -reminders-database-uri string
    	A valid sfomuseum/reminder/database.RemindersDatabase URI.
  -schedule string
    	A valid cron expression for the scheduled event. If defined as a 'YYYY-MM-DD' date string those value will be used to generate a new schedule (cron) expression in the form of: 0 0 {DAY} {MONTH} * {YEAR}
  -verbose
    	Enable verbose (debug) logging.
```

For example:

```
$> ./bin/add-reminder \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-schedule '0,15,30,45 * * * *' \
	-notify-before 'PT2M' \
	-message 'Hello world'
	
2024/10/31 11:39:11 INFO New reminder added id=1852057750180204544

$> ./bin/add-reminder \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-schedule '0,15,30,45 * * * *' \
	-notify-before 'PT2M' \
	-message 'Hello world 2'
	
2024/10/31 11:43:55 INFO New reminder added id=1852058940615954432

$> ./bin/add-reminder \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-schedule '2025-03-20' \
	-notify-before P14D \
	-message 'Remember that is 2025' \
	-verbose
	
2024/10/31 14:13:10 DEBUG Verbose logging enabled
2024/10/31 14:13:10 DEBUG Reassign schedule schedule="0 0 20 3 * 2025"
2024/10/31 14:13:10 INFO New reminder added id=1852096497911336960
```

### list-reminders

```
$> ./bin/list-reminders -h
  -reminders-database-uri string
    	A valid sfomuseum/reminder/database.RemindersDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```

For example:

```
$> ./bin/list-reminders \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db'

notify_before,message,deliver_to,id,schedule
PT2M,Hello world,,1852057750180204544,"0,15,30,45 * * * *"
PT2M,Hello world 2,,1852058940615954432,"0,15,30,45 * * * *"
```

### remove-reminders

```
$> ./bin/remove-reminders -h
  -id value
    	One or more reminder IDs to remove
  -reminders-database-uri string
    	A valid sfomuseum/reminder/database.RemindersDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```

For example:

```
$> ./bin/remove-reminders \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-id 1852058940615954432
	
2024/10/31 11:45:13 INFO Reminder has been removed id=1852058940615954432
```

### process-reminders

```
$> ./bin/process-reminders -h
  -frequency string
    	A valid ISO8601 duration string indicating how often to process reminders. Required if -mode daemon. (default "PT1M")
  -messenger-agent-uri value
    	One or more valid sfomuseum/go-messenger.Messenger URIs.
  -mode string
    	Valid options are: cli, daemon, lambda (default "cli")
  -reminders-database-uri string
    	A valid sfomuseum/reminder/database.RemindersDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```

For example:

```
$> ./bin/process-reminders \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-messenger-agent-uri stdout:// \
	-messenger-agent-uri beeep:// \
	-mode daemon \
	-verbose

2024/10/31 11:53:09 DEBUG Verbose logging enabled
2024/10/31 11:54:09 DEBUG Process reminders
2024/10/31 11:54:09 DEBUG Now reminder=1852057750180204544 t=2024-10-31T11:54:09.062-07:00
2024/10/31 11:54:09 DEBUG Next reminder=1852057750180204544 t=2024-10-31T12:00:00.000-07:00
2024/10/31 11:54:09 DEBUG Trigger reminder=1852057750180204544 t=2024-10-31T11:58:00.000-07:00

time passes...

2024/10/31 11:58:09 DEBUG Process reminders
2024/10/31 11:58:09 DEBUG Now reminder=1852057750180204544 t=2024-10-31T11:58:09.066-07:00
2024/10/31 11:58:09 DEBUG Next reminder=1852057750180204544 t=2024-10-31T12:00:00.000-07:00
2024/10/31 11:58:09 DEBUG Trigger reminder=1852057750180204544 t=2024-10-31T11:58:00.000-07:00
2024/10/31 11:58:09 DEBUG Reminder is due reminder=1852057750180204544
Hello world
2024/10/31 11:59:09 DEBUG Process reminders
2024/10/31 11:59:09 DEBUG Now reminder=1852057750180204544 t=2024-10-31T11:59:09.066-07:00
2024/10/31 11:59:09 DEBUG Next reminder=1852057750180204544 t=2024-10-31T12:00:00.000-07:00
2024/10/31 11:59:09 DEBUG Trigger reminder=1852057750180204544 t=2024-10-31T11:58:00.000-07:00
2024/10/31 11:59:09 DEBUG Reminder is due reminder=1852057750180204544
Hello world
```

## See also

* https://github.com/sfomuseum/go-messenger
* https://github.com/adhocore/gronx
* https://github.com/sfomuseum/iso8601duration
