// Package main provides main  INFO:
package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "GET /hello HTTP/1.1\r\n"

	fmt.Println(strings.Fields(strings.Split(str, "\r\n")[0])[0])
}
