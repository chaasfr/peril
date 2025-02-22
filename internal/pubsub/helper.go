package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int
const (
	Transient SimpleQueueType = iota
	Durable
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
	handler func(T),
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
				handler(obj)
				err = msg.Ack(false)
				if err != nil {
					log.Println(err)
					return err
				}
		}
		return nil
	}()

	return nil
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

