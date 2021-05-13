package app

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"sync"
	"time"
)

type Producer struct {
	Config             *Config
	Watcher            *fsnotify.Watcher
	done               chan struct{} // Channel for sending a "quit message" to the reader goroutine
	mu                 sync.Mutex
	isClosed           bool
	lines              *Lines
	workerNum          int
	sleepBeforeSending int
}

func NewProducer(config *Config, watcher *fsnotify.Watcher, lines *Lines, appConfig *Application, workerNum int) (*Producer, error) {

	w := &Producer{
		Config:    config,
		Watcher:   watcher,
		done:      make(chan struct{}),
		lines:     lines,
		workerNum: workerNum,
	}
	for _, v := range appConfig.Rules {
		w.sleepBeforeSending = v.SleepBeforeSending
	}

	go w.eventLoop()
	return w, nil
}

func (p *Producer) Close() error {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return nil
	}
	p.isClosed = true
	p.mu.Unlock()
	close(p.done)

	return nil
}

func (p *Producer) eventLoop() {
	log.Info("Start eventloop producer: ", p.workerNum)
	for {
		select {
		case event, ok := <-p.Watcher.Events:
			if !ok {
				return
			}
			log.Debug("event:", event)
			if event.Op&fsnotify.Create == fsnotify.Create {

				log.Debug("modified file: ", event.Name, " producer number: ", p.workerNum)

				b, err := ioutil.ReadFile(event.Name)
				if err != nil {
					log.Error(err)
				}
				time.Sleep(time.Second * time.Duration(p.sleepBeforeSending))
				p.lines.C <- b

			}
		case err, ok := <-p.Watcher.Errors:
			if !ok {
				return
			}
			log.Error("error:", err)
		}
	}
}
