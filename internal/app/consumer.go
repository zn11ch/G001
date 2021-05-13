package app

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Consumer struct {
	Config          *Config
	done            chan struct{}
	mu              sync.Mutex
	isClosed        bool
	lines           *Lines
	workerNum       int
	logEveryMessage bool
}

func NewConsumer(config *Config, lines *Lines, appConfig *Application, workerNum int) (*Consumer, error) {
	c := &Consumer{
		Config:    config,
		done:      make(chan struct{}),
		lines:     lines,
		workerNum: workerNum,
	}

	for _, v := range appConfig.Rules {
		c.logEveryMessage = v.LogEveryMessage
	}

	go c.eventLoop()
	return c, nil
}

func (c *Consumer) Close() error {
	c.mu.Lock()
	if c.isClosed {
		c.mu.Unlock()
		return nil
	}
	c.isClosed = true
	c.mu.Unlock()
	close(c.done)

	return nil
}

func (c *Consumer) eventLoop() {
	log.Info("Start eventloop consumer: ", c.workerNum)
	for {
		select {
		case data, ok := <-c.lines.C:
			if ok {
				person := &Person{}
				err := json.Unmarshal(data, person)
				if err != nil {
					log.Error(err)
				}
				log.Info(person, " consumer number: ", c.workerNum)
			}
		}
	}
}
