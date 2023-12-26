package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

func Main() {
	// Configuração do consumidor
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Criação do consumidor
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Criação do canal de sinais para interromper o consumo
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "consumer-grupo", config)

	if err != nil {
		log.Fatal(err)
	}

	defer consumerGroup.Close()

	// Inicialização do consumidor em uma goroutine
	go func() {
		for {
			select {
			case err := <-consumerGroup.Errors():
				log.Fatal(err)
			}
		}
	}()

	// Função de processamento de mensagens
	processMessage := &ConsumerGroupHandler{}

	// Início do loop de consumo
	consumerLoop := func() {
		for {
			if err := consumerGroup.Consume(context.Background(), []string{"meu-topico"}, processMessage); err != nil {
				log.Fatal(err)
			}
		}
	}

	// Uso de WaitGroup para esperar pelo sinal de interrupção
	var wg sync.WaitGroup
	wg.Add(1)

	// Início do loop de consumo em uma goroutine
	go func() {
		defer wg.Done()
		consumerLoop()
	}()

	// Aguarda sinal de interrupção
	select {
	case <-signals:
		// Recebeu sinal de interrupção, encerra a execução
		fmt.Println("Recebido sinal de interrupção. Encerrando...")
	}

	// Espera a conclusão das goroutines
	wg.Wait()
}

// ConsumerGroupHandler implementa a interface sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Executa qualquer inicialização necessária aqui
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// Executa qualquer limpeza necessária aqui
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Função de processamento de mensagens
	for mensagem := range claim.Messages() {
		fmt.Printf("Mensagem recebida: chave=%s, valor=%s, partição=%d, offset=%d\n",
			string(mensagem.Key), string(mensagem.Value), mensagem.Partition, mensagem.Offset)
		session.MarkMessage(mensagem, "")
	}

	return nil
}
