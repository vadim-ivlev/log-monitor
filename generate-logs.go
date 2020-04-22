/*
Программа генерирует лог в директории ./logs/logrus.log
через неравные промежутки времени.
Каждая строка имеет вид:
{"level":"info","msg":"Рейс задерживается на 1.449s","status":"INFO","time":"2020-04-22T15:21:07+03:00","wait":1449}

*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// имя файла лога
const logFile = "./logs/logrus.log"

// максимальное время задержки между записями в лог
const maxSleepingTime = 4000

// логгер для вывода на экран
var stdoutLog = log.New()

// лггер для вывода в файл
var fileLog = log.New()

func main() {
	// Инициализируем логгер
	initLogger()

	// Пока пользователь не нажал Ctrl-C выполняем Вечный цикл.
	for {
		// Вычисляем время задержки
		sleepTime := time.Duration(rand.Int63n(maxSleepingTime)) * time.Millisecond
		// Добавляем линию  в лог файл
		addLineToLog(sleepTime)
		// Ждём
		time.Sleep(sleepTime)
	}
}

// Инициализируем логгер
func initLogger() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	stdoutLog.SetFormatter(&log.JSONFormatter{})
	fileLog.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// log.SetOutput(os.Stdout)
	stdoutLog.Out = os.Stdout

	// // Устанавливаем вывод в файл
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// log.SetOutput(file)
		fileLog.Out = file
	} else {
		log.Info("Вывод в файл невозможен, используем stdout")
		fileLog.Out = os.Stdout
	}

	// Only log the warning severity or above.
	// log.SetLevel(log.WarnLevel)
}

// Добавляет строку в файл лога
func addLineToLog(t time.Duration) {
	msg := fmt.Sprintf("Рейс задерживается на %v", t)
	fields := logrus.Fields{
		"status": choose(t < maxSleepingTime*1000000/5, "ERROR", "INFO"),
		"wait":   int64(t / 1000000),
		"format": "JSON",
	}

	if t < maxSleepingTime*1000000/5 {
		fileLog.WithFields(fields).Error(msg)
		stdoutLog.WithFields(fields).Error(msg)
	} else {
		fileLog.WithFields(fields).Info(msg)
		stdoutLog.WithFields(fields).Info(msg)
	}

}

// Возвращает одну из строк в зависимости от условия
func choose(condition bool, s1, s2 interface{}) interface{} {
	if condition {
		return s1
	}
	return s2
}
