#!/bin/bash
docker exec -i -t -u root $(docker ps | grep docker_kafka | cut -d' ' -f1) bash -c "/opt/kafka/bin/kafka-console-consumer.sh --from-beginning --bootstrap-server kafka:9092 --topic=logs & sleep 5 "
