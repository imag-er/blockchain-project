package main

import (
	"fmt"
	"net/http"
	"time"
	"strings"
)

// 提交博客入口
func submitBlogEntry(content string, w http.ResponseWriter) {
	// 创建博客条目
	entry := BlogEntry{
		Timestamp: time.Now(),
		Content:   content,
	}

	// 记录挖矿开始时间
	startTime := time.Now()

	// 挖矿创建新区块
	newBlock := mineBlock(entry)

	// 计算挖矿耗时
	elapsedTime := time.Since(startTime)

	// 添加到区块链
	blockchain.Blocks = append(blockchain.Blocks, *newBlock)

	// 更新区块链
	updateBlockchain()

	// 验证是否被篡改
	checkBlogIntegrity(w)

	// 在写博客结束后返回添加成功以及挖矿信息
	response := fmt.Sprintf("博客添加成功！挖矿耗时：%s\n", elapsedTime)
	response += fmt.Sprintf("新区块信息：Index: %d, Timestamp: %s, Hash: %s, Nonce: %d\n",
		newBlock.Index, newBlock.Timestamp.Format("2006-01-02 15:04:05"), newBlock.Hash, newBlock.Nonce)

	w.Write([]byte(response))
}

// 验证博客的完整性
func checkBlogIntegrity(w http.ResponseWriter) {
	if isBlockchainTampered() {
		w.Write([]byte("警告：区块链数据已被篡改！"))
		blogTampered = true
	} else {
		w.Write([]byte("区块链未被篡改"))
	}
}

// 查看所有博客
func viewBlogEntries(username string) string {
	seenEntries := make(map[string]struct{})
	var Bloginfo strings.Builder

	for _, block := range blockchain.Blocks {
		// 检查区块的用户名是否与当前登录用户一致
		if block.BlogEntries.Entries != nil && block.BlogEntries.Entries[0].Content == fmt.Sprintf("%s的博客本", username) {
			for _, entry := range block.BlogEntries.Entries {
				// 去重逻辑
				entryKey := entry.Timestamp.Format("2006-01-02 15:04:05") + entry.Content
				if _, seen := seenEntries[entryKey]; !seen {
					Bloginfo.WriteString(fmt.Sprintf("[%s] %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.Content))
					seenEntries[entryKey] = struct{}{}
				}
			}
		}
	}

	// 返回包含博客信息的字符串
	return Bloginfo.String()
}
