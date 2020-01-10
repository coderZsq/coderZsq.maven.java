package main

import (
	"fmt"

	"./blc"
	"github.com/boltdb/bolt"
)

// 启动
func main() {
	bc := blc.CreateBlockChainWithGenesisBlock()
	//fmt.Printf("blockchain: %v\n", bc.Blocks[0])
	bc.AddBlock([]byte("a send 100 eth to b"))
	bc.AddBlock([]byte("b send 100 eth to c"))
	// 上链
	// AddBlock(height int64, prevBlockHash []byte, data []byte)
	//bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1, bc.Blocks[len(bc.Blocks)-1].Hash, []byte("alice send 10 btc to bob"))
	//bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1, bc.Blocks[len(bc.Blocks)-1].Hash, []byte("bob send 5 btc to tony"))
	//for _, block := range bc.Blocks {
	//	fmt.Printf("prevBlockHash: %x, currentHash: %x\n", block.PrevBlockHash, block.Hash)
	//}
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if nil != b {
			hash := b.Get([]byte("l"))
			fmt.Printf("value : %x\n", hash)
		}
		return nil
	})
}
