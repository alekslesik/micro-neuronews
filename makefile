# Run individual services
run-auth:
	go run ./auth-service/main.go

run-user:
	go run ./user-service/main.go

run-content:
	go run ./content-service/main.go

run-comments:
	go run ./comments-service/main.go

run-search:
	go run ./search-service/main.go

run-analytics:
	go run ./analytics-service/main.go

run-notifications:
	go run ./notifications-service/main.go

run-recommendation:
	go run ./recommendation-service/main.go

run-media:
	go run ./media-service/main.go

run-payments:
	go run ./payments-service/main.go

# Build and run everything with Docker Compose
build:
	docker-compose up --build

down:
	docker-compose down
