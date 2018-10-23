formula_dyno: *.go **/*.go
	@GOOS=linux GOARCH=arm GOARM=7 go build -v

.phony: build
build: formula_dyno

.phony: clean
clean:
	@go clean

.phony: test
test:
	@go test ./...

.phony: upload
upload: formula_dyno
	scp ./formula_dyno bbb:~/
