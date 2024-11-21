build:
	cargo build
	mkdir -p bin
	ln -sf ../target/debug/comms bin/comms
release:
	cargo build --release
	mkdir -p bin
	ln -sf ../target/release/comms bin/comms
