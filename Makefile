build:
	go build -o godelivery ./cmd/godelivery/main.go

up:
	docker-compose up -d

down:
	docker-compose down