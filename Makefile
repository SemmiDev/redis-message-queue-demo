run:
	go run cmd/main.go

get:
	hey -n 10000 -c 100 http://localhost:8080/health