build:
	$(eval VERSION := $(shell godep go run *.go --version))

	# GOOS=linux GOARCH=amd64 godep go build -o build/docker-rsync.v$(VERSION).linux.x86_64
	# (cd build && tar -cvzf docker-rsync.v$(VERSION).linux.x86_64.tar.gz docker-rsync.v$(VERSION).linux.x86_64)
	# rm build/docker-rsync.v$(VERSION).linux.x86_64

	GOOS=darwin GOARCH=amd64 godep go build -o build/docker-rsync.v$(VERSION).darwin.x86_64
	(cd build && tar -cvzf docker-rsync.v$(VERSION).darwin.x86_64.tar.gz docker-rsync.v$(VERSION).darwin.x86_64)
	rm build/docker-rsync.v$(VERSION).darwin.x86_64

	# TODO returns error: docker/docker/pkg/term/term.go:16: undefined: Termios
	# GOOS=windows GOARCH=amd64 godep go build -o build/docker-rsync.v$(VERSION).windows.x86_64
	# cd build && tar -cvzf docker-rsync.v$(VERSION).windows.x86_64.tar.gz docker-rsync.v$(VERSION).windows.x86_64)
	# rm build/docker-rsync.v$(VERSION).windows.x86_64

clean:
	rm -r build/*

install:
	godep go install

test:
	GOTEST=1 godep go test -v ./...


.PHONY: build clean test install