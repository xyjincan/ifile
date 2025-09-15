package main

import (
	"fmt"
	"ifile/tool"
)

func main() {
	message := fmt.Sprintf("%v starting...", "ifile")
	fmt.Println(message)
	tool.StartHttpd()
}
