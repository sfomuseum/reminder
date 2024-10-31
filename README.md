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

### add-reminder

```
$> ./bin/add-reminder \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-schedule '0,15,30,45 * * * *' \
	-notify-before 'PT2M' \
	-message 'Hello world'
	
2024/10/31 11:39:11 INFO New reminder added id=1852057750180204544

#> ./bin/add-reminder \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-schedule '0,15,30,45 * * * *' \
	-notify-before 'PT2M' \
	-message 'Hello world 2'
	
2024/10/31 11:43:55 INFO New reminder added id=1852058940615954432
```

### list-reminders

```
$> ./bin/list-reminders \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db'

notify_before,message,deliver_to,id,schedule
PT2M,Hello world,,1852057750180204544,"0,15,30,45 * * * *"
PT2M,Hello world 2,,1852058940615954432,"0,15,30,45 * * * *"
```

### remove-reminders

```
$> ./bin/remove-reminders \
	-reminders-database-uri 'sql://sqlite3?dsn=reminders.db' \
	-id 1852058940615954432
	
2024/10/31 11:45:13 INFO Reminder has been removed id=1852058940615954432
```

### process-reminders

#### Example

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

