GODEBUG:=GODEBUG=gocacheverify=1
LOCDIR:=$(PWD)

install_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.0


lint:
	golangci-lint run -v --timeout 3m0s

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./schoolserver ./cmd/schoolserver/main.go

test:
	go test ./... -v

