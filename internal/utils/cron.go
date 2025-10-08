package utils

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"
)

// PreventSleepCron starts a cron job that runs every 5 seconds and every 5 minutes
func PreventSleepCron() {
	c := cron.New(cron.WithSeconds()) 

	// Every 5 minutes
	c.AddFunc("*/5 * * * *", func() {
		pingAPI()
	})

	c.Start()
}

func pingAPI() {
	resp, err := http.Get("https://restaurant-backend-srp6.onrender.com/ping")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Ping successful:", resp.Status)
}
