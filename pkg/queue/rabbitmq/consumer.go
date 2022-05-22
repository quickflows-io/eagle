package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Consumer define consumer for rabbitmq
type Consumer struct {
	addr          string
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	quit          chan struct{}
	exchange      string
	routingKey    string
	queueName     string
	consumerTag   string
	autoDelete    bool                    // Whether to delete automatically
	handler       func(body []byte) error // Business custom consumption function
}

// NewConsumer instance a consumer
func NewConsumer(addr, exchange, queueName string, autoDelete bool, handler func(body []byte) error) *Consumer {
	return &Consumer{
		addr:        addr,
		exchange:    exchange,
		routingKey:  "",
		queueName:   queueName,
		consumerTag: "consumer",
		autoDelete:  autoDelete,
		handler:     handler,
		quit:        make(chan struct{}),
	}
}

// Start start a service
func (c *Consumer) Start() error {
	if err := c.Run(); err != nil {
		return err
	}

	go c.ReConnect()

	return nil
}

// Stop a consumer
func (c *Consumer) Stop() {
	close(c.quit)

	if !c.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := c.channel.Cancel(c.consumerTag, true); err != nil {
			log.Println("rabbitmq consumer - channel cancel failed: ", err)
		}

		if err := c.conn.Close(); err != nil {
			log.Println("rabbitmq consumer - connection close failed: ", err)
		}
	}
}

// Run .
func (c *Consumer) Run() error {
	var err error
	if c.conn, err = OpenConnection(c.addr); err != nil {
		return err
	}

	if c.channel, err = NewChannel(c.conn).Create(); err != nil {
		_ = c.conn.Close()
		return err
	}

	// bind queue
	if _, err = c.channel.QueueDeclare(c.queueName, true, c.autoDelete, false, false, nil); err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	if err = c.channel.QueueBind(c.queueName, c.routingKey, c.exchange, false, nil); err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	var delivery <-chan amqp.Delivery
	// NOTE: autoAck param
	delivery, err = c.channel.Consume(c.queueName, c.consumerTag, true, false, false, false, nil)
	if err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	go c.Handle(delivery)

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))

	return nil
}

// Handle handle data
func (c *Consumer) Handle(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		log.Printf("Consumer received a message: %s in queue: %s", d.Body, c.queueName)
		log.Printf("got %dB delivery: [%v] %q", len(d.Body), d.DeliveryTag, d.Body)
		go func(delivery amqp.Delivery) {
			if err := c.handler(delivery.Body); err == nil {
				// NOTE: If there are now 10 messages, they are all processed concurrently, if the 10th message is processed first,
				// Then the first 9 messages will be confirmed by delivery.Ack(true). When the next 9 messages are processed,
				// Execute delivery.Ack(true) again, which will obviously lead to repeated confirmation of the message
				// Report 406 PRECONDITION_FAILED error, so here is false
				_ = delivery.Ack(false)
			} else {
				// Re-queue, otherwise unacknowledged messages will continue to occupy memory
				_ = delivery.Reject(true)
			}
		}(d)
	}
	log.Println("handle: async deliveries channel closed")
}

// ReConnect .
func (c *Consumer) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				log.Fatalf("rabbitmq consumer - connection NotifyClose: %+v", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				log.Fatalf("rabbitmq consumer - channel NotifyClose: %+v", err)
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			//Turn off SubMsg message delivery
			if err := c.channel.Cancel(c.consumerTag, true); err != nil {
				log.Fatalf("rabbitmq consumer - channel cancel failed: %+v", err)
			}
			if err := c.conn.Close(); err != nil {
				log.Fatalf("rabbitmq consumer - conn cancel failed: %+v", err)
			}
		}

		// IMPORTANT: Notify must be cleared, otherwise dead connections will not be released
		for err := range c.channelNotify {
			println(err)
		}
		for err := range c.connNotify {
			println(err)
		}

	quit:
		for {
			select {
			case <-c.quit:
				return
			default:
				log.Println("rabbitmq consumer - reconnect")

				if err := c.Run(); err != nil {
					log.Printf("rabbitmq consumer - failCheck: %+v", err)

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}

				break quit
			}
		}
	}
}
