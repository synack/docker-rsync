install:
	godep go install

test:
	GOTEST=1 godep go test -v ./...


.PHONY: test install