install:
	GOOS=windows GOARCH=amd64 go install ../zero/zctl
	GOOS=linux GOARCH=amd64 go install ../zero/zctl
	go install ../zero/zctl