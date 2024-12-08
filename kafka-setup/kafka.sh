#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# File paths
COMPOSE_FILE="kafka-compose.yaml"

# Kafka topic details
TOPIC_NAME="verve-streaming"
PARTITIONS=1
REPLICATION_FACTOR=1
BROKER="localhost:9092"

# Functions
start_kafka() {
    echo "Starting Kafka using Docker Compose..."
    docker-compose -f $COMPOSE_FILE up -d
    echo "Kafka started successfully."
}

create_topic() {
    echo "Creating Kafka topic: $TOPIC_NAME..."
    docker exec -it kafka kafka-topics --create \
        --topic $TOPIC_NAME \
        --bootstrap-server $BROKER \
        --partitions $PARTITIONS \
        --replication-factor $REPLICATION_FACTOR
    echo "Topic '$TOPIC_NAME' created successfully."
}

start_and_create() {
    start_kafka
    # Adding a small delay to ensure Kafka is fully started before creating the topic
    sleep 5
    create_topic
}

# Script execution
case "$1" in
    start)
        start_kafka
        ;;
    create-topic)
        create_topic
        ;;
    start_and_create)
        start_and_create
        ;;
    *)
        echo "Usage: $0 {start|create-topic|list-topics}"
        exit 1
        ;;
esac
