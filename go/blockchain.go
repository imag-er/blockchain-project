package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// 查看区块链信息的函数，返回包含区块链信息的字符串
func viewBlockchainInfo(username string) string {
	var blockchainInfo strings.Builder

	for _, block := range blockchain.Blocks {
		// 检查区块的用户名是否与当前登录用户一致
		if block.BlogEntries.Entries != nil && block.BlogEntries.Entries[0].Content == fmt.Sprintf("%s的博客本", username) {
			blockchainInfo.WriteString(fmt.Sprintf("Index: %d\n", block.Index))
			blockchainInfo.WriteString(fmt.Sprintf("Timestamp: %s\n", block.Timestamp.Format("2006-01-02 15:04:05")))
			blockchainInfo.WriteString(fmt.Sprintf("Hash: %s\n", block.Hash))
			blockchainInfo.WriteString(fmt.Sprintf("PrevHash: %s\n", block.PrevHash))
			blockchainInfo.WriteString(fmt.Sprintf("Nonce: %d\n", block.Nonce))
			blockchainInfo.WriteString("Blog Entries:\n")

			for _, entry := range block.BlogEntries.Entries {
				blockchainInfo.WriteString(fmt.Sprintf("[%s] %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.Content))
			}

			blockchainInfo.WriteString("--------------------------------------------------------------------------------------------------------------------\n")
		}
	}

	// 返回包含区块链信息的字符串
	return blockchainInfo.String()
}

// 处理验证区块链的请求
func validateBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	// 打印接收到的请求信息
	//log.Printf("Received request: %s %s\n", r.Method, r.URL)

	// 鉴别用户是否登录
	if !loggedIn {
		http.Error(w, "用户未登录", http.StatusUnauthorized)
		return
	}

	// 调用验证区块链是否被篡改的函数
	if isBlockchainTampered() {
		// 如果区块链被篡改，返回警告信息
		responseData := map[string]string{
			"status":  "warning",
			"message": "警告：区块链数据已被篡改！",
		}

		// 将结构体转换为 JSON 格式的字符串
		responseJSON, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "无法序列化响应数据", http.StatusInternalServerError)
			return
		}

		// 打印要发送到前端的响应信息
		//log.Printf("Sent response (length=%d): %q\n", len(responseJSON), responseJSON)
		// 设置响应头为 application/json
		w.Header().Set("Content-Type", "application/json")
		// 发送到前端
		w.Write(responseJSON)
	} else {
		// 如果区块链未被篡改，返回成功信息
		responseData := map[string]string{
			"status":  "success",
			"message": "区块链未被篡改",
		}

		// 将结构体转换为 JSON 格式的字符串
		responseJSON, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "无法序列化响应数据", http.StatusInternalServerError)
			return
		}

		// 打印要发送到前端的响应信息
		//log.Printf("Sent response (length=%d): %q\n", len(responseJSON), responseJSON)
		// 设置响应头为 application/json
		w.Header().Set("Content-Type", "application/json")
		// 发送到前端
		w.Write(responseJSON)
	}
}

// 加载区块链信息和博客
func loadBlockchain(username string) {
	// 从文件中加载区块链信息
	blockchainFileName := getBlockchainFileName(username)
	blockchainData, err := os.ReadFile(blockchainFileName)
	if err != nil {
		// 如果文件不存在，创建一个新文件
		if os.IsNotExist(err) {
			createGenesisBlock(username)
			return
		}
		fmt.Println("无法加载区块链文件:", err)
		return
	}

	err = json.Unmarshal(blockchainData, &blockchain)
	if err != nil {
		fmt.Println("无法解析区块链文件:", err)
		return
	}

	// 从文件中加载博客信息
	blogFileName := getBlogFileName(username)
	blogData, err := os.ReadFile(blogFileName)
	if err != nil {
		fmt.Println("无法加载博客文件:", err)
		return
	}

	err = json.Unmarshal(blogData, &getCurrentBlock().BlogEntries)
	if err != nil {
		fmt.Println("无法解析博客文件:", err)
		return
	}
}

// 验证区块链是否被篡改
func isBlockchainTampered() bool {
	for i := 1; i < len(blockchain.Blocks); i++ {
		prevBlock := blockchain.Blocks[i-1]
		currentBlock := blockchain.Blocks[i]

		// 检查当前区块的前一个哈希是否等于前一个区块的哈希
		if currentBlock.PrevHash != prevBlock.Hash {
			return true
		}
	}
	return false
}

func mineBlock(entry BlogEntry) *Block {
	// 挖矿
	currentBlock := getCurrentBlock()
	newBlock := mine(currentBlock, entry)

	// 更新区块链
	updateBlockchain()

	return newBlock
}

func mine(prevBlock *Block, entry BlogEntry) *Block {
	// 工作量证明（增加挖矿次数）
	maxAttempts := 1000000
	for nonce := 0; nonce < maxAttempts; nonce++ {
		hash := calculateHash(*prevBlock, nonce)
		if hashIsValid(hash) {
			// 将原有的博客条目也包含在 BlogEntries 中
			newBlogEntries := append(prevBlock.BlogEntries.Entries, entry)
			return &Block{
				Index:       prevBlock.Index + 1,
				Timestamp:   time.Now(),
				BlogEntries: BlogBook{Entries: newBlogEntries},
				PrevHash:    prevBlock.Hash,
				Hash:        hash,
				Nonce:       nonce,
			}
		}
	}
	// 如果未能找到符合条件的哈希值，可以考虑返回错误或者进行其他处理
	return nil
}

func hashIsValid(hash string) bool {
	// 简化版的验证条件
	return strings.HasPrefix(hash, "000")
}

func calculateBlogHash(blog BlogBook) string {
	data, err := json.Marshal(blog)
	if err != nil {
		fmt.Println("计算哈希值时发生错误:", err)
		return ""
	}

	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func calculateHash(block Block, nonce int) string {
	data := fmt.Sprintf("%d%d%s%d%s", block.Index, block.Timestamp.UnixNano(), block.PrevHash, nonce, calculateBlogHash(block.BlogEntries))
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func getCurrentBlock() *Block {
	// 获取当前区块
	return &blockchain.Blocks[len(blockchain.Blocks)-1]
}

func updateBlockchain() {
	// 将区块链和博客写入文件
	blockchainData, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		fmt.Println("更新区块链时发生错误:", err)
		return
	}

	err = os.WriteFile(getBlockchainFileName(currentUser.Username), blockchainData, 0755)
	if err != nil {
		fmt.Println("写入区块链文件时发生错误:", err)
		return
	}

	blogData, err := json.MarshalIndent(getCurrentBlock().BlogEntries, "", "  ")
	if err != nil {
		fmt.Println("更新博客时发生错误:", err)
		return
	}

	err = os.WriteFile(getBlogFileName(currentUser.Username), blogData, 0755)
	if err != nil {
		fmt.Println("写博客时发生错误:", err)
		return
	}
}

// 创世区块
func createGenesisBlock(username string) Block {
	// 检查文件是否已存在，如果存在则不再创建
	blockchainFileName := getBlockchainFileName(username)
	blogFileName := getBlogFileName(username)

	if _, err := os.Stat(blockchainFileName); err == nil {
		fmt.Println("区块链文件已存在，无需再次创建。")
		return Block{}
	}

	if _, err := os.Stat(blogFileName); err == nil {
		fmt.Println("博客文件已存在，无需再次创建。")
		return Block{}
	}

	genesisBlock := mineGenesisBlock(username)
	writeBlockToFile(genesisBlock, username)

	// 将创世区块添加到全局的 blockchain 变量中
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)

	return genesisBlock
}

func mineGenesisBlock(username string) Block {
	// 挖矿创造创世区块
	nonce := 0
	for {
		genesisBlock := Block{
			Index:       0,
			Timestamp:   time.Now(),
			BlogEntries: BlogBook{Entries: []BlogEntry{{Timestamp: time.Now(), Content: fmt.Sprintf("%s的博客本", username)}}},
			PrevHash:    "",
			Nonce:       nonce,
		}

		// 计算哈希值
		hash := calculateHash(genesisBlock, nonce)

		// 检查哈希值是否符合条件
		if hashIsValid(hash) {
			// 符合条件的哈希值找到，返回创世区块
			genesisBlock.Hash = hash
			return genesisBlock
		}

		// 尝试下一个 nonce
		nonce++
	}
}

// 写入区块和博客到文件的函数
func writeBlockToFile(block Block, username string) {
	blockchainData, err := json.MarshalIndent(Blockchain{Blocks: []Block{block}}, "", "  ")
	if err != nil {
		fmt.Println("写入区块到文件时发生错误:", err)
		return
	}

	err = os.WriteFile(getBlockchainFileName(username), blockchainData, 0644)
	if err != nil {
		fmt.Println("写入区块链文件时发生错误:", err)
		return
	}

	blogData, err := json.MarshalIndent(block.BlogEntries, "", "  ")
	if err != nil {
		fmt.Println("写入博客到文件时发生错误:", err)
		return
	}

	err = os.WriteFile(getBlogFileName(username), blogData, 0644)
	if err != nil {
		fmt.Println("写入博客文件时发生错误:", err)
		return
	}
}

func getBlogFileName(username string) string {
	return blogFolder + username + "_blog.json"
}

func getBlockchainFileName(username string) string {
	return blockchainFolder + username + "_blockchain.json"
}
