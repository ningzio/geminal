package main

import (
	"log"
	"os"

	"github.com/ningzio/geminal/internal"
	"github.com/ningzio/geminal/internal/llm"
	"github.com/ningzio/geminal/internal/repo"
	"github.com/ningzio/geminal/tui"
)

func main() {
	f, err := os.OpenFile("geminal.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	r, err := repo.NewRepository()
	if err != nil {
		log.Fatalf("init repo: %s", err)
	}
	h := internal.NewHandler(
		&llm.Mock{},
		r,
		internal.NewChromaRenderer(),
	)

	app, err := tui.NewApplication(h)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
