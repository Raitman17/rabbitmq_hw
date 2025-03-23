package main

import (
	"bufio"
	"log"
	"time"
    "os"
	"rabbitmq_hw/storage"
	"math/rand"
)

var fileName string = "text.txt"

func main() {
	rand.Seed(time.Now().UnixNano())

	ch := storage.NewRabbitChannel()
	ch.Init()
	defer ch.Close()

	input, err := os.Open(fileName)
	if err != nil {
		log.Panic("Input file not found.")
		return
	}
	defer input.Close()

	in := bufio.NewScanner(input)

	for in.Scan() {
		randomMs := time.Duration(rand.Intn(400))
		time.Sleep(randomMs * time.Millisecond)
		ch.Publish(in.Text())
		log.Printf("Producer sent message %s.\n", in.Text())
	}
	log.Println("Producer finished sending messages.")
}
