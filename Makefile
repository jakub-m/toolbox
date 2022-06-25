gofiles=$(shell find . -name \*.go)

default: \
    bin/tcalc \
    bin/uniqcut

bin/tcalc: $(gofiles)
	go build -o bin/tcalc cli/tcalc/main.go
bin/uniqcut: $(gofiles)
	go build -o bin/uniqcut cli/uniqcut/main.go
clean:
	rm -frv bin/
