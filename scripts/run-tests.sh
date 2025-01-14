#!/bin/bash  
  
set -e  
  
# Function to check service availability  
check_service() {  
  SERVICE=$1  
  PORT=$2  
  echo "Waiting for $SERVICE:$PORT to be available..."  
  while ! nc -z "$SERVICE" "$PORT"; do  
    sleep 1  
  done  
  echo "$SERVICE:$PORT is available."  
}  
  
# Wait for backenda, backendb, and proxy services  
check_service backenda 60408  
check_service backendb 60409  
check_service proxy 8080  
  
# Run the integration tests  
echo "All services are up and running. Executing tests..."  
go test ./tests/...  
