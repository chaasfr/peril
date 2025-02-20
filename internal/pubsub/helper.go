package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

func DeclareAndBind(
	conn *amqp.Connection,
	exchange, queueName, key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, *amqp.Queue, error) {
	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	durable := simpleQueueType == 1
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