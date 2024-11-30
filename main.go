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
var lastPredictions []string  // 保存上一次的预测结果
var predictionIndex int       // 当前预测结果的索引

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
			if len(lastPredictions) == 0 {
				// 如果没有上次的预测结果，重新搜索
				lastPredictions = dictTrie.Search(builder.String())
				predictionIndex = 0
			} else {
				// 循环切换预测结果
				predictionIndex = (predictionIndex + 1) % len(lastPredictions)
			}
			if len(lastPredictions) > 0 {
				builder.Reset()
				builder.WriteString(lastPredictions[predictionIndex])
			}
		default:
			builder.WriteRune(char)
			// 显示预测结果
			lastPredictions = dictTrie.Search(builder.String())
			predictionIndex = 0  // 重置预测索引
			if len(lastPredictions) > 0 {
				// 最多显示3个预测结果
				displayPredictions := lastPredictions
				if len(displayPredictions) > 3 {
					displayPredictions = displayPredictions[:3]
				}
				// 清除当前行并显示输入
				fmt.Printf("\r\033[K\033[1;32mTrans-CLI>\033[0m %s", builder.String())
				// 清除下一行并显示预测
				fmt.Printf("\n\033[K预测: %s", strings.Join(displayPredictions, ", "))
				// 回到输入行
				fmt.Printf("\033[1A\r\033[1;32mTrans-CLI>\033[0m %s", builder.String())
			} else {
				lastPredictions = nil  // 清空预测结果
				// 如果没有预测结果，清除所有预测显示
				fmt.Printf("\r\033[K\033[1;32mTrans-CLI>\033[0m %s", builder.String())
				fmt.Printf("\n\033[K")  // 清除预测行
				fmt.Printf("\033[1A")   // 回到输入行
			}
		}

		// 清空预测行并重新显示当前输入
		fmt.Printf("\r\033[K")  // 清除当前行
		fmt.Printf("\n\033[K")  // 清除预测行
		fmt.Printf("\r\033[1A\033[1;32mTrans-CLI>\033[0m %s", builder.String())
	}

}
