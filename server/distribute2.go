package server

import (
	"edgefusion-video-push/config"
	"edgefusion-video-push/service"
	"github.com/robfig/cron"
	"log"
)

var distPush *CommandStatus

// Consume2 按照传入2个对象进行推流控制，在前置new出要推送的对象，传输进来直接进行转发，不在转发里进行new操作
// 队列只为处理数据是否活跃，来决定是否要拉起ffmpeg
func Consume2(push *CommandStatus, queue *service.Queue, cfg config.Config) {
	distPush = push
	done := make(chan CommandStatus)
	if push != nil {
		// 默认进方法先执行一次
		go func() {
			if err := distPush.PushRtmp(done); err != nil {
				log.Println("推流命令启动执行失败.", err)
			}
		}()
		c := cron.New()
		if err := c.AddFunc("@every 22s", func() {
			// running为ture说明command执行结束，需要重新开始
			log.Printf("push 状态: %v,队列状态: %v", distPush.Running, queue.Status())
			if distPush.Running && queue.Status() == 0 {
				distPush = RetryPush(cfg)
				if err := distPush.PushRtmp(done); err != nil {
					log.Println("推流命令启动执行失败.", err)
				}
			}
		}); err != nil {
			log.Println("定时队列状态监测启动失败....", err)
		}
		c.Start()
		go func() {
			for {
				select {
				case tr, ok := <-done:
					if ok {
						distPush = &tr
					}
				}
			}
		}()
	}
}
