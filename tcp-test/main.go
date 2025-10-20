// package main implements tcp test
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Request struct {
	Method   string
	Path     string
	Params   string
	Protocol string
	Headers  map[string]string
}

func (r Request) String() string {
	jsonHeaders, _ := json.MarshalIndent(r.Headers, "", "    ")
	return fmt.Sprintf(`
Method: %v
Path: %v
Protocol: %v
Params: %v
Headers: 
%v
	`, r.Method, r.Path, r.Protocol, r.Params, string(jsonHeaders))
}

func main() {
	//#nosec
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("[error]: %v", err)
			return
		}
	}()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			return
		}

		log.Printf("[error]: %v", err)
		return
	}

	reqStr := string(buf[:n])

	reqSlice := strings.Split(reqStr, "\r\n")

	req := parseRequest(reqSlice)

	fmt.Printf("REQUEST: \n%v\n", req)

	if req.Method == "POST" && req.Headers["Content-Type"] == "application/json" {
		reqBody := reqSlice[len(reqSlice)-1]
		jsonMap := handleJSONBody(reqBody)
		jsonRespBytes, err := json.Marshal(jsonMap)
		if err != nil {
			log.Fatalf("Failed to marshal JSON response: %v", err)
		}

		jsonPretty, _ := json.MarshalIndent(jsonMap, "", "    ")

		fmt.Fprintf(conn, "HTTP/1.1 200 OK \r\nConnection: close\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", len(jsonRespBytes), jsonRespBytes)
		fmt.Printf("RESPONSE: \n\n%v\n", string(jsonPretty))
		return
	}
	body := "Hello TCP!"
	fmt.Printf("RESPONSE: \n\n\t%v\n", body)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nConnection: close\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
}

func handleJSONBody(body string) map[string]any {
	jsonMap := map[string]any{}
	err := json.Unmarshal([]byte(body), &jsonMap)
	if err != nil {
		log.Fatalf("Could not parse JSON: %v", err)
	}

	return jsonMap
}

func parseRequest(reqSlice []string) *Request {
	firstLine := strings.Fields(reqSlice[0])
	var (
		method   = firstLine[0]
		path     = strings.Split(firstLine[1], "?")[0]
		params   = strings.Split(firstLine[1], path)[1]
		protocol = firstLine[2]
	)

	headers := map[string]string{}
	for i := 1; i < len(reqSlice); i++ {
		if !strings.Contains(reqSlice[i], ":") {
			break
		}

		keyVal := strings.Split(reqSlice[i], ":")
		headers[keyVal[0]] = strings.Trim(keyVal[1], " ")
	}

	req := Request{
		Method:   method,
		Path:     path,
		Protocol: protocol,
		Params:   params,
		Headers:  headers,
	}

	return &req
}
