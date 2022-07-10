package services

import (
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Publisher struct {
	sync.Mutex
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

var rabbitUrl = os.Getenv("RABBIT_URL")

func (p *Publisher) Publish(data []byte) error {
	p.Lock()
	defer p.Unlock()
	for i := 0; i < 10; i += 1 {
		if p.connection == nil || p.connection.IsClosed() {
			time.Sleep(1 * time.Second)
			var err error
			p.connection, err = amqp.Dial(rabbitUrl)
			if err != nil {
				log.Println(err)
				continue
			}
			p.channel, err = p.connection.Channel()
			if err != nil {
				log.Println(err)
				p.connection.Close()
				continue
			}
			p.queue, err = p.channel.QueueDeclare("stop", true, false, false, false, nil)
			if err != nil {
				log.Println(err)
				p.channel.Close()
				p.connection.Close()
				continue
			}
		}
		if err := p.channel.Publish("", p.queue.Name, false, false, amqp.Publishing{Body: data}); err != nil {
			log.Println(err)
			p.channel.Close()
			p.connection.Close()
			continue
		}
		return nil
	}
	return http.ErrLineTooLong
}
