GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/add-reminder cmd/add-reminder/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/remove-reminders cmd/remove-reminders/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-reminders cmd/list-reminders/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/process-reminders cmd/process-reminders/main.go

lambda:
	@make lambda-process-reminders


lambda-process-reminders:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f process-reminders.zip; then rm -f process-reminders.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/process-reminders/main.go
	zip process-reminders.zip bootstrap
	rm -f bootstrap
