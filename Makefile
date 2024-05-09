GOBUILD=go build
GOENV=GOOS=linux GOARCH=amd64


bin/nesd:
	${GOBUILD} -o $@ github.com/nesd/cmd/test
