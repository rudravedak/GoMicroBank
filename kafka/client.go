package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Client struct {
	producer sarama.SyncProducer
}

func NewClient(brokers []string) *Client {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Retry/backoff mekanizması
	maxRetries := 30
	retryInterval := 2 * time.Second

	var producer sarama.SyncProducer
	var err error

	for i := 0; i < maxRetries; i++ {
		producer, err = sarama.NewSyncProducer(brokers, config)
		if err == nil {
			break
		}

		log.Printf("Kafka'ya bağlanılamadı (deneme %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("%v sonra tekrar denenecek...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		panic(fmt.Sprintf("Kafka'ya bağlanılamadı (%d deneme sonrası): %v", maxRetries, err))
	}

	log.Println("Kafka'ya başarıyla bağlanıldı!")
	return &Client{
		producer: producer,
	}
}

func (c *Client) Close() error {
	return c.producer.Close()
}

func (c *Client) SendMessage(topic string, message interface{}) error {
	// Mesajı JSON'a çevir
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Kafka mesajını oluştur
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(jsonData),
		Timestamp: time.Now(),
	}

	// Mesajı gönder
	_, _, err = c.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
