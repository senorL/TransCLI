package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/eiannone/keyboard"

	"github.com/senorL/TransCLI/conf"
	"github.com/senorL/TransCLI/translate"
	"github.com/senorL/TransCLI/trie"
)

//go:embed asset/dict.txt
var dict string

var history []string
var index int
var dictTrie *trie.Trie

func main() {
	conf.GetConf()

	dictTrie = trie.NewTrie()
	for _, word := range strings.Split(dict, "\n") {
		word = strings.TrimSpace(word)
		if word != "" {
			dictTrie.Insert(word)
		}
	}

	var builder strings.Builder
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	fmt.Printf("\r\033[1;32mTrans-CLI>\033[0m %s", builder.String())
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case keyboard.KeyEsc:
			return // 按下 ESC 键退出程序
		case keyboard.KeyEnter:
			fmt.Printf("\n\033[1;32mTrans-CLI>\033[0m ") // 换行
			translate.TranForBaidu(builder.String())     // 调用翻译功能
			history = append(history, builder.String())  // 记录历史
			index = len(history) - 1                     // 更新历史记录索引
			builder.Reset()                              // 清空输入内容
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			// 删除最后一个字符
			str := builder.String()
			if len(str) > 0 {
				builder.Reset()
				builder.WriteString(str[:len(str)-1])
			}
		case keyboard.KeySpace:
			builder.WriteRune(' ')
		case keyboard.KeyArrowDown, keyboard.KeyArrowUp, keyboard.KeyArrowLeft, keyboard.KeyArrowRight:
			// TODO: 处理方向键
			// 方向键控制history记录，可以进行上下翻页
			if key == keyboard.KeyArrowDown {
				if index < len(history)-1 {
					index++
					builder.Reset()
					builder.WriteString(history[index])
				}
			} else if key == keyboard.KeyArrowUp {
				if index > 0 {
					index--
					builder.Reset()
					builder.WriteString(history[index])
				}
			}
		case keyboard.KeyTab:
			// 处理预测功能
			predictions := dictTrie.Search(builder.String())
			if len(predictions) > 0 {
				// 使用第一个预测结果
				builder.Reset()
				builder.WriteString(predictions[0])
			}
		default:
			builder.WriteRune(char)
			// 显示预测结果
			predictions := dictTrie.Search(builder.String())
			if len(predictions) > 0 {
				// 最多显示3个预测结果
				if len(predictions) > 3 {
					predictions = predictions[:3]
				}
				fmt.Printf("\r\033[1;32mTrans-CLI>\033[0m %s\033[K", builder.String())
				fmt.Printf("\n预测: %s", strings.Join(predictions, ", "))
				fmt.Printf("\r\033[1A\033[1;32mTrans-CLI>\033[0m %s", builder.String())
			} else {
				fmt.Printf("\r\033[1;32mTrans-CLI>\033[0m %s\033[K", builder.String())
			}
		}

		// 清空预测行并重新显示当前输入
		fmt.Printf("\r\033[K")  // 清除当前行
		fmt.Printf("\n\033[K")  // 清除预测行
		fmt.Printf("\r\033[1A\033[1;32mTrans-CLI>\033[0m %s", builder.String())
	}

}
