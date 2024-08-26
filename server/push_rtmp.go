package server

import (
	"log"
	"os/exec"
	"time"
)

// CommandStatus 用于封装命令执行的状态
type CommandStatus struct {
	cmd       *exec.Cmd
	Running   bool
	Success   bool
	Timestamp time.Time
}

func NewPushRtmp(sysAddr, userAddr string) (*CommandStatus, error) {
	cmd := exec.Command("ffmpeg",
		"-f", "h264",
		"-i", "udp://127.0.0.1:65525",
		"-vcodec", "copy",
		"-an", // 这个参数用于禁用音频
		"-f", "flv",
		sysAddr,
		"-vcodec", "copy",
		"-an", // 这个参数用于禁用音频
		"-f", "flv",
		userAddr)
	log.Printf("直播推流[ 双 ]路转发启动sysAddr: %s; userAddr: %s \n", sysAddr, userAddr)
	// 定义ffmpeg命令
	return &CommandStatus{
		cmd:     cmd,
		Running: false,
	}, nil
}

func NewOnePushRtmp(addr string) (*CommandStatus, error) {
	cmd := exec.Command("ffmpeg",
		"-f", "h264",
		"-i", "udp://127.0.0.1:65525",
		"-vcodec", "copy",
		"-an", // 这个参数用于禁用音频
		"-f", "flv",
		addr)
	log.Printf("直播推流[ 单 ]路转发启动addr: %s;  \n", addr)
	// 定义ffmpeg命令
	return &CommandStatus{
		cmd:     cmd,
		Running: false,
	}, nil
}

func (c *CommandStatus) PushRtmp(done chan CommandStatus) error {
	// 输出命令详情（可选）
	log.Printf("Executing command: %s\n", c.cmd.String())
	// 启动命令
	err := c.cmd.Start()
	if err != nil {
		done <- CommandStatus{Running: true, Success: false, Timestamp: time.Now()}
		return err
	}
	// 等待命令完成
	err = c.cmd.Wait()
	if err != nil {
		done <- CommandStatus{Running: true, Success: false, Timestamp: time.Now()}
		return err
	}
	select {
	// 如果命令成功完成或者终止发出终止信息
	case done <- CommandStatus{Running: true, Success: true, Timestamp: time.Now()}:
	}
	return nil
}
