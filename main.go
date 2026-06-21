package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"github.com/Omkardalvi01/redis_go.git/internal/cli"
	"github.com/Omkardalvi01/redis_go.git/internal/engine"
	"github.com/Omkardalvi01/redis_go.git/internal/httpapi"
)

func main() {
	eng, err := engine.New(engine.Config{
		LogPath:   "log.txt",
		EnableAOF: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", httpapi.New(eng))

	log.Println("server listening on :8000")

	go cli.Run(bufio.NewScanner(os.Stdin), eng, os.Stdout, os.Stderr)

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
