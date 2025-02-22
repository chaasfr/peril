package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int
const (
	Transient SimpleQueueType = iota
	Durable
)

type AckType int
const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

var SimpleQueueTypeName = map[SimpleQueueType]string{
	Transient: "transient",
	Durable:   "durable",
}

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T)  error {
jsonBytes, err := json.Marshal(val)
if err != nil {
	return err
}
ch.PublishWithContext(
	context.Background(),
	exchange, key,
	false,
	false,
	amqp.Publishing{
		ContentType: "aplication/json",
		Body: jsonBytes,
	})

return nil
}

func SubscribeJson[T any](
	conn * amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	amqpChan, amqpQ, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	deliveryChan, err := amqpChan.Consume(amqpQ.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() error {
			for msg := range deliveryChan {
				var obj T
				err := json.Unmarshal(msg.Body, &obj)
				if err != nil {
					log.Println(err)
					return err
				}
				ackType := handler(obj)
				err = ackmessage(ackType, &msg)
				if err != nil {
					log.Println(err)
					return err
				}
		}
		return nil
	}()

	return nil
}

func ackmessage(ackType AckType, msg *amqp.Delivery) error {
	var err error
	switch(ackType){
	case Ack:
		err = msg.Ack(false)
		log.Println("message Acked")
	case NackRequeue:
		err = msg.Nack(false, true)
		log.Println("message Nacked and requeued")
	case NackDiscard:
		err = msg.Nack(false, false)
		log.Println("message Nacked and discarded")
	default:
		err = fmt.Errorf("unknown acktype")
	}
	return err
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange, queueName, key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, *amqp.Queue, error) {
	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	durable := simpleQueueType == Durable
	q, err := amqpChan.QueueDeclare(
		queueName,
		durable,
		!durable,
		!durable,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	err = amqpChan.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, nil, err
	}

	return amqpChan, &q, nil
}

