package consumerproducer

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func Main() {
	writer := kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "mensagens",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	defer writer.Close()

	writer.WriteMessages(context.Background(), kafka.Message{
		Key:       []byte("key"),
		Value:     []byte("value"),
		Partition: 2,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "my-group",
		Topic:    "mensagens",
		MinBytes: 0,    // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer reader.Close()
	readMessages(reader)
}

func readMessages(reader *kafka.Reader) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		if err = reader.CommitMessages(context.Background(), m); err != nil {
			log.Panic(err)
		}
	}

}
