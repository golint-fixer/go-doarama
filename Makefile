all:
	go test -v ./...
	golint ./...
	go test -cover -race ./...
	go vet ./...
	test -z "$(go fmt -s ./...)"
	go generate ./...
	git diff --exit-code
