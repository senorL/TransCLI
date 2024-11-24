package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/senorL/TransCLI/conf"
	"github.com/senorL/TransCLI/translate"
)

func main() {
	// 实例化一个cli应用
	conf.GetConf()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\033[1;32mTrans-CLI>\033[0m ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.TrimRight(input, "\r\n")

		if input == "exit" {
			break
		}
		translate.TranForBaidu(input)
	}
}
