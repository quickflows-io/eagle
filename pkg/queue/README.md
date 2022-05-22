## message queue

- RocketMQ
- RabbitMQ
- Kafka
- Nats


## effect

- System decoupling
- Asynchronous processing
- Shaving peaks and filling valleys

## client

- RocketMQ Go client: https://github.com/apache/rocketmq-client-go
- kafka-go: https://github.com/segmentio/kafka-go
- Nats.go: github.com/nats-io/nats.go

> If it is Alibaba Cloud RocketMQ: You can use the official library


## Precautions

- Consumers need to do idempotent processing because they may receive the same message multiple times
- Record logs when consuming to facilitate subsequent locating problems. It is best to add the unique identifier of the request, such as fields such as request_id or trace_id
- Try to consume in batches, which can greatly improve the consumption throughput


## Reference

- [RocketMQ official website](https://rocketmq.apache.org/)
- [RocketMQ Documentation](https://rocketmq.apache.org/docs/quick-start/)
- [RocketMQ Go client](https://github.com/apache/rocketmq-client-go)
- [RocketMQ Go client documentation](https://github.com/apache/rocketmq-client-go/blob/master/docs/Introduction.md)
- [Alibaba Cloud RocketMQ](https://cn.aliyun.com/product/rocketmq)
- https://github.com/GSabadini/go-message-broker/blob/master/main.go
- [Automatically recovering RabbitMQ connections in Go applications](https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59)