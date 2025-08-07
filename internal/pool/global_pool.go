package pool

import (
	"log"

	"github.com/panjf2000/ants"
)

var logPool *ants.Pool

func init() {
	var err error
	logPool, err = ants.NewPool(2000) // max 2000 concurrent logging jobs
	if err != nil {
		log.Fatal("Gagal membuat goroutine pool", err)
	}
}
