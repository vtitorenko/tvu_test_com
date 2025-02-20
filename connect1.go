package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tarm/serial"
)

func main() {
	// Настройка параметров COM-порта
	port := "/dev/ttyUSB0"
	cfg := &serial.Config{
		Name:        port,
		Baud:        9600,
		ReadTimeout: 10 * 60 * 1000, // 10 минут в миллисекундах
	}

	// Открытие COM-порта
	s, err := serial.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// Настройка обработчика сигнала завершения
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChannel
		s.Close()
		os.Exit(0)
	}()

	// Цикл чтения/записи данных регистров
	for {
		// Чтение данных из регистра с адресом 0x01
		addr := []byte{0x01}
		_, err := s.Write(addr)
		if err != nil {
			log.Println("Ошибка записи адреса регистра:", err)
			continue
		}
		time.Sleep(10 * time.Millisecond) // Задержка для стабильной работы

		// Чтение данных из COM-порта
		var data [2]byte
		_, err = s.Read(data[:])
		if err != nil {
			log.Println("Ошибка чтения данных:", err)
			continue
		}
		value := int(data[0])<<8 + int(data[1])
		fmt.Printf("Чтение данных из регистра 0x01: 0x%04X\n", value)

		// Запись данных в регистр с адресом 0x02
		value = 0x1234
		data[0] = byte(value >> 8)
		data[1] = byte(value & 0xFF)
		_, err = s.Write(data[:])
		if err != nil {
			log.Println("Ошибка записи данных:", err)
			continue
		}
		fmt.Printf("Запись данных в регистр 0x02: 0x%04X\n", value)
	}
}
