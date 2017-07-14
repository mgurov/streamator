package main

import (
	"sync"
	"net/http"
	"github.com/sirupsen/logrus"
	"fmt"
	"strconv"
)

type restServer struct {
	quit chan interface{}
}

func startRestServer(port int, wg *sync.WaitGroup, data dataProvider) *restServer {
	wg.Add(1)

	http.HandleFunc("/ticks", func(w http.ResponseWriter, r *http.Request) {

		formatter := logrus.JSONFormatter{}
		logRecords := data.Get()

		w.Header().Add("Content-type", "text/json")

		fmt.Fprint(w, "[")
		for i, rec := range logRecords {
			recBytes, err := formatter.Format(rec)
			if err != nil {
				log.Info("Could not convert log record ", rec, err)
				continue
			}
			w.Write(recBytes)
			if i != len(logRecords)-1 {
				fmt.Fprint(w, ",")
			}
		}
		fmt.Fprint(w, "]")
	})

	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: nil}

	go func() {
		log.Info("Starting at port ", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
		wg.Done()
	}()

	s := &restServer{
		quit: make(chan interface{}),
	}

	go func() {
		<-s.quit
		log.Info("Stopping the web server")
		server.Close()
	}()

	return s
}

func (s *restServer) Stop() {
	s.quit <- nil
}
