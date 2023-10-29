default: \
    bin/tscalc \
    bin/uniqcut \
    bin/cloudlogs-json-format \
    bin/wall-of-text \
    bin/diff-of-two \
	bin/yamljson


bin/tscalc: $(shell find tscalc/cli -type f)
	go build -o bin/tscalc tscalc/cli/*
bin/uniqcut: $(shell find uniqcut/cli -type f)
	go build -o bin/uniqcut uniqcut/cli/*
bin/cloudlogs-json-format: $(shell find cloudlogs-json-format/cli -type f)
	go build -o bin/cloudlogs-json-format cloudlogs-json-format/cli/*
bin/wall-of-text: $(shell find wall-of-text/cli -type f)
	go build -o bin/wall-of-text wall-of-text/cli/*
bin/diff-of-two: $(shell find diff-of-two/cli -type f)
	go build -o bin/diff-of-two diff-of-two/cli/*
bin/yamljson: $(shell find yamljson/cli -type f)
	go build -o bin/yamljson yamljson/cli/*
test:
	go test ./...
clean:
	rm -frv bin/
