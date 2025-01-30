# Makefile  
  
.PHONY: build up buld-up down test clean
  
# Build all services  
build:  
	docker-compose -f deploy/docker-compose.yaml build  
  
# Build Docker images  
up:
	docker-compose -f deploy/docker-compose.yaml up -d

build-up:
	docker-compose -f deploy/docker-compose.yaml up --build

# Stop all services  
down:
	docker-compose -f deploy/docker-compose.yaml down --volumes --remove-orphans
  
# Run Integration Tests  
test:
	docker-compose up tests  
  
# Clean build artifacts  
clean:  
	docker-compose down --volumes --remove-orphans  
	docker system prune -f  
