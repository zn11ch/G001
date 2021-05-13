package app

import "sync"

type Lines struct {
	C    chan []byte
	once sync.Once
}

func NewLinesChannel() *Lines {
	return &Lines{C: make(chan []byte)}
}
