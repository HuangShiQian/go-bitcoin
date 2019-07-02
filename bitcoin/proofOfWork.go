package main

import (
	"math/big"
	"fmt"
	"bitcoin/utils"
	"bytes"
	"crypto/sha256"
)

//实现挖矿功能 pow

// 字段：
// ​	区块：blockPkg
// ​	目标值：target
// 方法：
// ​	run计算
// ​	功能：找到nonce，从而满足哈希币目标值小

type ProofOfWork struct {
	// ​区块：blockPkg
	block *Block
	// ​目标值：target，这个目标值要与生成哈希值比较
	target *big.Int//结构，提供了方法：比较，把哈希值设置为big.Int类型

}

//创建ProofOfWork
//block由用户提供
//target目标值由系统提供
func NewProofOfWork(block *Block)*ProofOfWork  {
	pow:=ProofOfWork{
		block:block,// 注意如果这边只写一个，前面一定要注明block，如果没注明，必须要把block和target对应的都写上
	}

	//难度值先写死，不去推导，后面补充推导方式
	targetStr:="0001000000000000000000000000000000000000000000000000000000000000"
	tmpBigInt:=new(big.Int)

	//将我们的难度值赋值给bigint
	tmpBigInt.SetString(targetStr,16)

	pow.target=tmpBigInt
	return &pow

}

//挖矿函数，不断变化nonce，使得sha256(数据+nonce) < 难度值
//返回：区块哈希，nonce
func (pow *ProofOfWork)Run()([]byte,uint64)  {
	//定义随机数
	var nonce uint64
	var hash [32]byte
	fmt.Println("开始挖矿...")

	for {
		fmt.Printf("%x\r", hash[:])
		// 1. 拼接字符串 + nonce
		data:=pow.PrepareData(nonce)
		// 2. 哈希值 = sha256(data)
		hash=sha256.Sum256(data)

		//将hash转换为bigInt类型
		tmpInt:=new(big.Int)
		tmpInt.SetBytes(hash[:])

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//当前计算的哈希.Cmp(难度值)
		if tmpInt.Cmp(pow.target)==-1{
			fmt.Printf("挖矿成功,hash :%x, nonce :%d\n", hash[:], nonce)
			break
		}else {
			//如果不小于难度值
			nonce++
		}
	}
	// 	return 哈希，nonce
	return hash[:],nonce
}


//拼接nonce和block数据
func (pow *ProofOfWork)PrepareData(nonce uint64)[]byte  {
	b:=pow.block
	tmp := [][]byte{
		utils.UintToByte(b.Version), //将uint64转换为[]byte
		b.Prehash,
		b.MerkleRoot,
		utils.UintToByte(b.TimeStamp),
		utils.UintToByte(b.Bits),
		utils.UintToByte(nonce),
		//utils.UintToByte(b.nonce),
		//b.Hash,   它不应该参与哈希运算 而应该是我们填进去的
		//b.Data,
	}
	//使用join方法，将二维切片转为1维切片
	data:=bytes.Join(tmp,[]byte{})

	return data

}

func (pow *ProofOfWork)IsValid()bool  {
	// 	1. 获取区块
	// 2. 拼装数据（blockPkg + nonce）
	data:=pow.PrepareData(pow.block.Nonce)
	// 3. 计算sha256
	hash:=sha256.Sum256(data)
	// 4. 与难度值比较
	tmpInt:=new(big.Int)
	tmpInt.SetBytes(hash[:])
	// if tmpInt.Cmp(pow.target) == -1 {
	// 	return true
	// }
	// return false

	//满足条件，返回true
	return tmpInt.Cmp(pow.target)==-1
}