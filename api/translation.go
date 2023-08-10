package api

import (
	"io"
	"log"
	"net"
	"github.com/henson/proxypool/pkg/storage"
	"fmt"
	"time"
)



func handleClient(client net.Conn) error {
	defer client.Close()

	data := make([]byte, 1024)
	_, err := client.Read(data)
	if err != nil {
		fmt.Printf("[-]Local receiving client: %v\n", err)
		return nil
	}
	fmt.Printf("[*%s]: Accept data...\n", time.Now())

	ips:=storage.ProxyRandom()
	mbsocket, err := net.DialTimeout("tcp", ips.Data, 3*time.Second)
	if err != nil {
		fmt.Printf("[-]RE_Connect...\n")
		return nil
	}
	defer mbsocket.Close()

	done := make(chan struct{})

	// Forward data from client to proxy
	go func() {
		_, err := io.Copy(mbsocket, client)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("[-]Error while sending data to proxy: %v\n", err)
			}
		}
		close(done)
	}()

	// Forward data from proxy to client
	_, err = io.Copy(client, mbsocket)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("[-]Error while sending data back to client: %v\n", err)
		}
	}

	<-done
	return nil
}


func  handleClient1(client net.Conn) {
	defer client.Close()

	data := make([]byte, 1024)
	_, err := client.Read(data)
	if err != nil {
		fmt.Printf("[-]Local receiving client: %v\n", err)
		return 
	}
	fmt.Printf("[*%s]: Accept data...\n", time.Now())

	ips:=storage.ProxyRandom()
	mbsocket, err := net.DialTimeout("tcp", ips.Data, 3*time.Second)
	if err != nil {
		fmt.Printf("[-]RE_Connect...\n")
		return
	}
	defer mbsocket.Close()

	_, err = mbsocket.Write(data)
	if err != nil {
		fmt.Printf("[-]Sent to the proxy server: %v\n", err)
		return
	}

	data1 := make([]byte, 1024)
	_, err = mbsocket.Read(data1)
	if err != nil {
		fmt.Printf("[-]Back to the client: %v\n", err)
		return
	}
	fmt.Printf("[*%s]: Send data...\n", time.Now())

	_, err = client.Write(data1)
	if err != nil {
		fmt.Printf("[-]Error sending data back to client: %v\n", err)
	}
	
}

func copyData(dst net.Conn, src net.Conn) {
	const timeout = 5 * time.Second
	buffer := make([]byte, 4096)

	defer dst.Close()
	defer src.Close()

	for {
		// 设置读取超时
		src.SetReadDeadline(time.Now().Add(timeout))
		n, err := src.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				 
				break
			}
			// 其他读取错误，退出
			log.Printf("Error reading from source: %s", err)
			return
		}

		if n > 0 {
			// 写入数据到目标连接
			if _, err := dst.Write(buffer[:n]); err != nil {
				log.Printf("Error writing to destination: %s", err)
				return
			}
		}
	}
}




func Translation(listenAddr string) {
	
	
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to bind on address: %s", listenAddr)
	}
	defer listener.Close()

	
	
	
	
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}

		 go handleClient(client)

		
	}

}
