package main

import (
	"log"
	"math/rand"
	"os"
	"rabbitmq_hw/storage"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	consumerName := os.Getenv("CONSUMER_NAME")

	ch := storage.NewRabbitChannel()
	ch.InitExchange()
	ch.InitQueue(consumerName)
	defer ch.Close()

	db := storage.NewDatabase()
	db.ConnectToDB()
	defer db.Close()

	rdb := storage.NewRedis()
	rdb.ConnectToRedis()
	defer rdb.Close()

	ms := 250

	msgs := ch.GetMessageChan()
loop:
	for {
		select {
		case msg := <-msgs:
			randomMs := time.Duration(rand.Intn(ms))
			time.Sleep(randomMs * time.Millisecond)
			headers := msg.Headers
			idemKey := headers["X-Idempotency-Key"].(string)
			val := rdb.Get(idemKey)
			if val == nil {
				db.Insert(consumerName, string(msg.Body))
				rdb.Set(idemKey, string(msg.Body))
				log.Printf("%s cached message: %s", consumerName, msg.Body)
				ms += 30
			} else {
				db.Insert(consumerName, *val)
				if ms > 30 {
					ms -= 30
				}
			}
			log.Printf("%s wrote the %s into the database\n", consumerName, msg.Body)
		case <-time.After(5 * time.Second):
			break loop
		}
		randomMs := time.Duration(rand.Intn(ms))
		time.Sleep(randomMs * time.Millisecond)
	}

	log.Printf("%s finished writing messages to database.\n", consumerName)
}
