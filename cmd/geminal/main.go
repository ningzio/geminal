package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ningzio/geminal/internal"
	"github.com/ningzio/geminal/internal/llm"
	"github.com/ningzio/geminal/internal/repo"
	"github.com/ningzio/geminal/tui"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	geminalDir := filepath.Join(homeDir, ".geminal")
	err = os.MkdirAll(geminalDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(filepath.Join(geminalDir, "geminal.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	apiKey := os.Getenv("API_KEY")
	ai, err := llm.NewGeminiAI(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	r, err := repo.NewRepository()
	if err != nil {
		log.Fatalf("init repo: %s", err)
	}
	h := internal.NewHandler(
		ai,
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
