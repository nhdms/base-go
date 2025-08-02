package app

import (
	"fmt"
	"github.com/nhdms/base-go/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type PublisherInterface interface {
	PublishSimple(exchName string, data []byte) (err error)
	PublishRouting(exchName, routingKey string, data []byte) (err error)
	PublishRoutingPersist(exchName, routingKey string, data []byte) (err error)
	PublishDirectToQueue(queueName string, data []byte) (err error)

	Close() error
}

type Publisher struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	mutex      sync.Mutex
	amqpURI    string
	config     *RabbitMQConfig
	notifyChan chan *amqp.Error // Channel for monitoring channel health
}

func (p *Publisher) PublishRoutingPersist(exchName, routingKey string, data []byte) (err error) {
	if err := p.ensureConnection(); err != nil {
		return err
	}

	return p.channel.Publish(
		exchName,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: 2, // Persistent
			Timestamp:    time.Now(),
		},
	)
}

// NewPublisher creates a new RabbitMQ publisher instance
func NewPublisher() (PublisherInterface, error) {
	config, amqpURI, err := ReadRabbitMQConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read RabbitMQ config: %w", err)
	}

	publisher := &Publisher{
		config:     config,
		amqpURI:    amqpURI,
		notifyChan: make(chan *amqp.Error),
	}

	err = publisher.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return publisher, nil
}

// createChannel creates a new channel and sets up notifications
func (p *Publisher) createChannel() error {
	var err error
	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	// Set up notification for channel closure
	closeChan := make(chan *amqp.Error, 1)
	p.channel.NotifyClose(closeChan)

	// Monitor channel health in a separate goroutine
	go func() {
		select {
		case err := <-closeChan:
			if err != nil {
				p.notifyChan <- err
			}
		}
	}()

	logger.DefaultLogger.Debugw("create new channel!")
	return nil
}

// connect establishes connection to RabbitMQ with retry mechanism
func (p *Publisher) connect() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.conn != nil && !p.conn.IsClosed() && p.channel != nil {
		// Test channel health with a lightweight operation
		err := p.channel.Flow(true) // Flow control check
		if err == nil {
			return nil // Both connection and channel are healthy
		}
	}

	// Close existing connections if any
	if p.channel != nil {
		p.channel.Close()
		p.channel = nil
	}
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}

	var err error
	maxRetries := 3
	retryDelay := time.Second

	for i := 0; i < maxRetries; i++ {
		p.conn, err = amqp.Dial(p.amqpURI)
		if err == nil {
			err = p.createChannel()
			if err == nil {
				return nil
			}
		}

		if i < maxRetries-1 { // Don't sleep on the last iteration
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("failed to connect after %d retries: %w", maxRetries, err)
}

// ensureConnection makes sure there's an active connection and channel before publishing
func (p *Publisher) ensureConnection() error {
	// Check if we need to reconnect
	if p.conn == nil || p.conn.IsClosed() || p.channel == nil {
		return p.connect()
	}

	//Test channel health
	if p.channel.IsClosed() {
		// Channel is dead, try to recreate it
		p.mutex.Lock()
		defer p.mutex.Unlock()

		// Double-check after acquiring lock
		err := p.createChannel()
		if err == nil {
			return nil
		}

		// If channel creation failed or connection is dead, try full reconnect
		return p.connect()
	}

	return nil
}

// PublishSimple publishes a message to an exchange without routing key
func (p *Publisher) PublishSimple(exchName string, data []byte) error {
	if err := p.ensureConnection(); err != nil {
		return err
	}

	return p.channel.Publish(
		exchName, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Timestamp:   time.Now(),
		},
	)
}

// PublishRouting publishes a message to an exchange with a routing key
func (p *Publisher) PublishRouting(exchName, routingKey string, data []byte) error {
	if err := p.ensureConnection(); err != nil {
		return err
	}

	return p.channel.Publish(
		exchName,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Timestamp:   time.Now(),
		},
	)
}

// PublishDirectToQueue publishes a message directly to a queue
func (p *Publisher) PublishDirectToQueue(queueName string, data []byte) error {
	if err := p.ensureConnection(); err != nil {
		return err
	}

	// Ensure queue exists
	//_, err := p.channel.QueueDeclare(
	//	queueName, // name
	//	true,      // durable
	//	false,     // delete when unused
	//	false,     // exclusive
	//	false,     // no-wait
	//	nil,       // arguments
	//)
	//if err != nil {
	//	return fmt.Errorf("failed to declare queue: %w", err)
	//}

	return p.channel.Publish(
		"",        // exchange (empty for direct queue publishing)
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Timestamp:   time.Now(),
		},
	)
}

// Close closes the RabbitMQ connection and channel
func (p *Publisher) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
		p.channel = nil
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		p.conn = nil
	}

	close(p.notifyChan)
	return nil
}
