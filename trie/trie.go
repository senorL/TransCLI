package trie

// TrieNode 表示压缩前缀树的节点
type TrieNode struct {
	prefix   string           // 存储压缩的前缀
	children map[rune]*TrieNode
	isEnd    bool
	word     string
}

// Trie 压缩前缀树结构
type Trie struct {
	root *TrieNode
}

// NewTrie 创建新的压缩前缀树
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			prefix:   "",
			children: make(map[rune]*TrieNode),
		},
	}
}

// findCommonPrefix 找到两个字符串的最长公共前缀
func findCommonPrefix(s1, s2 string) string {
	i := 0
	for i < len(s1) && i < len(s2) && s1[i] == s2[i] {
		i++
	}
	return s1[:i]
}

// Insert 向压缩前缀树中插入单词
func (t *Trie) Insert(word string) {
	if word == "" {
		return
	}

	node := t.root
	current := word

	for len(current) > 0 {
		firstChar := rune(current[0])
		child, exists := node.children[firstChar]

		if !exists {
			// 创建新节点存储剩余部分
			newNode := &TrieNode{
				prefix:   current,
				children: make(map[rune]*TrieNode),
				isEnd:    true,
				word:     word,
			}
			node.children[firstChar] = newNode
			return
		}

		// 找到当前节点prefix和待插入单词的最长公共前缀
		commonPrefix := findCommonPrefix(child.prefix, current)
		if commonPrefix == child.prefix {
			// 当前节点的prefix是待插入单词的前缀，继续向下插入剩余部分
			current = current[len(commonPrefix):]
			node = child
			continue
		}

		// 需要分裂当前节点
		newNode := &TrieNode{
			prefix:   commonPrefix,
			children: make(map[rune]*TrieNode),
		}

		// 更新原节点
		child.prefix = child.prefix[len(commonPrefix):]
		newNode.children[rune(child.prefix[0])] = child

		// 插入剩余部分
		remainingPart := current[len(commonPrefix):]
		if remainingPart != "" {
			newNode.children[rune(remainingPart[0])] = &TrieNode{
				prefix:   remainingPart,
				children: make(map[rune]*TrieNode),
				isEnd:    true,
				word:     word,
			}
		} else {
			newNode.isEnd = true
			newNode.word = word
		}

		node.children[firstChar] = newNode
		return
	}

	node.isEnd = true
	node.word = word
}

// Search 搜索前缀，返回所有可能的补全词
func (t *Trie) Search(prefix string) []string {
	if prefix == "" {
		return nil
	}

	node := t.root
	current := prefix

	// 查找前缀
	for len(current) > 0 {
		firstChar := rune(current[0])
		child, exists := node.children[firstChar]
		if !exists {
			return nil
		}

		if len(current) < len(child.prefix) {
			// 检查child.prefix是否以current开头
			if child.prefix[:len(current)] != current {
				return nil
			}
			current = ""
		} else {
			// 检查child.prefix是否匹配current的开头部分
			if current[:len(child.prefix)] != child.prefix {
				return nil
			}
			current = current[len(child.prefix):]
		}
		node = child
	}

	// 收集所有可能的补全词
	var results []string
	t.collectWords(node, &results)
	return results
}

// collectWords 收集从当前节点开始的所有完整单词
func (t *Trie) collectWords(node *TrieNode, results *[]string) {
	if node.isEnd {
		*results = append(*results, node.word)
	}

	for _, child := range node.children {
		t.collectWords(child, results)
	}
}
