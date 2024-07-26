package main

import (
	"github.com/assimon/ai-anti-bot/bot"
	"github.com/assimon/ai-anti-bot/pkg/os"
	"log"
)

func main() {
	go func() {
		if err := bot.Start(); err != nil {
			log.Fatalf("%s\n", err.Error())
		}
	}()
	os.WaitSignalChina()
}
