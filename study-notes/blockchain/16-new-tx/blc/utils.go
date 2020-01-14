package blc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

// IntToHex 实现int64转[]byte
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if nil != err {
		log.Panicf("int transact to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// JSONToSlice 标准JSON格式转切片
// mac下需要添加引号
// $ ./bc send -from "[\"a\"]" -to "[\"b\"]" -amount "[\"20\"]"
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	// json
	if err := json.Unmarshal([]byte(jsonString), &strSlice); nil != err {
		log.Printf("json to []string failed! %v\n", err)
	}
	return strSlice
}
