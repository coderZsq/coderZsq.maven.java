package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
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
	txInput := &TxInput{[]byte{}, -1, nil, nil}
	// 输出:
	// value
	// address
	//txOutput := &TxOutput{10, StringToHash160(address)}
	txOutput := NewTxOutput(10, address)
	// 输入输出组装交易
	txCoinbase := &Transaction{nil, []*TxInput{txInput}, []*TxOutput{txOutput}}
	// 交易哈希生成
	txCoinbase.HashTransaction()
	return txCoinbase
}

// HashTransaction 生成交易哈希(交易序列化)
// 不同时间生成的交易哈希值不同
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// 设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed %v\n", err)
	}
	// 添加时间戳标识, 不添加会导致所有的coinbase交易哈希完全相同
	tm := time.Now().UnixNano()
	// 用于生成哈希的源数据
	txHashBytes := bytes.Join([][]byte{result.Bytes(), IntToHex(tm)}, []byte{})
	// 生成哈希值
	hash := sha256.Sum256(txHashBytes)
	tx.TxHash = hash[:]
}

// NewSimpleTransaction 生成普通转账交易
func NewSimpleTransaction(from string, to string, amount int, bc *BlockChain, txs []*Transaction, nodeID string) *Transaction {
	var txInputs []*TxInput   // 输入列表
	var txOutputs []*TxOutput // 输出列表
	// 调用可花费UTXO函数
	money, spendableUTXODic := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("money: %v\n", money)
	// 获取钱包集合对象
	wallets := NewWallets(nodeID)
	// 查找对应的钱包结构
	wallet := wallets.Wallets[from]
	// 输入
	for txHash, indexArray := range spendableUTXODic {
		txHashesBytes, err := hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string to []byte failed! %v\n", err)
		}
		// 遍历索引列表
		for _, index := range indexArray {
			txInpt := &TxInput{txHashesBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInpt)
		}
	}
	// 输出 (转账源)
	//txOutput := &TxOutput{amount, to}
	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)
	// 输出 (找零)
	if money > amount {
		//txOutput = &TxOutput{money - amount, from}
		txOutput = NewTxOutput(money-amount, from)
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足...\n")
	}
	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction() // 生成一笔完整的交易

	// 对交易进行签名
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

// IsCoinbaseTransaction 判断指定的交易是否是一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}

// Sign 交易签名
// prevTxs: 代表当前交易的输入所引用的所有OUTPUT所属的交易
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	// 处理输入, 保证交易的正确性
	// 检查tx中每一个输入所引用的交易哈希是否包含在prevTxs中
	// 如果没有包含在里面, 则说明该交易被人修改了
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("ERROR: Prev transaction is not correct!\n")
		}
	}
	// 提取需要签名的属性
	txCopy := tx.TrimmedCopy()
	// 处理交易副本的输入
	for vinID, vin := range txCopy.Vins {
		// 获取关联交易
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 找到发送者(当前输入引用的哈希--输出的哈希)
		txCopy.Vins[vinID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 生成交易副本的哈希
		txCopy.TxHash = txCopy.Hash()
		// 调用核心签名函数
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if nil != err {
			log.Panicf("sign to transaction [%x] failed! %v\n", txCopy.TxHash, err)
		}
		// 组成交易签名
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[vinID].Signature = signature
	}
}

// TrimmedCopy 交易拷贝, 生成一个专门用于交易签名的副本
func (tx *Transaction) TrimmedCopy() Transaction {
	// 重新组装生成一个新的交易
	var inputs []*TxInput
	var outputs []*TxOutput
	// 组装input
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TxInput{vin.TxHash, vin.Vout, nil, nil})
	}
	// 组装output
	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TxOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transaction{tx.TxHash, inputs, outputs}
	return txCopy
}

// Hash 设置用于签名的交易的哈希
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(tx.Serialize())
	return hash[:]
}

// Serialize 交易序列化
func (tx *Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	// 新建编码对象
	encoder := gob.NewEncoder(&buffer)
	// 编码(序列化)
	if err := encoder.Encode(tx); nil != err {
		log.Panicf("serialize the tx to []bytes failed! %v\n", err)
	}
	return buffer.Bytes()
}

// Verify 验证签名
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	// 检查能否找到交易哈希
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("VERIFY ERROR: transaction verify failed!\n")
		}
	}
	// 提取相同的交易签名属性
	txCopy := tx.TrimmedCopy()
	// 使用相同的椭圆
	curve := elliptic.P256()
	// 遍历tx输入, 对每笔输入所引用的输出进行验证
	for vinID, vin := range tx.Vins {
		// 获取关联交易
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 找到发送者(当前输入引用的哈希--输出的哈希)
		txCopy.Vins[vinID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 由需要验证的数据生成的交易哈希, 必须要和签名时的数据完全一致
		txCopy.TxHash = txCopy.Hash()
		// 在比特币中, 签名时一个数值对, r, s代表签名
		// 所以要从输入的signature中获取
		// 获取r, s. r, s长度相等
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		// 获取公钥
		// 公钥由X, Y坐标组成
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(pubKeyLen / 2)])
		y.SetBytes(vin.PublicKey[(pubKeyLen / 2):])
		// 调用验证签名核心函数
		rawPublicKey := ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(&rawPublicKey, txCopy.TxHash, &r, &s) {
			return false
		}
	}
	return true
}
