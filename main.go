package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/eiannone/keyboard"

	"github.com/senorL/TransCLI/conf"
	"github.com/senorL/TransCLI/history"
	"github.com/senorL/TransCLI/prediction"
	"github.com/senorL/TransCLI/translate"
)

//go:embed asset/dict.txt
var dict string

func main() {
	conf.GetConf()
	prediction.LoadDict(dict)
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
		case keyboard.KeyEnter:
			fmt.Printf("\n\033[1;32mTrans-CLI>\033[0m ") // Trans-CLI>
			translate.TranForBaidu(builder.String())     // 翻译
			history.AddHistory(builder.String())         // 历史记录
			builder.Reset()                              // Trans-CLI>
		case keyboard.KeyArrowDown, keyboard.KeyArrowUp:
			if key == keyboard.KeyArrowUp {
				history.GetUpHistory(&builder)
			}
			if key == keyboard.KeyArrowDown {
				history.GetDownHistory(&builder)
			}
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			str := builder.String()
			if len(str) > 0 {
				builder.Reset()
				builder.WriteString(str[:len(str)-1])
			}
		case keyboard.KeyTab:
			prediction.KeyTab(&builder)
		case keyboard.KeySpace:
			builder.WriteRune(' ')
		default:
			builder.WriteRune(char)
		}
		fmt.Printf("\r\033[1;32mTrans-CLI>\033[0m %s\033[0;90m%s\033[0m\033[K", builder.String(), prediction.Predict(&builder))
	}

}
