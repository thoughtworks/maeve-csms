#!/bin/bash

# Change to the directory where your docker-compose.yml is located
cd .

# Start Docker Compose in detached mode
docker-compose up -d

# Check if Docker Compose started successfully
if [ $? -eq 0 ]; then
    echo "Docker Compose started successfully"
else
    echo "Failed to start Docker Compose"
    exit 1
fi

# Wait for a few seconds for services to initialize (adjust as needed)
sleep 10

docker-compose down lb && docker-compose up lb -d

# Define the URL of the health endpoint
HEALTH_ENDPOINT="http://localhost:9410/health"

# Check the health endpoint in a loop until it returns a 200 status code
echo "Waiting for the health endpoint to become available..."
while true; do
    STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_ENDPOINT)
    if [ $STATUS_CODE -eq 200 ]; then
        echo "Health endpoint is available (HTTP 200)"
        break
    else
        echo "Health endpoint is not yet available (HTTP $STATUS_CODE)"
        sleep 5  # Wait for 5 seconds before checking again
    fi
done


# Change to the directory where your docker-compose.yml is located
cd e2e-tests

# Start Docker Compose in detached mode
make up

# Check if Docker Compose started successfully
if [ $? -eq 0 ]; then
    echo "Docker Compose started successfully"
else
    echo "Failed to start Docker Compose"
    exit 1
fi

# Wait for a few seconds for services to initialize (adjust as needed)
sleep 20


# Run your test command
# Replace the command below with your actual test command
echo "Running test command..."
# Example: docker-compose exec <service_name> <test_command>
cd test-driver
go test -v ./... -count=1

# Check the exit status of the test command
if [ $? -eq 0 ]; then
    echo "Test completed successfully"
else
    echo "Test failed"
fi

# Stop Docker Compose
docker-compose down
cd ../.. && docker-compose down