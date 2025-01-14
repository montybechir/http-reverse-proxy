# Makefile  
  
.PHONY: build docker-build docker-compose-up docker-compose-down test clean  
  
# Build all services  
build:  
	docker-compose build  
  
# Build Docker images  
docker-build:  
	docker-compose build  
  
# Start all services  
docker-compose-up:  
	docker-compose up -d  
  
# Stop all services  
docker-compose-down:  
	docker-compose down --volumes --remove-orphans  
  
# Run Integration Tests  
test:  
	docker-compose up tests  
  
# Clean build artifacts  
clean:  
	docker-compose down --volumes --remove-orphans  
	docker system prune -f  
