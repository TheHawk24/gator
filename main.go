package main

import (
	"fmt"

	"github.com/TheHawk24/gator/internal/config"
)

func main() {

	config_file := config.Read()
	config_file.Current_Username = "Lehloenya"
	config_file.SetUser()

	config_username := config.Read()
	fmt.Println(config_username)
}
