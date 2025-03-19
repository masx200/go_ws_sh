package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/masx200/go_ws_sh/go_ws_sh"
)

func main() {
	// 定义一个字符串变量来存储config参数的值
	var config string
	// var reconnect_delay int
	var mode string
	// var serverip string
	// 使用flag.String定义一个名为config的命令行参数，该参数有一个默认值和一个帮助说明
	flag.StringVar(&config, "config", "", "the configuration file")
	flag.StringVar(&mode, "mode", "", "server or client mode")
	// flag.IntVar(&reconnect_delay, "reconnect_delay", 0, "client reconnect delay seconds")
	// 打印用法信息并检查参数是否有效

	// flag.StringVar(&serverip, "serverip", "", "server ip")
	flag.Parse() // 解析命令行参数
	if mode == "" {
		log.Fatal("No mode provided.")
	}
	if config == "" {
		log.Fatal("No config provided.")
	}
	fmt.Printf("mode : %s\n", mode)
	// 输出config参数的值
	fmt.Printf("Config : %s\n", config)
	if mode == "server" {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in panic", r)
			}
		}()
		go_ws_sh.Server_start(config)
	} else {
		if mode == "client" {
			// for {
			// 	if reconnect_delay > 0 {
			// 		fmt.Printf("reconnect_delay : %d\n", reconnect_delay)
			//recover from panic
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in panic", r)
				}
			}()
			go_ws_sh.Client_start(config/* , serverip */)

			// 	//sleep time delay
			// 	time.Sleep(time.Duration(reconnect_delay) * time.Second)

			// } else {
			// 	return
			// }
			// }

		}
	}
}
