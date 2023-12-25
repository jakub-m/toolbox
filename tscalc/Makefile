gofiles=$(shell find . -name \*.go)
bin=bin/tscalc
$(bin): test $(gofiles)
	go build -o $(bin) cli/*.go
test: $(gofiles)
	go test -v ./...
clean:
	rm -fv bin
.phony: test clean

