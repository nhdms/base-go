package app

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Consumer interface {
	HandleMessage(msg *message.Message) error
	Init() error
	SetPublisher(p PublisherInterface)
	Close()
	GetName() string
}

// RabbitMQConfig holds the RabbitMQ configuration fields.
type RabbitMQConfig struct {
	Username  string
	Password  string
	Host      string
	Port      int
	VHost     string
	Heartbeat int
}

type RmqExchQueueInfo struct {
	Name    string `json:"name,omitempty" mapstructure:"name"`
	Disable bool   `json:"enable" mapstructure:"disable"`

	Exchange   string `json:"exchange,omitempty" mapstructure:"exchange"`
	Queue      string `json:"queue,omitempty" mapstructure:"queue"`
	RoutingKey string `json:"routing_key,omitempty" mapstructure:"routing_key"`
	Type       string `json:"type,omitempty" mapstructure:"type"`

	AutoDelete bool `json:"auto_delete,omitempty" mapstructure:"auto_delete"`
	Durable    bool `json:"durable,omitempty" mapstructure:"durable"`
	Exclusive  bool `json:"exclusive,omitempty" mapstructure:"exclusive"`
	Qos        int  `json:"qos,omitempty" mapstructure:"qos"`

	WorkerCount        int      `json:"worker_count,omitempty" mapstructure:"worker_count"`
	AdditionalBindings []string `json:"additional_bindings" mapstructure:"additional_bindings"` // format exchange:routing_key

	AmqpConfig amqp.Config `json:"amqp_config" mapstructure:"-"`
}

/*
*
[consumers]
[consumers.task_consumer]
exchange = "task_consumer"
queue = "task_consumer"
routing_key = "task_consumer"
type = "direct"
auto_delete = false
durable = true
exclusive = false
qos = 1
*/

func GetTaskDefinitions(taskNames ...string) map[string]RmqExchQueueInfo {
	resp := make(map[string]RmqExchQueueInfo)
	_, amqpURI, err := ReadRabbitMQConfig()
	if err != nil {
		logger.DefaultLogger.Fatalf("Can not read RabbitMQ config: %v", err)
	}

	for _, taskName := range taskNames {
		sub := viper.Sub("consumers." + taskName)
		if sub == nil {
			logger.DefaultLogger.Warnf("Task %s not found in config", taskName)
			continue
		}

		var info RmqExchQueueInfo
		err := sub.Unmarshal(&info)
		if err != nil {
			logger.DefaultLogger.Warnf("Can not unmarshal task %s config by name %v", err, taskName)
			continue
		}

		if info.Disable {
			continue
		}

		info.AmqpConfig = amqp.Config{
			Connection: amqp.ConnectionConfig{
				AmqpURI:   amqpURI,
				Reconnect: amqp.DefaultReconnectConfig(),
			},
			Marshaler: &amqp.DefaultMarshaler{},
			Exchange: amqp.ExchangeConfig{
				GenerateName: amqp.GenerateExchangeNameConstant(info.Exchange),
				Type:         info.Type,
				Durable:      info.Durable,
				AutoDeleted:  info.AutoDelete,
			},
			Queue: amqp.QueueConfig{
				GenerateName: amqp.GenerateQueueNameConstant(info.Queue),
				Durable:      info.Durable,
				AutoDelete:   info.AutoDelete,
				Exclusive:    info.Exclusive,
			},
			QueueBind: amqp.QueueBindConfig{
				GenerateRoutingKey: func(topic string) string {
					return info.RoutingKey
				},
			},
			Publish: amqp.PublishConfig{
				GenerateRoutingKey: func(topic string) string {
					return info.RoutingKey
				},
			},
			Consume: amqp.ConsumeConfig{
				Qos: amqp.QosConfig{
					PrefetchCount: info.Qos,
				},
			},
			TopologyBuilder: &amqp.DefaultTopologyBuilder{},
		}

		resp[taskName] = info
		logger.DefaultLogger.Infof("Loaded task %s config", taskName)
	}
	return resp
}

/**
 * @params name: task name
 * @params handler: consumer handler
 * consul config template:
	[consumers]
	[consumers.hello_iam_go] # hello_iam_go: task's name, must be matched with the consumer's name
	exchange = "task_consumer" # exchange name
	queue = "task_consumer" # queue name
	routing_key = "task_consumer" # routing key
	type = "direct" # exchange type
	auto_delete = false # whether the queue should be deleted when the consumer is closed
	durable = true # whether the queue should be durable
	exclusive = false # whether the queue should be exclusive (locked to one consumer)
	#disable=true # whether the consumer should be disabled
	qos = 10 # number of unacknowledged messages the consumer will request from the broker
    worker_count = 5 # number of consumer workers (worker count should be less than or equal to the qos)
*/

func StartNewConsumer(handler Consumer) error {
	name := handler.GetName()
	config := GetTaskDefinitions(name)
	if len(config) == 0 {
		return fmt.Errorf("task %s not found in config", name)
	}

	consumerConfig, exists := config[name]
	if !exists {
		return fmt.Errorf("task %s config not found", name)
	}

	err := handler.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize gRPC client for task %s: %w", name, err)
	}

	defer handler.Close()
	// Capture OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	lg := watermill.NewStdLogger(false, false)
	//publisher, err := amqp.NewPublisher(consumerConfig.AmqpConfig, lg)
	//if err != nil {
	//	return fmt.Errorf("failed to create publisher for task %s: %w", name, err)
	//}
	//defer publisher.Close()
	publisher, err := NewPublisher()
	if err != nil {
		return fmt.Errorf("failed to create producer for task %s: %w", name, err)
	}
	defer publisher.Close()

	router, err := message.NewRouter(message.RouterConfig{}, lg)
	if err != nil {
		return fmt.Errorf("failed to create router for task %s: %w", name, err)
	}

	router.AddMiddleware(middleware.Recoverer)

	//publisher.Publish("", message.NewMessage(name, []byte("Hello, Athena!")))
	// init redis....
	handler.SetPublisher(publisher)

	amqpConfig := consumerConfig.AmqpConfig
	consumer, err := amqp.NewSubscriber(amqpConfig, lg)
	if err != nil {
		return fmt.Errorf("failed to create subscriber for task %s: %w", name, err)
	}

	err = initAdditionalBindings(consumer, consumerConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize additional bindings for task %s: %w", name, err)
	}

	workerCount := consumerConfig.WorkerCount
	if workerCount <= 0 {
		workerCount = 1
	}

	for i := 0; i < workerCount; i++ {
		ex := fmt.Sprintf("worker-%v-%d", name, i+1)
		router.AddNoPublisherHandler(
			ex,
			ex,
			consumer,
			handler.HandleMessage,
		)
	}

	// Listen for shutdown signal
	go func() {
		<-signalChan
		logger.DefaultLogger.Infof("Shutdown signal received")
		_ = publisher.Close()
		handler.Close()
		_ = router.Close()
	}()

	err = router.Run(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func initAdditionalBindings(consumer *amqp.Subscriber, config RmqExchQueueInfo) (err error) {
	if len(config.AdditionalBindings) == 0 {
		return
	}

	channel, err := consumer.Connection().Channel()
	if err != nil {
		logger.DefaultLogger.Error("Failed to create channel for additional bindings", err)
		return
	}
	defer channel.Close()
	for _, binding := range config.AdditionalBindings {
		split := strings.Split(binding, ":")
		ex := split[0]
		routingKey := ""
		if len(split) > 1 {
			routingKey = split[1]
		}

		err = channel.QueueBind(config.Queue, routingKey, ex, false, nil)
		if err != nil {
			logger.DefaultLogger.Error("Failed to bind additional queue", err)
			return
		}
		logger.DefaultLogger.Infof("Bound additional queue %s to exchange %s with routing key %s", config.Queue, ex, routingKey)
	}

	return nil
}

func ReadRabbitMQConfig() (*RabbitMQConfig, string, error) {
	config := &RabbitMQConfig{
		Username:  viper.GetString("rabbitmq.username"),
		Password:  viper.GetString("rabbitmq.password"),
		Host:      viper.GetString("rabbitmq.host"),
		Port:      viper.GetInt("rabbitmq.port"),
		VHost:     viper.GetString("rabbitmq.vhost"),
		Heartbeat: viper.GetInt("rabbitmq.heartbeat"),
	}

	// Validate required fields
	if config.Username == "" || config.Password == "" || config.Host == "" || config.Port == 0 {
		return nil, "", fmt.Errorf("missing RabbitMQ configuration fields")
	}

	// Set default heartbeat if not specified
	if config.Heartbeat == 0 {
		config.Heartbeat = 60 // Default heartbeat in seconds
	}

	// Build the AMQP URI with heartbeat
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d/%s?heartbeat=%d",
		config.Username, config.Password, config.Host, config.Port, strings.TrimPrefix(config.VHost, "/"), config.Heartbeat)

	return config, amqpURI, nil
}
