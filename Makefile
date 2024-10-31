GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/add-reminder cmd/add-reminder/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/remove-reminders cmd/remove-reminders/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-reminders cmd/list-reminders/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/process-reminders cmd/process-reminders/main.go
