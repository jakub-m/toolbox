default: \
	mkdir_bin \
    bin/tscalc \
    bin/pocket \
    bin/uniqcut \
    bin/cloudlogs-json-format \
    bin/wall-of-text \
    bin/diff-of-two \
	bin/yamljson


mkdir_bin:
	mkdir -p bin
bin/tscalc:
	cd tscalc && $(MAKE) bin/tscalc && cd ../bin && ln -s ../tscalc/bin/tscalc
bin/pocket:
	cd pocket && $(MAKE) bin/pocket && cd ../bin && ln -s ../pocket/bin/pocket
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
	cd tscalc && $(MAKE) clean
	cd pocket && $(MAKE) clean
	rm -frv bin/
.phony: clean mkdir_bin
