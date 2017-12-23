package messages

import (
	"github.com/bsm/sarama-cluster"
	"github.com/Shopify/sarama"
	"github.com/satori/go.uuid"
	"os"
	"syscall"
	"fmt"
	"os/signal"
	"log"
)



type Consumer struct {
	consumer *cluster.Consumer
	callbacks ConsumerCallbacks
}

func NewConsumer(callbacks ConsumerCallbacks,brokerList []string, groupId string, topics []string) *Consumer {
	consumer := Consumer{callbacks:callbacks}

	config := cluster.NewConfig()
	config.ClientID = uuid.NewV4().String()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConsumer, err := cluster.NewConsumer(brokerList, groupId, topics, config)
	if err != nil {
		panic(err)
	}
	consumer.consumer = saramaConsumer
	return &consumer

}

func (c *Consumer) Consume() {
	// Create signal channel
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// Consume all channels, wait for signal to exit
	go func(){
		for {
			select {
			case msg, more := <-c.consumer.Messages():
				if more {
					if c.callbacks.OnDataReceived!=nil {
						c.callbacks.OnDataReceived(msg.Value)
					}
					fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
					c.consumer.MarkOffset(msg, "")
				}
			case ntf, more := <-c.consumer.Notifications():
				if more {
					log.Printf("Rebalanced: %+v\n", ntf)
				}
			case err, more := <-c.consumer.Errors():
				if more {
					if c.callbacks.OnError!=nil {
						c.callbacks.OnError(err)
					}
					//logger.Printf("Error: %s\n", err.Error())
				}
			case <-sigchan:
				return
			}
		}
	}()

}

func (c *Consumer) Close(){
	c.consumer.Close()
}

type ConsumerCallbacks struct {
	OnDataReceived func(msg []byte)
	OnError func(err error)
}
