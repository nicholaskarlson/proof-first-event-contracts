.PHONY: verify test demo fmt

verify: test demo

test:
	go test -count=1 ./...

demo:
	go run ./cmd/eventcontracts demo --out ./out

fmt:
	gofmt -w ./cmd ./contract ./internal
