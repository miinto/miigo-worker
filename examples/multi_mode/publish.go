package main

import (
"github.com/streadway/amqp"
"strconv"
)

func main () {
	var index int
	size := 2

	conn, _ := amqp.Dial("amqp://miinto:miinto@localhost:5672/")
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	index = 0
	for index < size {
		index++
		ch.Publish("", "go-generic-0-0", false, false, amqp.Publishing{Body: []byte(`{"_type":"Command\\BasicCommand","_data":{"foo":"bar`+strconv.Itoa(index)+`"}}`)})
	}

	index = 0
	for index < size {
		index++
		ch.Publish("", "go-generic-1-0", false, false, amqp.Publishing{Body: []byte(`{"_type":"Command\\ComplexCommand","_data":{"foo":"foocontent","bar":{"foo":"inner_foobar`+strconv.Itoa(index)+`"}}}`)})
	}
}