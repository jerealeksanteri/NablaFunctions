run:
	go run main.go

load:
	curl -X POST -F "code=@(zip_file)" http://localhost:8080/api/load

execute:
	curl http://localhost:8080/api/execute?\functionId\=$(fn)
