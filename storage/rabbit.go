package storage

import (
	"log"
	"fmt"
	"github.com/joho/godotenv"
    "os"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
	"github.com/google/uuid"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type AbsChannel interface {
	InitExchange()
	InitQueue(queueName string)
	Publish(text string)
	Close()
	GetMessageChan() <-chan amqp.Delivery
}

type RabbitChannel struct {
	ch *amqp.Channel
	conn *amqp.Connection
	q amqp.Queue
}

func NewRabbitChannel() AbsChannel {
	return &RabbitChannel{}
}

func (rc *RabbitChannel) InitExchange() {
	connStr := fmt.Sprintf("amqp://%s:%s@host.docker.internal:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_PORT"),
	)
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(connStr)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	failOnError(err, "Failed to connect to RabbitMQ")
	rc.conn = conn

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	rc.ch = ch

	err = ch.ExchangeDeclare(
		"fanout_ex",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
}

func (rc *RabbitChannel) InitQueue(queueName string) {
	q, err := rc.ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	rc.q = q

	err = rc.ch.QueueBind(
		q.Name,
		"",
		"fanout_ex",
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
}

func (rc *RabbitChannel) Publish(text string) {
	err := rc.ch.Publish(
		"fanout_ex",
		"",
		false,
		false,
		amqp.Publishing{
			Headers: amqp.Table{
				"X-Idempotency-Key": uuid.New().String(),
			},
			ContentType: "text/plain",
			Body: []byte(text),
	})
	failOnError(err, "Failed to publish a message")
}

func (rc *RabbitChannel) Close() {
	rc.ch.Close()
	rc.conn.Close()
}

func (rc *RabbitChannel) GetMessageChan() <-chan amqp.Delivery {
	msgs, err := rc.ch.Consume(
		rc.q.Name,
		"",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-prefetch-count": 1,
		},
	)
	failOnError(err, "Failed to get a message")

	return msgs
}