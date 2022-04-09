package zrpc

import "github.com/SunMaybo/zero/common/ip"

type Register struct {
	client CenterClient
}

func NewRegister(client CenterClient) *Register {
	return &Register{
		client: client,
	}
}

func (r *Register) DoRegister(instance ServiceInstance) {
	instance.Address = ip.LocalHostIP()
	if r.client != nil {
		if err := r.client.DoRegister(instance); err != nil {
			panic(err)
		}
	}
}
func (r *Register) UnRegister(instance ServiceInstance) {
	instance.Address = ip.LocalHostIP()
	if r.client != nil {
		if err := r.client.DeRegister(instance); err != nil {
			panic(err)
		}
	}
}
