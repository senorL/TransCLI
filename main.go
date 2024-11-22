package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/senorL/TransCLI/conf"
	"github.com/senorL/TransCLI/translate"
)

func main() {
	conf.GetConf()

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
		default:
			builder.WriteRune(char)
		}

		// 使用 ANSI 转义序列清空当前行，然后重新打印内容
		fmt.Printf("\r\033[1;32mTrans-CLI>\033[0m %s\033[K", builder.String())
	}
}
