package main

import (

	"fmt"
)

func (cli *CLI)AddBlock(data string)  {
	//bc,_:=blockPkg.GetBlockChainInstance()
	//
	//err:=bc.AddBlock(data)
	//if err!=nil{
	//	fmt.Println("AddBlock Failed!",err)
	//	return
	//}
	//fmt.Println("添加区块成功！")
}

func (cli *CLI)CreateBlockChain(address string)  {
	if !isValidAddress(address){
		fmt.Println("传入地址无效，无效地址为:",address)
		return
	}

	err:=CreateBlockChain(address)
	if err!=nil{
		fmt.Println("CreateBlockChain Failed!",err)
		return
	}
	fmt.Println("创建区块链成功！")
}

func (cli *CLI)Print()  {

	bc,_:=GetBlockChainInstance()

	//调用迭代器 返回blockchain
	it:=bc.NewIterator()
	for {
		// 调用Next 方法 获取区块 游标左移
		block := it.Next()

		fmt.Printf("\n+++++++++++++++++++++++++++\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevHash : %x\n", block.Prehash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Bits : %d\n", block.Bits)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Transactions[0].TXInputs[0].ScriptSig)//矿工写入的数据

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		// 退出条件
		if block.Prehash == nil {
			fmt.Println("区块链遍历结束")
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	if !isValidAddress(address) {
		fmt.Println("传入地址无效，无效地址为:", address)
		return
	}
	bc, err := GetBlockChainInstance()
	if err!=nil{
		fmt.Println("getBalance err:",err)
		return
	}

 	defer bc.Db.Close()
	//得到查询地址的公钥哈希
	pubKeyHash := getPubKeyHashFromAddress(address)

	//获取所有相关的utxo集合
	utxos := bc.FindMyUTXO(pubKeyHash)
	total := 0.0

	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("'%s'的金额为:%f\n", address, total)
}

func (cli *CLI) send(from, to string, amount float64, miner, data string) {
	//校验from
	if !isValidAddress(from) {
		fmt.Println("传入from地址无效，无效地址为:", from)
		return
	}
	//校验to
	if !isValidAddress(to) {
		fmt.Println("传入to地址无效，无效地址为:", to)
		return
	}
	//校验miner
	if !isValidAddress(miner) {
		fmt.Println("传入miner地址无效，无效地址为:", miner)
		return
	}
	fmt.Println("from:", from)
	fmt.Println("to:", to)
	fmt.Println("amount:", amount)
	fmt.Println("miner:", miner)
	fmt.Println("data:", data)


	bc,err:= GetBlockChainInstance()
	if err!=nil{
		fmt.Println("send err",err)
		return
	}

	defer bc.Db.Close()

	//创建挖矿交易
	coinbaseTx := NewCoinbaseTx(miner, data)

	//创建txs数组，将有效交易添加进来
	txs := []*Transaction{coinbaseTx}

	//创建普通交易
	tx := NewTransaction(from, to, amount, bc)
	if tx != nil {
		fmt.Println("找到一笔有效的转账交易!")
		txs = append(txs, tx)
	} else {
		fmt.Println("注意，找到一笔无效的转账交易, 不添加到区块!")
	}

	//调用AddBlock
	err = bc.AddBlock(txs)
	if err != nil {
		fmt.Println("添加区块失败，转账失败!")
	}

	fmt.Println("添加区块成功，转账成功!")
}

func (cli *CLI)createWallet()  {
	wm:=NewWalletManager()
	if wm==nil{
		fmt.Println("createWallet失败!")
		return
	}

	address:=wm.createWallet()

	if len(address)==0{
		fmt.Println("创建钱包失败！")
		return
	}

	fmt.Println("新钱包地址为:",address)
}

func (cli *CLI)listAddress()  {
	wm:=NewWalletManager()
	if wm==nil{
		fmt.Println("NewWalletManager 失败!")
		return
	}

	addresses := wm.listAddresses()
	for _, address := range addresses {
		fmt.Printf("%s\n", address)
	}
}



func (cli *CLI) printTx() {
	bc, err := GetBlockChainInstance()
	if err != nil {
		fmt.Println("getBalance err:", err)
		return
	}

	defer bc.Db.Close()

	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Println("\n+++++++++++++++++ 区块分割 +++++++++++++++")

		for _, tx := range block.Transactions {
			//直接打印交易
			fmt.Println(tx)
		}

		if len(block.Prehash) == 0 {
			break
		}
	}
}
