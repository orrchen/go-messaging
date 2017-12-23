package messages

import (
	"github.com/Shopify/sarama"
	"log"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"
	"encoding/json"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	callbacks     ProducerCallbacks
	topic         string
}

func NewProducer(callbacks ProducerCallbacks,brokerList []string,topic string,certFile *string,keyFile *string,caFile *string,verifySsl *bool ) *Producer {
	producer := Producer{ callbacks: callbacks, topic: topic}

	config := sarama.NewConfig()
	tlsConfig := createTlsConfiguration(certFile,keyFile,caFile,verifySsl)
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	saramaProducer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
		panic(err)
	}
	go func() {
		for err := range saramaProducer.Errors() {
			if producer.callbacks.OnError!=nil {
				producer.callbacks.OnError(err)
			}
		}
	}()
	producer.asyncProducer = saramaProducer
	return &producer
}

func createTlsConfiguration(certFile *string,keyFile *string,caFile *string,verifySsl *bool)(t *tls.Config) {
	if certFile!=nil && keyFile!=nil && caFile!=nil && *certFile != "" && *keyFile != "" && *caFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: *verifySsl,
		}
	}
	// will be nil by default if nothing is provided
	return t
}

func (p *Producer) Produce(payload interface{}) {
	value := message{
		value: payload,
	}
	value.ensureEncoded()
	log.Println("producing: ", string(value.encoded))
	p.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		//Key:   sarama.StringEncoder(r.RemoteAddr),
		Value: &value,
	}
}

type message struct {
	value interface{}
	encoded []byte
	err     error
}

func (ale *message) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale.value)
		if ale.err!=nil {
			log.Println(ale.err)
		}
	}
}

func (ale *message) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *message) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

func (p *Producer) Close() error{
	log.Println("Producer.Close()")
	if err := p.asyncProducer.Close(); err != nil {
		return err
	}
	return nil
}

type ProducerCallbacks struct {
	OnError func(error)
}



