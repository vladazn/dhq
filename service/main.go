package main

import (
	"fmt"
	"github.com/vladazn/dhq/service/app"
	"os"
)

func main() {
	p, _ := os.Getwd()
	app.Run(fmt.Sprintf("%s/service/config/config.yml", p))
}
