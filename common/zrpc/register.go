package zrpc

import (
	"github.com/SunMaybo/zero/common/center"
	"github.com/SunMaybo/zero/common/ip"
)

type Register struct {
	client center.Center
}

func NewRegister(client center.Center) *Register {
	return &Register{
		client: client,
	}
}

func (r *Register) DoRegister(instance center.ServiceInstance) {
	instance.Address = ip.LocalHostIP()
	if r.client != nil {
		if err := r.client.DoRegister(instance); err != nil {
			panic(err)
		}
	}
}
func (r *Register) UnRegister(instance center.ServiceInstance) {
	instance.Address = ip.LocalHostIP()
	if r.client != nil {
		if err := r.client.DeRegister(instance); err != nil {
			panic(err)
		}
	}
}
