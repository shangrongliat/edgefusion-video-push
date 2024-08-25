package server

import (
	"sync"

	"edgefusion-video-push/config"
	"edgefusion-video-push/service"
)

// 按照传入2个对象进行推流控制，在前置new出要推送的对象，传输进来直接进行转发，不在转发里进行new操作
func Consume(group *sync.WaitGroup, queue *service.Queue, tr TransmitInfo, cfg *config.Config) {
	done := make(chan CommandStatus)
	for {
		select {
		case <-queue.ItemChan:
			data, ok := queue.Pull()
			if !ok && data != nil {
				if video, ok := data.([]byte); ok {
					//如果取数据成功
					//根据配置启动分发策略

				}
			}
		}
	}
}
