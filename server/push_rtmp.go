package server

import (
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// CommandStatus 用于封装命令执行的状态
type CommandStatus struct {
	cmd       *exec.Cmd
	Running   bool
	Success   bool
	Timestamp time.Time
}

func NewPushRtmp(remoteAddr string) *CommandStatus {
	cmd := exec.Command("ffmpeg",
		"-f", "h264",
		"-i", "udp://127.0.0.1:65515",
		"-vcodec", "copy",
		"-an", // 这个参数用于禁用音频
		"-f", "flv",
		remoteAddr)
	// 定义ffmpeg命令
	return &CommandStatus{
		cmd:     cmd,
		Running: false,
	}
}

func (c *CommandStatus) PushRtmp(group *sync.WaitGroup, done chan CommandStatus) {
	defer group.Done()
	// 输出命令详情（可选）
	fmt.Printf("Executing command: %s\n", c.cmd.String())
	// 启动命令
	err := c.cmd.Start()
	if err != nil {
		done <- CommandStatus{Running: true, Success: false, Timestamp: time.Now()}
		return
	}
	// 等待命令完成
	err = c.cmd.Wait()
	if err != nil {
		done <- CommandStatus{Running: true, Success: false, Timestamp: time.Now()}
		return
	}
	// 如果命令成功完成
	done <- CommandStatus{Running: true, Success: true, Timestamp: time.Now()}
}
