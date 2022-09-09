install:
	GOOS=windows go install ../zero/zctl
	GOOS=linux go install ../zero/zctl
	GOOS=drawin  go install ../zero/zctl
	GOARCH=arm64 go install ../zero/zctl
