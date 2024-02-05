gofiles=$(shell find . -name \*.go)
bin=bin/pocket
$(bin): test $(gofiles)
	go build -o $(bin) cli/*.go
test: $(gofiles)
	go test ./...
clean:
	rm -frv bin/
.phony: test clean

