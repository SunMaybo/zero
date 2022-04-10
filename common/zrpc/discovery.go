package zrpc

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
	"strings"
	"time"
)

const (
	defaultFreq = time.Minute * 30
)

type ResolverBuilder struct {
	center CenterClient
}

func NewResolverBuilder(center CenterClient) *ResolverBuilder {
	return &ResolverBuilder{center: center}
}

func (e *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		center: e.center,
		target: target,
		cc:     cc,
		store:  map[string]struct{}{},
		stopCh: make(chan struct{}, 1),
		rn:     make(chan struct{}, 1),
		t:      time.NewTicker(defaultFreq),
	}
	// 调用 start 初始化地址
	go r.start()
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (e *ResolverBuilder) Scheme() string { return e.center.GetSchema() }

type Resolver struct {
	center CenterClient
	target resolver.Target
	cc     resolver.ClientConn
	store  map[string]struct{}
	stopCh chan struct{}
	// rn channel is used by ResolveNow() to force an immediate resolution of the target.
	rn chan struct{}
	t  *time.Ticker
}

func (r *Resolver) start() {
	//nacos://com.ikurento.user.UserService?cluster=DEFAULT&group=DEFAULT&tenant=DEFAULT
	if r.target.URL.Scheme != r.center.GetSchema() {
		r.cc.ReportError(errors.Errorf("scheme must be %s", r.center.GetSchema()))
		return
	}
	clusters := r.target.URL.Query().Get("cluster")
	if clusters == "" {
		clusters = "DEFAULT"
	}
	groups := r.target.URL.Query().Get("group")
	if groups == "" {
		groups = "DEFAULT_GROUP"
	}
	rch := make(chan []ServiceInstance, 1)
	err := r.center.Subscribe(&SubscribeParam{
		ServiceName: r.target.URL.Hostname(),
		Clusters:    strings.Split(strings.TrimSpace(clusters), ","),
		GroupName:   groups,
		SubscribeCallback: func(services []ServiceInstance) {
			rch <- services
		},
	})
	if err != nil {
		r.cc.ReportError(errors.Wrap(err, "failed to subscribe service"))
	}
	for {
		select {
		case <-r.rn:
			r.resolveNow()
		case <-r.t.C:
			r.ResolveNow(resolver.ResolveNowOptions{})
		case <-r.stopCh:
			close(r.rn)
			close(r.stopCh)
			return
		case wresp := <-rch:
			r.store = make(map[string]struct{})
			for _, service := range wresp {
				r.store[service.Address+":"+fmt.Sprintf("%d", service.Port)] = struct{}{}
			}
			r.updateTargetState()
		}
	}
}
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}
func (r *Resolver) Close() {
	r.t.Stop()
	r.stopCh <- struct{}{}
}
func (r *Resolver) resolveNow() {
	clusters := r.target.URL.Query().Get("cluster")
	if clusters == "" {
		clusters = "DEFAULT"
	}
	groups := r.target.URL.Query().Get("group")
	if groups == "" {
		groups = "DEFAULT_GROUP"
	}
	instances, err := r.center.SelectInstances(SelectInstancesParam{
		ServiceName: r.target.URL.Hostname(),
		Clusters:    strings.Split(strings.TrimSpace(clusters), ","),
		GroupName:   groups,
		HealthyOnly: true,
	})
	if err != nil {
		r.cc.ReportError(errors.Wrap(err, "failed to select instances"))
		return
	}
	r.store = make(map[string]struct{})
	for _, instance := range instances {
		r.store[instance.Address+":"+fmt.Sprintf("%d", instance.Port)] = struct{}{}
	}
	r.updateTargetState()
}

func (r *Resolver) updateTargetState() {
	addrs := make([]resolver.Address, len(r.store))
	i := 0
	for k := range r.store {
		addrs[i] = resolver.Address{Addr: k}
		i++
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
