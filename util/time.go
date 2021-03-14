package util

import (
	"log"
	"time"
)

func LogTime(label string, task func()) {
	start := time.Now()
	task()
	log.Println("Completed", label, "in", time.Now().Sub(start))
}
