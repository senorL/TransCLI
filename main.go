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
var lastPredictions []string // 保存上一次的预测结果
var predictionIndex int      // 当前预测结果的索引

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
			words := strings.Fields(builder.String())
			if len(words) == 0 {
				continue
			}

			// 只对最后一个单词进行预测
			prefix := words[len(words)-1]
			if len(lastPredictions) == 0 {
				lastPredictions = dictTrie.Search(prefix)
				predictionIndex = 0
			} else {
				predictionIndex = (predictionIndex + 1) % len(lastPredictions)
			}

			if len(lastPredictions) > 0 {
				// 保留之前的单词，只替换最后一个单词
				words[len(words)-1] = lastPredictions[predictionIndex]
				builder.Reset()
				builder.WriteString(strings.Join(words, " "))
			}

		default:
			builder.WriteRune(char)
			// 显示预测结果
			words := strings.Fields(builder.String())
			if len(words) > 0 {
				// 只对最后一个单词进行预测
				prefix := words[len(words)-1]
				lastPredictions = dictTrie.Search(prefix)
			} else {
				lastPredictions = nil
			}

			predictionIndex = 0 // 重置预测索引
			if len(lastPredictions) > 0 {
				displayPredictions := lastPredictions
				if len(displayPredictions) > 3 {
					displayPredictions = displayPredictions[:3]
				}

				// 显示当前输入和预测结果
				fmt.Printf("\r\033[K\033[1;32mTrans-CLI>\033[0m %s", builder.String())
				fmt.Printf("\n\033[K预测: %s", formatPredictions(displayPredictions))
				fmt.Printf("\033[1A\r\033[1;32mTrans-CLI>\033[0m %s", builder.String())
			} else {
				// 清除预测显示
				fmt.Printf("\r\033[K\033[1;32mTrans-CLI>\033[0m %s", builder.String())
				fmt.Printf("\n\033[K")
				fmt.Printf("\033[1A")
			}
		}
	}

}

func formatPredictions(predictions []string) string {
	var result strings.Builder
	for i, prediction := range predictions {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(prediction)
	}
	return result.String()
}
