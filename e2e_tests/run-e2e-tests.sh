#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

# Get the directory where the CSMS is located
DEFAULT_CSMS_DIR="${SCRIPT_DIR}"/..
CSMS_DIR="${1:-$DEFAULT_CSMS_DIR}"

# Define paths relative to the script's location
EVEREST_DIR="$CSMS_DIR/e2e_tests"
TEST_DIR="$CSMS_DIR/e2e_tests/test_driver"


# Function to start Docker Compose
start_docker_compose_for_maeve_csms() {
    cd "$CSMS_DIR"
    (cd config/certificates && make)
    chmod 755 $CSMS_DIR/config/certificates/csms.key
    docker-compose up -d
    if [ $? -eq 0 ]; then
        echo "Docker Compose started successfully"
    else
        echo "Failed to start Docker Compose"
        stop_docker_compose_for_maeve_csms
        exit 1
    fi
}

# Function to start Docker Compose
start_docker_compose_for_everest() {
        source "$SCRIPT_DIR/everest/scripts/copy-csms-cert.sh"
        source "$SCRIPT_DIR/everest/scripts/setup-everest.sh"
        cd "$EVEREST_DIR"
        make up
        if [ $? -ne 0 ]; then
            echo "Failed to start Docker Compose for tests"
            stop_docker_compose_for_everest
            exit 1
        fi

        echo "Waiting for services to initialize..."
        sleep 10
}

# Function to stop Docker Compose
stop_docker_compose_for_everest() {
    cd "$EVEREST_DIR" && docker-compose down
}

stop_docker_compose_for_maeve_csms() {
    cd "$CSMS_DIR" && docker-compose down
}

# Function to check health endpoint
check_health_endpoint() {
    HEALTH_ENDPOINT="http://localhost:9410/health"
    echo "$(date +"%Y-%m-%d %H:%M:%S"):Waiting for the health endpoint to become available..."
    while true; do
        STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_ENDPOINT)
        if [ $STATUS_CODE -eq 200 ]; then
            echo "$(date +"%Y-%m-%d %H:%M:%S"):Health endpoint is available (HTTP 200)"
            break
        else
            echo "$(date +"%Y-%m-%d %H:%M:%S"):Health endpoint is not yet available (HTTP $STATUS_CODE)"
            sleep 5
        fi
    done
}

# Function to run tests
run_tests() {
    echo "Running test command..."
    cd "$TEST_DIR"
    go test --tags=e2e -v ./... -count=1
    TEST_RESULT=$?
    cd "$CSMS_DIR"
    if [ $TEST_RESULT -eq 0 ]; then
        echo "Tests completed successfully"
    else
        echo "Tests failed"
    fi

    stop_docker_compose_for_everest
    stop_docker_compose_for_maeve_csms
}

# Main script execution
start_docker_compose_for_maeve_csms
check_health_endpoint
start_docker_compose_for_everest
run_tests
