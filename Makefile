default: \
    bin/tscalc \
    bin/uniqcut \
    bin/cloudlogs-json-format \
    bin/wall-of-text \
    bin/diff-of-two \
	bin/yamljson


bin/tscalc: $(shell find cli/tscalc -type f)
	go build -o bin/tscalc tscalc/cli/main.go
bin/uniqcut: $(shell find cli/uniqcut -type f)
	go build -o bin/uniqcut cli/uniqcut/main.go
bin/cloudlogs-json-format: $(shell find cli/cloudlogs-json-format -type f)
	go build -o bin/cloudlogs-json-format cli/cloudlogs-json-format/main.go
bin/wall-of-text: $(shell find cli/wall-of-text -type f)
	go build -o bin/wall-of-text cli/wall-of-text/*
bin/diff-of-two: $(shell find cli/diff-of-two -type f)
	go build -o bin/diff-of-two cli/diff-of-two/*
bin/yamljson: $(shell find cli/yamljson -type f)
	go build -o bin/yamljson cli/yamljson/*
test:
	go test ./...
clean:
	rm -frv bin/
