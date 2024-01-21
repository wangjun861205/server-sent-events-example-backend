package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerFactory = func(keeper *AddrKeeper) http.Handler

func PullMsgFactory(keeper *AddrKeeper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := mux.Vars(r)["uid"]
		outChan := make(chan string)
		go func() {
			for msg := range outChan {
				if _, err := w.Write([]byte(fmt.Sprintf("event: message\ndata: %s\n\n", msg))); err != nil {
					log.Println(err)
				}
				w.(http.Flusher).Flush()
			}
		}()
		keeper.register(uid, outChan)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Encoding", "none")

		<-r.Context().Done()
		keeper.unregister(uid)
		close(outChan)
	})
}

func SendMsgFactory(keeper *AddrKeeper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		uid := mux.Vars(r)["uid"]
		outChan := keeper.getAddr(uid)
		if outChan == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("用户不存在"))
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
