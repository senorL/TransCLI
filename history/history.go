package history

import (
	"container/list"
	"strings"
)

var history = list.New()
var current *list.Element

func AddHistory(item string) {
	history.PushBack(item)
	current = history.Back()
}

func GetUpHistory(builder *strings.Builder) {
	if current != nil && current.Prev() != nil {
		current = current.Prev()
	}
	builder.Reset()
	builder.WriteString(current.Value.(string))
}

func GetDownHistory(builder *strings.Builder) {
	if current != nil && current.Next() != nil {
		current = current.Next()
	}
	builder.Reset()
	builder.WriteString(current.Value.(string))
}
