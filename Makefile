run:
	go run main.go

submit:
	curl -X POST -F "code=@(zip_file)" http://localhost:8080/api/submit

execute:
	curl http://localhost:8080/api/execute?\fn\=$(fn)
	