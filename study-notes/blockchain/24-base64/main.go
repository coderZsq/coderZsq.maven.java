package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	msg := "coderZsq"
	// 编码
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Printf("encode result: %v\n", encoded)
	b, err := base64.StdEncoding.DecodeString("Y29kZXJac3E=")
	if nil != err {
		panic(err)
	}
	fmt.Printf("decode result: %s\n", b)
}
