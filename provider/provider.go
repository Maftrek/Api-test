package provider

import (
	"fmt"
	"github.com/adjust/rmq"
	stan "github.com/nats-io/go-nats-streaming"
	"log"
)

// Provider interface
type Provider interface {
	ConnectNatsStreaming(dpURL string) error
	GetNatsConnectionStreaming() (stan.Conn, error)
}

type provider struct {
	snats    stan.Conn
	rmqRedis rmq.Connection
}

func NewNATS() Provider {
	return &provider{}
}

func NewRMQ() Provider {
	return &provider{}
}

// метод инициализации к серверу NATS streaming
func (p *provider) ConnectNatsStreaming(dpURL string) error {
	clientID := "nats"
	snc, err := stan.Connect("test-cluster", clientID, stan.NatsURL(dpURL), stan.MaxPubAcksInflight(1))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, dpURL)
		return err
	}
	p.snats = snc
	fmt.Printf(`{"clientID": %s}`, clientID)
	return nil
}

// метод для возвращения коннекшена к NATS streaming
func (p *provider) GetNatsConnectionStreaming() (stan.Conn, error) {
	return p.snats, nil
}
