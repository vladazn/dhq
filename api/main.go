package main

import (
	"fmt"
	"github.com/vladazn/dhq/api/app"
	"os"
)

func main() {
	p, _ := os.Getwd()
	app.Run(fmt.Sprintf("%s/api/config/config.yml", p))
}
