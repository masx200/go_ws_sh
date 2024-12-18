package go_ws_sh

import (
	"io"
	"os"
	"os/exec"
)

func pipe_std() {
	cmd := exec.Command("pwsh")

	// 设置标准输入、输出和错误流
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// 处理标准输出和错误流
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	// 处理标准输入流
	io.Copy(stdin, os.Stdin)

	// 等待命令执行完毕
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
