/*
Программа генерирует лог в директории logs/generated.log
через неравные промежутки времени.
Каждая строка имеет вид:


*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const logFile = "./logs/generated.log"
const maxSleepingTime = 2000

func main() {
	for {
		sleepTime := time.Duration(rand.Int63n(maxSleepingTime)) * time.Millisecond
		addLineToLog(sleepTime)
		time.Sleep(sleepTime)
	}
}

func addLineToLog(t time.Duration) {

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	message := fmt.Sprintf(" - generation time - %d", t/1000000)
	if t < maxSleepingTime*1000000/5 {
		message = "- ERROR" + message
	} else {
		message = "- INFO" + message
	}

	logger.Printf(message)
	fmt.Println(message)
}
