#!/bin/bash

# Function to start Docker Compose
start_docker_compose() {
    cd . && docker-compose up -d
    if [ $? -eq 0 ]; then
        echo "Docker Compose started successfully"
    else
        echo "Failed to start Docker Compose"
        exit 1
    fi
}

# Function to stop Docker Compose
stop_docker_compose_for_everest() {
    cd e2e-tests && docker-compose down
}

stop_docker_compose_for_maeve_csms() {
    cd .. && docker-compose down
}

# Function to check health endpoint
check_health_endpoint() {
    HEALTH_ENDPOINT="http://localhost:9410/health"
    echo "Waiting for the health endpoint to become available..."
    while true; do
        STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_ENDPOINT)
        if [ $STATUS_CODE -eq 200 ]; then
            echo "Health endpoint is available (HTTP 200)"
            break
        else
            echo "Health endpoint is not yet available (HTTP $STATUS_CODE)"
            sleep 5
        fi
    done
}

# Function to run tests
run_tests() {
    cd e2e-tests
    make up
    if [ $? -ne 0 ]; then
        echo "Failed to start Docker Compose for tests"
        stop_docker_compose_for_everest
        exit 1
    fi

    echo "Waiting for services to initialize..."
    sleep 20

    echo "Running test command..."
    cd test-driver
    go test -v ./... -count=1
    TEST_RESULT=$?
    cd ../..

    if [ $TEST_RESULT -eq 0 ]; then
        echo "Tests completed successfully"
    else
        echo "Tests failed"
    fi

    pwd
    stop_docker_compose_for_everest
    pwd
    stop_docker_compose_for_maeve_csms
}

# Main script execution
start_docker_compose
check_health_endpoint
run_tests
