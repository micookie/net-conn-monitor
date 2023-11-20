package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

const (
	pingTimeout            = 10 * time.Millisecond
	pingInterval           = 10 * time.Second
	maxConsecutiveFailures = 3
)

var (
	targets       = []string{"qq.com", "114.114.114.114", "baidu.com"}
	failureCount  = 0
	failureStart  time.Time
	monitoringURL = "http://www.pushplus.plus/send"
	token         = os.Getenv("PUSH_TOKEN")
)

func main() {
	sendNotification("网络连接监控程序已启动", time.Now().Format(time.DateTime))

	for {
		success := pingTargets()

		if success {
			if failureCount > 0 {
				duration := time.Since(failureStart).Round(time.Second)
				sendNotification("网络故障恢复", fmt.Sprintf("网络已恢复<br>断开时间：%s<br>恢复时间：%s <br>持续时长：%s", failureStart.Format(time.DateTime), time.Now().Format(time.DateTime), formatDuration(duration)))
				failureCount = 0
			}
		} else {
			if failureCount == 0 {
				failureStart = time.Now()
			}
			failureCount++

			if failureCount == maxConsecutiveFailures {
				sendNotification("网络故障", "网络故障，请及时处理。")
			}
		}

		time.Sleep(pingInterval)
	}
}

func pingTargets() bool {
	for _, target := range targets {
		pinger, err := ping.NewPinger(target)
		pinger.SetPrivileged(true)
		if err != nil {
			log.Printf("创建 Pinger 失败：%s\n", err)
			return false
		}

		pinger.Count = 1
		pinger.Timeout = pingTimeout

		err = pinger.Run()
		if err != nil {
			log.Printf("%s ping失败：%s\n", target, err)
			return false
		}

		stats := pinger.Statistics()
		if stats.PacketsRecv > 0 {
			log.Printf("%s ping成功\n", target)
		} else {
			log.Printf("%s ping超时\n", target)
			return false
		}
	}
	return true
}

func sendNotification(title, content string) {
	message := fmt.Sprintf(`{"token":"%s","title":"%s","content":"%s"}`, token, title, content)

	resp, err := http.Post(monitoringURL, "application/json", strings.NewReader(message))
	if err != nil {
		log.Println("发送消息失败：", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取响应失败：", err)
		return
	}

	log.Println("消息发送成功")
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}
