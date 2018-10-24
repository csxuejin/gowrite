linux:
	GOOS=linux GOARCH=amd64 go build -o gowrite *.go

mac:
	go build -o gowrite *.go