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
}
``

## Databases


## Tools

### process

#### Example

```
$> go ./bin/process \
	-reminders-database-uri csv:///usr/local/sfomuseum/reminder/fixtures/reminders.csv \
	-messenger-agent-uri stdout:// \
	-mode daemon \
	-verbose
	
2024/10/31 10:30:31 DEBUG Verbose logging enabled
2024/10/31 10:31:31 DEBUG Process reminders
2024/10/31 10:31:31 DEBUG Now t=2024-10-31T10:31:31.789-07:00
2024/10/31 10:31:31 DEBUG Next t=2024-10-31T10:45:00.000-07:00
2024/10/31 10:31:31 DEBUG Trigger t=2024-10-31T10:43:00.000-07:00

...time passes

2024/10/31 10:43:31 DEBUG Process reminders
2024/10/31 10:43:31 DEBUG Now t=2024-10-31T10:43:31.790-07:00
2024/10/31 10:43:31 DEBUG Next t=2024-10-31T10:45:00.000-07:00
2024/10/31 10:43:31 DEBUG Trigger t=2024-10-31T10:43:00.000-07:00
hello world
2024/10/31 10:44:31 DEBUG Process reminders
2024/10/31 10:44:31 DEBUG Now t=2024-10-31T10:44:31.790-07:00
2024/10/31 10:44:31 DEBUG Next t=2024-10-31T10:45:00.000-07:00
2024/10/31 10:44:31 DEBUG Trigger t=2024-10-31T10:43:00.000-07:00
hello world
2024/10/31 10:45:31 DEBUG Process reminders
2024/10/31 10:45:31 DEBUG Now t=2024-10-31T10:45:31.791-07:00
2024/10/31 10:45:31 DEBUG Next t=2024-10-31T11:00:00.000-07:00
2024/10/31 10:45:31 DEBUG Trigger t=2024-10-31T10:58:00.000-07:00
```

## See also