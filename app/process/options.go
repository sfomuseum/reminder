package process

type RunOptions struct {
	RemindersDatabaseURI string
	MessengerAgentURIs   []string
	Mode                 string
	Verbose              bool
}
