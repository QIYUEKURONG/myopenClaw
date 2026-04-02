package main

import (
	"bufio"
	"context"
	"fmt"
	"myopenclaw/agent"
	"myopenclaw/gateway"
	"myopenclaw/types"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // ← 添加这行
	if err != nil {
		fmt.Println("警告: 未找到 .env 文件")
	}
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ 错误: DEEPSEEK_API_KEY 未设置！")
		return
	}

	fmt.Println("MyOpenClaw 启动中......")

	gw := gateway.NewGateway()
	runtime := agent.NewRuntime()
	gw.RunTime = runtime

	//消息接收
	scanner := bufio.NewScanner(os.Stdin)
	uid := "Cli-User-Id" + uuid.New().String()

	for {
		fmt.Print(">>>>>>>>你: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			time.Sleep(time.Second * 1)
			continue
		}

		if input == "exit" {
			break
		}

		message := types.Message{
			ID:          uuid.New().String(),
			Content:     input,
			SessionID:   "",
			UserID:      uid,
			Channel:     "iTerm",
			CreatedTime: time.Now(),
		}

		ctx := context.Background()
		Response, err := gw.HandleMessage(ctx, &message)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("\nAgent: %s\n", Response.Content)
	}

}
