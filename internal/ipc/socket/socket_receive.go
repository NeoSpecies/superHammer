package ipc

import (
	"bigHammer/internal/plugin"
	"bigHammer/internal/shared"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSocket(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// log.Println("Connection closed by peer immediately after connecting")
		} else {
			log.Println("Error reading data from socket:", err)
		}
		return
	}

	pluginInterface, err := shared.GlobalContainer.Resolve("plugin")
	if err != nil {
		log.Println("Error resolving plugin from container:", err)
		return
	}

	var req plugin.Request
	if err := json.Unmarshal(data[:n], &req); err != nil {
		log.Println("Error parsing request:", err)
		return
	}

	pluginInstance, ok := pluginInterface.(plugin.ServicePlugin)
	if !ok {
		log.Println("The provided plugin service does not implement the ServicePlugin interface.")
		return
	}

	response := pluginInstance.HandleRequest(req)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling response:", err)
		return
	}
	_, err = conn.Write(responseJSON)
	if err != nil {
		log.Println("Error writing data to socket:", err)
		return
	}

	// 主循环，持续读取数据
	for {
		_, err := conn.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// log.Println("Connection closed by peer")
				break
			} else {
				log.Println("Error reading data from socket:", err)
				return
			}
		}

	}
}
