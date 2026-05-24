run-http:
	go build -tags dynamic cmd/http/main.go
	ENVIRONMENT=staging ./main