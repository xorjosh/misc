package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	serverURL   = "ws://***.***.***.***:****/socket.io/?EIO=3&transport=websocket"
	logFileName = "bot_connections.log"
)

func main() {
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	for {
		err := connectAndMonitor(serverURL, logFile)
		if err != nil {
			log.Println("Connection error:", err)
			log.Println("Waiting for 15 seconds before attempting to reconnect...")
			time.Sleep(15 * time.Second)
		}
	}
}

func connectAndMonitor(serverURL string, logFile *os.File) error {
	c, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	registerMessage := `42["registerAdmin",{"usrType":"admin","usrName":"kiran","pass":"suthar"}]`
	if err := c.WriteMessage(websocket.TextMessage, []byte(registerMessage)); err != nil {
		return err
	}

	fmt.Println("Connected to server. Monitoring for bot joins...")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return err
		}

		
		processMessage(string(message), logFile)
	}
}

func processMessage(message string, logFile *os.File) {
	
	if len(message) > 2 && message[:2] == "42" {
		
		data := message[2:]

		if data != "null" {
			

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			logMessage := fmt.Sprintf("[%s] Bot joined: %s\n", timestamp, data)
			fmt.Print(logMessage)

			
			if _, err := logFile.WriteString(logMessage); err != nil {
				log.Println("Error writing to log file:", err)
			}
		}
	}
}

