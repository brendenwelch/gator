package main

import (
	"fmt"
	"github.com/brendenwelch/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config")
	}

	if err = cfg.SetUser("brenden"); err != nil {
		fmt.Println("error setting user in config")
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("error reading config")
	}

	fmt.Println("Config contents:")
	fmt.Printf("- Database URL: %v\n", *cfg.Db_url)
	fmt.Printf("- Username: %v\n", *cfg.Current_user_name)
}
