package queueHandler

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitHandler struct {
	chanel   *amqp091.Channel
	exchange string
	logger   *logrus.Logger
}

func NewRabbitHandler(url string, log *logrus.Logger) *RabbitHandler {
	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	exchange := "users storage files app"
	err = ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &RabbitHandler{ch, exchange, log}
}

func (rh *RabbitHandler) Publish(ev any) {
	body, err := json.Marshal(ev)
	if err != nil {
		rh.logger.Error(err)
	}
	err = rh.chanel.Publish(rh.exchange, "", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		rh.logger.Error(err)
	}
}
