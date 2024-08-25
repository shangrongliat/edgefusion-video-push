package server

import (
	"net"
	"time"

	"edgefusion-video-push/service"
)

// 按照传入2个对象进行推流控制，在前置new出要推送的对象，传输进来直接进行转发，不在转发里进行new操作
func Consume(listen *Listener, queue *service.Queue, transmit *TransmitInfo, push *CommandStatus) {
	done := make(chan CommandStatus)
	var addrs []*net.UDPAddr
	if transmit != nil {
		addrs = append(addrs, transmit.RemoteAddr)
	}
	if push != nil {
		go push.PushRtmp(done)
		addrs = append(addrs, push.transmitInfo.RemoteAddr)
		go func() {
			select {
			case tr := <-done:
				// running为ture说明command执行结束，需要重新开始
				if tr.Running && queue.Status() == 0 {
					push.PushRtmp(done)
					// 推流停止后没隔60秒进行一次重试，基于UDP是否能接收到数据
					time.Sleep(60 * time.Second)
				}
			}
		}()
	}
	pushExc(listen, queue, addrs)
}

func pushExc(listen *Listener, queue *service.Queue, addr []*net.UDPAddr) {
	for {
		select {
		case <-queue.DataChan:
			data, ok := queue.Pull()
			if ok && data != nil {
				if video, ok := data.([]byte); ok {
					//如果取数据成功
					//根据配置启动分发策略
					listen.Transmit(video, addr...)
				}
			}
		}
	}
}
