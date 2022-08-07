up:
	@echo "Starting docker ..."
	docker-compose up -d
	@echo "Docker started"

down:
	@echo "Stopping docker ..."
	docker-compose down
	@echo "Docker stopped"

run:
	@echo "Running application ..."
	go mod tidy
	go run main.go
	@echo "Application started"