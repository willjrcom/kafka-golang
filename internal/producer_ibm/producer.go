package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func Main() {
	// Configuração do produtor
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Criação do produtor
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	notification := Notification{
		From:    User{ID: 1, Name: "Bruno"},
		To:      User{ID: 2, Name: "will"},
		Message: "seja bem vindo",
	}

	notificationJson, _ := json.Marshal(notification)

	// Envia uma mensagem para o tópico "meu-topico"
	msg := &sarama.ProducerMessage{
		Topic: "meu-topico",
		Key:   sarama.StringEncoder("1"),
		Value: sarama.StringEncoder(notificationJson),
	}

	for i := 0; i < 1000; i++ {
		_, _, err := producer.SendMessage(msg)

		if err != nil {
			log.Fatal(err)
		}

	}

	// fmt.Printf("Mensagem enviada para a partição %d com offset %d\n", partition, offset)
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Notification struct {
	From    User   `json:"from"`
	To      User   `json:"to"`
	Message string `json:"message"`
}

func sendKafkaMessage(producer sarama.SyncProducer,
	users []User, ctx *gin.Context) error {
	message := ctx.PostForm("message")

	notification := Notification{
		From:    User{ID: 1, Name: "Bruno"},
		To:      User{ID: 2, Name: "will"},
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "Message-Topic",
		Key:   sarama.StringEncoder(strconv.Itoa(2)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Println("Message sent to partition %d with offset %d", partition, offset)
	return nil
}

func sendMessageHandler(producer sarama.SyncProducer,
	users []User) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := sendKafkaMessage(producer, users, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to send message",
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully!",
		})
	}
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"},
		config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

// func Main() {
// 	users := []User{
// 		{ID: 1, Name: "Emma"},
// 		{ID: 2, Name: "Bruno"},
// 		{ID: 3, Name: "Rick"},
// 		{ID: 4, Name: "Lena"},
// 	}

// 	producer, err := setupProducer()
// 	if err != nil {
// 		log.Fatalf("failed to initialize producer: %v", err)
// 	}
// 	defer producer.Close()

// 	gin.SetMode(gin.ReleaseMode)
// 	router := gin.Default()
// 	router.POST("/send", sendMessageHandler(producer, users))

// 	fmt.Printf("Kafka PRODUCER 📨 started at http://localhost%s\n",
// 		":8080")

// 	if err := router.Run(":8080"); err != nil {
// 		log.Printf("failed to run the server: %v", err)
// 	}
// }
