package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startTime = time.Now()

func main() {
	// Ganti dengan Token Bot Anda
	botToken := "8778956451:AAFkfHFetDkvyjNM7ilk55THyiwmWBjju8Y" 

	// Arahkan ke domain Railway Local API Server Anda
	customEndpoint = "https://telegram-bot-api-production-6598.up.railway.app/bot%s/%s"
	// Inisialisasi bot langsung dengan Custom Endpoint
	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint(botToken, customEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "Markdown"

		switch update.Message.Command() {
		case "start":
			msg.Text = "Sistem aktif. Gunakan /ping untuk metrik server."

		case "me":
			user := update.Message.From
			msg.Text = fmt.Sprintf("ID: `%d`\nUsername: @%s\nName: %s %s", 
				user.ID, user.UserName, user.FirstName, user.LastName)

		case "ping":
			msg.Text = getSystemMetrics()
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

// getSystemMetrics menarik metrik performa backend secara real-time
func getSystemMetrics() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime).Round(time.Second)
	
	// Konversi byte ke Megabyte
	allocMB := float64(m.Alloc) / 1024 / 1024
	sysMB := float64(m.Sys) / 1024 / 1024

	return fmt.Sprintf(
		"*Backend Performance*\n"+
			"Uptime: `%s`\n"+
			"Goroutines: `%d`\n"+
			"RAM Alloc: `%.2f MB`\n"+
			"RAM Sys: `%.2f MB`",
		uptime, runtime.NumGoroutine(), allocMB, sysMB,
	)
}