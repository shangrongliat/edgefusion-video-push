package server

import (
	"fmt"
	"log"

	"github.com/robfig/cron"
)

// 域名定时解析
type AnalysisDomain struct {
	f      *Forward
	domain string
	cron   *cron.Cron // 定时器
}

func NewAnalysis(f *Forward, domain string) *AnalysisDomain {
	return &AnalysisDomain{
		f:      f,
		domain: domain,
		cron:   cron.New(),
	}
}

func (a *AnalysisDomain) Start() {
	job := func() {
		resolver := DomainResolver(a.domain)
		if resolver != a.f.TransmitAddrIp {
			transmit, err := NewTransmit(fmt.Sprintf("%s:65506", resolver))
			if err != nil {
				log.Printf("视频[ 透传转发 ] 启动失败,转发地址: %s:65506 \n", resolver)
			}
			//重新覆盖本地转发地址
			a.f.SetTransmitAddr(transmit, 1)
		}
	}
	if err := a.cron.AddFunc("@every 10s", job); err != nil {
		return
	}
	a.cron.Start()
}
