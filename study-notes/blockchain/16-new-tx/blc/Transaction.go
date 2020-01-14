package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// 交易管理文件

// Transaction 定义一个交易基本结构
type Transaction struct {
	// 交易哈希(标识)
	TxHash []byte
	// 输入列表
	Vins []*TxInput
	// 输出列表
	Vouts []*TxOutput
}

// NewCoinbaseTransaction 实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	// 输入:
	// coinbase特点
	// txHash: nil
	// vout: -1(为了对是否是coinbase交易进行判断)
	// ScriptSig: 系统奖励
	txInput := &TxInput{[]byte{}, -1, "system reward"}
	// 输出:
	// value
	// address
	txOutput := &TxOutput{10, address}
	// 输入输出组装交易
	txCoinbase := &Transaction{nil, []*TxInput{txInput}, []*TxOutput{txOutput}}
	// 交易哈希生成
	txCoinbase.HashTransaction()
	return txCoinbase
}

// HashTransaction 生成交易哈希(交易序列化)
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// 设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed %v\n", err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// NewSimpleTransaction 生成普通转账交易
func NewSimpleTransaction(from string, to string, amount int) *Transaction {
	var txInputs []*TxInput   // 输入列表
	var txOutputs []*TxOutput // 输出列表
	// 输入
	txInput := &TxInput{[]byte("44765ac9d57ed63865e9a018d4ec3fb4a3dfb6aa1fa19f6fa3787f8440a42d84"), 0, from}
	txInputs = append(txInputs, txInput)
	// 输出 (转账源)
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)
	// 输出 (找零)
	if amount < 10 {
		txOutput = &TxOutput{10 - amount, from}
		txOutputs = append(txOutputs, txOutput)
	}
	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction()
	return &tx
}
