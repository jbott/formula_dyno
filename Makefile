formula_dyno: *.go **/*.go
	GOOS=linux GOARCH=arm GOARM=7 go build

.phony: build
build: formula_dyno

.phony: upload
upload: formula_dyno
	scp ./formula_dyno bbb:~/
