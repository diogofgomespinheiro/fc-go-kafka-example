package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	deliveryChannel := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("message", "test", producer, nil, deliveryChannel)

	event := <-deliveryChannel
	msg := event.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		fmt.Println("Error sending the message")
	} else {
		fmt.Println("Message sent:", msg.TopicPartition)
	}

	producer.Flush(1000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "go-kafka_kafka:9092",
	}

	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}

	return p
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChannel chan kafka.Event) error {
	message := &kafka.Message{
		Value:          []byte(msg),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
	}

	err := producer.Produce(message, deliveryChannel)
	if err != nil {
		return err
	}

	return nil
}
