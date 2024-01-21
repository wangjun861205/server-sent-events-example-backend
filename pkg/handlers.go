package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type HandlerFactory = func(keeper *AddrKeeper) http.Handler

func PullMsgFactory(keeper *AddrKeeper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		segs := strings.Split(r.URL.Path, "/")
		id := segs[len(segs)-2]
		outChan := make(chan string)
		go func() {
			for msg := range outChan {
				if _, err := w.Write([]byte(fmt.Sprintf("event: message\ndata: %s\n\n", msg))); err != nil {
					log.Println(err)
				}
				w.(http.Flusher).Flush()
			}
		}()
		keeper.register(id, outChan)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Encoding", "none")

		<-r.Context().Done()
		keeper.unregister(id)
		close(outChan)
	})
}

func SendMsgFactory(keeper *AddrKeeper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		segs := strings.Split(r.URL.Path, "/")
		id := segs[len(segs)-2]
		outChan := keeper.getAddr(id)
		if outChan == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("addr not exists"))
			return
		}
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		outChan <- string(msg)
	})
}
