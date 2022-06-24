gofiles=$(shell find . -name \*.go)

bin/tcalc: $(gofiles)
	go build -o bin/tcalc cli/tcalc/main.go
clean:
	rm -frv bin/
