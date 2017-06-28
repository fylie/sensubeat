package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/fylie/sensubeat/beater"
)

func main() {
	err := beat.Run("sensubeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
