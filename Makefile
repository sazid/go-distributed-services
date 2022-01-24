all:
	go build -o bin/logservice app/cmd/logservice/main.go
	go build -o bin/gradingservice app/cmd/gradingservice/main.go
	go build -o bin/registryservice app/cmd/registryservice/main.go
