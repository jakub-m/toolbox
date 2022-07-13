gofiles=$(shell find . -name \*.go)

default: \
    bin/tcalc \
    bin/uniqcut \
    bin/cloudlogs-json-format

bin/tscalc: $(gofiles)
	go build -o bin/tscalc cli/tcalc/main.go
bin/uniqcut: $(gofiles)
	go build -o bin/uniqcut cli/uniqcut/main.go
bin/cloudlogs-json-format: $(gofiles)
	go build -o bin/cloudlogs-json-format cli/cloudlogs-json-format/main.go
clean:
	rm -frv bin/
