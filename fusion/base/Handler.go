package base

import (
	"log"
	"net"
	"runtime/debug"
	"sync"
)

var wg sync.WaitGroup

func WaitHandler() {
	wg.Wait()
}

func SafeHandler(handler func()) func() {
	return func() {
		defer func() {
			defer wg.Done()
			if err := recover(); err != nil {
				log.Println(err)
				debug.PrintStack()
			}
		}()
		wg.Add(1)
		handler()
	}
}

func SafeConnHandler(handler func(conn net.Conn)) func(conn net.Conn) {
	return func(conn net.Conn) {
		defer func() {
			defer wg.Done()
			if err := recover(); err != nil {
				log.Println(err)
				debug.PrintStack()
			}
		}()
		wg.Add(1)
		handler(conn)
	}
}
