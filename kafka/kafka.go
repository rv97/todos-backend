package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	bootstrapURL string
	p *kafka.Producer
}

type Topics string

type Message *kafka.Message

const (
	TodoCreated Topics = "todo-created"
	TopicTwo = "xyz"
)

func NewKafka(bootstrapURL string) *Kafka {
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapURL})
	if err != nil {
		log.Panic("Err: Kafka Connection creation failed %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		log.Panic("Err: Something happened in ParseDuration %w", err)
	}

	results, err := a.CreateTopics(ctx, []kafka.TopicSpecification{
		{
			Topic: string(TodoCreated),
			NumPartitions: 1,
			ReplicationFactor: 1,
		},
	}, kafka.SetAdminOperationTimeout(maxDur))

	if err != nil {
		log.Panic("Err: Failed to create topics %w", err)
	}

	for _,result := range results {
		log.Printf("%s\n", result)
	}

	a.Close()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapURL, "go.delivery.reports": false})

	if err != nil {
		log.Panic("Err: Failed to create producer: %w", err)
	}

	log.Printf("Created Producer %v\n", p)


	return &Kafka{
		bootstrapURL,
		p,
	}
}

func (t *Kafka) Produce(topic Topics, key string, value []byte) {
	// go func() {
	// 	for e := range p.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			// The message delivery report, indicating success or
	// 			// permanent failure after retries have been exhausted.
	// 			// Application level retries won't help since the client
	// 			// is already configured to do that.
	// 			m := ev
	// 			if m.TopicPartition.Error != nil {
	// 				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	// 			} else {
	// 				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
	// 					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	// 			}
	// 		case kafka.Error:
	// 			// Generic client instance-level errors, such as
	// 			// broker connection failures, authentication issues, etc.
	// 			//
	// 			// These errors should generally be considered informational
	// 			// as the underlying client will automatically try to
	// 			// recover from any errors encountered, the application
	// 			// does not need to take action on them.
	// 			fmt.Printf("Error: %v\n", ev)
	// 		default:
	// 			fmt.Printf("Ignored event: %s\n", ev)
	// 		}
	// 	}
	// }()

	top := string(topic)
	err := t.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &top, Partition: kafka.PartitionAny},
		Value:          value,
		Headers:        []kafka.Header{{Key: key, Value: value}},
	}, nil)

	if err != nil {
		log.Printf("Failed to produce message: %v\n", err)
	}

	for t.p.Flush(10000) > 0 {
		fmt.Print("Still waiting to flush outstanding messages\n")
	}

	// t.p.Close()
}

func (t *Kafka) Consume(groupId string, topics []Topics, handler func(value Message)) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": t.bootstrapURL,
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"group.id":              groupId,
		"session.timeout.ms":    6000,
		// Start reading from the first message of each assigned
		// partition if there are no previously committed offsets
		// for this group.
		"auto.offset.reset": "earliest",
		// Whether or not we store offsets automatically.
		"enable.auto.offset.store": false,
	})

	if err != nil {
		log.Panic("Err: Consumer not created %w", err)
	}

	log.Printf("Created Consumer %v\n", c)

	defer c.Close()

	consumerTopics := []string{}
	for _,topic := range topics {
		consumerTopics = append(consumerTopics, string(topic))
	}
	err = c.SubscribeTopics(consumerTopics, nil)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		
		default:
			ev := c.Poll(100)

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if e.Headers != nil {
					log.Printf("%% Headers: %v\n", e.Headers)
				}

				// We can store the offsets of the messages manually or let
				// the library do it automatically based on the setting
				// enable.auto.offset.store. Once an offset is stored, the
				// library takes care of periodically committing it to the broker
				// if enable.auto.commit isn't set to false (the default is true).
				// By storing the offsets manually after completely processing
				// each message, we can ensure atleast once processing.
				_, err := c.StoreMessage(e)
				if err != nil {
					log.Printf("Error storing offset after message %s:\n",
						e.TopicPartition)
				}
				handler(e)
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}