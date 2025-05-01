package akafka

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func CreateTopicIfNotExists(topic string, servers string) error {
	conn, err := kafka.Dial("tcp", servers)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		if err.Error() != "topic already exists" {
			return err
		}
	}

	return nil
}

func Consume(topics []string, servers string, msgChan chan *kafka.Message) {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{servers},
		GroupID:     "cali-group",
		Topic:       topics[0],
		StartOffset: kafka.FirstOffset,
		MaxWait:     1 * time.Second,
	})

	defer r.Close()
	fmt.Println("Consumindo mensagens do tópico:", topics[0])
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		msgChan <- &msg
	}
}

func Publish(topic string, servers string, msg []byte) error {
	err := CreateTopicIfNotExists(topic, servers)
	if err != nil {
		return err
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{servers},
		Topic:   topic,
	})
	defer w.Close()

	err = w.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
	if err != nil {
		return err
	}

	return nil
}
