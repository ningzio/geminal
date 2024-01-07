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

	// apiKey := os.Getenv("API_KEY")
	// ai, err := llm.NewGeminiAI(apiKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r, err := repo.NewRepository()
	if err != nil {
		log.Fatalf("init repo: %s", err)
	}
	h := internal.NewHandler(
		// ai,
		&llm.Mock{},
		r,
		internal.NewChromaRenderer(),
		// internal.NewMarkdownRenderer(),
	)

	app, err := tui.NewApplication(h)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// func foo() {
// 	f, err := os.ReadFile("/Users/ningzi/workspace/personal/geminal/internal/llm/mark.log")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	render := internal.NewChromaRenderer()

// 	render.Content(os.Stdout, f)
// }
