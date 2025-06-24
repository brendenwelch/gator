package main

import (
	"fmt"
	"log"

	"github.com/brendenwelch/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	if err = cfg.SetUser("brenden"); err != nil {
		log.Fatalf("error setting user in config: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Println("Config contents:")
	fmt.Printf("- Database URL: %v\n", cfg.Db_url)
	fmt.Printf("- Username: %v\n", cfg.Current_user_name)
}
