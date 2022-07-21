package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ready-to-go/framework/cgin"
	"time"
)

/* func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(DemoHandler))

	srv := http.Server{Addr: ":8080", Handler: mux}
	srv.ListenAndServe()

}

func DemoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("got:", r.URL.Path)
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, "Hello World!")

} */

func main() {

	e := cgin.NewEngine()

	e.Get("/", func(c *cgin.Context) {
		time.Sleep(time.Second * 10)
		c.Json(200, "hello world!!")
		log.Println("req done")

	})

	srv := http.Server{Addr: ":8080", Handler: e}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		srv.ListenAndServe()
	}()

	<-sig

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println("shutdown:", err)
	}

}
