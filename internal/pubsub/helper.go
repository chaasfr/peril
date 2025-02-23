package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
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

func ackmessage(ackType AckType, msg *amqp.Delivery) error {
	var err error
	switch ackType {
	case Ack:
		err = msg.Ack(false)
	case NackRequeue:
		err = msg.Nack(false, true)
	case NackDiscard:
		err = msg.Nack(false, false)
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
		amqp.Table{"x-dead-letter-exchange": "peril_dlx"},
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

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false,
		false,
		amqp.Publishing{
			ContentType: "aplication/json",
			Body:        jsonBytes,
		})
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buf bytes.Buffer
	gobEncoder := gob.NewEncoder(&buf)

	if err := gobEncoder.Encode(val); err != nil {
		return err
	}
	return ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false,
		false,
		amqp.Publishing{
			ContentType: "aplication/gob",
			Body:        buf.Bytes(),
		})
}

func SubscribeJson[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	unmarshaller := func(msgBody []byte) (T, error) {
		var obj T
		buffer := bytes.NewBuffer(msgBody)
		decoderJson := json.NewDecoder(buffer)
		err := decoderJson.Decode(&obj)
		return obj, err
	}
	err := Subscribe(conn, exchange, queueName, key, simpleQueueType, handler, unmarshaller)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	unmarshaller := func(msgBody []byte) (T, error) {
		var obj T
		buffer := bytes.NewBuffer(msgBody)
		decoderGob := gob.NewDecoder(buffer)
		err := decoderGob.Decode(&obj)
		return obj, err
	}
	err := Subscribe(conn, exchange, queueName, key, simpleQueueType, handler, unmarshaller)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	amqpChan, amqpQ, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	err = amqpChan.Qos(10, 0, false) //limit prefetch to 10 messages for this chan
	if err != nil {
		return err
	}

	deliveryChan, err := amqpChan.Consume(amqpQ.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() error {
		for msg := range deliveryChan {
			obj, err := unmarshaller(msg.Body)
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
