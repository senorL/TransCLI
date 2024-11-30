package trie

// TrieNode 表示前缀树的节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	word     string
}

// Trie 前缀树结构
type Trie struct {
	root *TrieNode
}

// NewTrie 创建新的前缀树
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

// Insert 向前缀树中插入单词
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if _, exists := node.children[ch]; !exists {
			node.children[ch] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[ch]
	}
	node.isEnd = true
	node.word = word
}

// Search 搜索前缀，返回所有可能的补全词
func (t *Trie) Search(prefix string) []string {
	node := t.root
	for _, ch := range prefix {
		if _, exists := node.children[ch]; !exists {
			return nil
		}
		node = node.children[ch]
	}

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
