package main

import (
	"log"
	"net"

	"github.com/tarm/serial"
)

func writePort(portNum *serial.Port, data string) {
	buf := []byte(data)

	n, err := portNum.Write(buf)
	if err != nil {
		log.Fatalf("Ошибка при записи данных в COM порт: %v", err)
	} else {
		log.Printf("Запись в COM порт успешна, записанно %d байт", n)
	}

}

func readFromPort(port *serial.Port, conn net.Conn) {
	buf := make([]byte, 128)
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Fatalf("Ошибка чтения данных из COM порта: %v", err)
		}
		if n > 0 {
			data := make([]byte, n)
			_, err = conn.Write(data)
			if err != nil {
				log.Fatalf("Ошибка при передаче данных клиенту: %v", err)
				return
			}
		}
	}
}

func handleConnection(conn net.Conn, port *serial.Port) {
	defer conn.Close()
	log.Printf("Подключен клиент: %s", conn.RemoteAddr())
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalf("Ошибка при чтении данных от клиента: %v", err)
				break
			}
			if n > 0 {
				data := string(buf[:n])
				writePort(port, data)

			}
		}
	}()
	go readFromPort(port, conn)
	select {}

}

func main() {
	config := &serial.Config{
		Name:     "COM1",
		Baud:     9600,
		Size:     8,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
	}

	listener, err := net.Listen("tcp", ":15675")
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	defer listener.Close()

	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Ошибка открытия порта: %v", err)
	}
	defer port.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Ошибка при принятии соединения: %v", err)
	}
	handleConnection(conn, port)

}
