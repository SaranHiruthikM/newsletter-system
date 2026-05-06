package main

import (
	"fmt"

	"github.com/SaranHiruthikM/newsletter-system/internal/config"
)

func main() {
	cfg := config.Load()

	fmt.Println("APP PORT:", cfg.App.Port)
	fmt.Println("RedisHOST", cfg.Redis.Host)
	fmt.Println("EmailProvider:", cfg.Email.Provider)
}
