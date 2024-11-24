package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/senorL/TransCLI/conf"
	"github.com/senorL/TransCLI/internal/transfer"
	"github.com/urfave/cli/v2"
)

func main() {
	// 实例化一个cli应用
	conf.GetConf()

	app := &cli.App{
		Name:     "TransCLI",
		Usage:    "终端翻译助手",
		Commands: []*cli.Command{},
		Action: func(c *cli.Context) error {
			transfer.TranForBaidu(c.Args().Slice())
			return nil
		},
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\033[1;32mTrans-CLI>\033[0m ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		args := strings.Split(input, " ")
		err := app.Run(append([]string{os.Args[0]}, args...))
		if err != nil {
			fmt.Println(err)
		}
	}
}
