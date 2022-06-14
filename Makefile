bin=bin/tcalc
gofiles=$(shell find . -name \*.go)
$(bin): $(gofiles)
	go build -o $(bin) main.go
clean:
	rm -fv $(bin)
