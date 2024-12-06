package prediction

import (
	"strings"
)

var trie = NewTrie()
var predict string

func Predict(builder *strings.Builder) string {
	words := strings.Fields(builder.String())
	if len(words) == 0 {
		predict = ""
		return ""
	}
	prefix := words[len(words)-1]
	if len(prefix) < 3 {
		predict = ""
		return ""
	}
	predictions := trie.Search(prefix)
	if len(predictions) > 0 {
		predict = predictions[0][len(prefix):]
	} else {
		predict = ""
	}
	return predict
}

func KeyTab(builder *strings.Builder) {
	builder.WriteString(predict)
}

func LoadDict(dict string) {
	words := strings.Fields(dict)
	for _, word := range words {
		trie.Insert(word)
	}
}
