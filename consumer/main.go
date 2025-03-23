package main

import (
	"log"
	"math/rand"
	"rabbitmq_hw/storage"
	"time"
	"os"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	tableName := os.Getenv("TABLE_NAME")

	ch := storage.NewRabbitChannel()
	ch.Init()
	defer ch.Close()

	db := storage.NewDatabase()
	db.ConnectToDB()
	defer db.Close()

	msgs := ch.GetMessageChan()
	loop:
		for {
			select {
			case msg := <- msgs:
				db.Insert(tableName, string(msg.Body))
				log.Printf("%s wrote the %s into the database\n", tableName, msg.Body)
				msg.Ack(false)
			case <-time.After(3 * time.Second):
				break loop
			}
			randomMs := time.Duration(rand.Intn(800))
			time.Sleep(randomMs * time.Millisecond)
		}

	log.Println("Consumers finished writing messages to database.")
}
