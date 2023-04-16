package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/patrickmn/go-cache"
)

type DNSResolver struct {
	cache *cache.Cache
}

func NewDNSResolver() *DNSResolver {
	return &DNSResolver{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (d *DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	if ip, ok := d.cache.Get(name); ok {
		return ctx, ip.(net.IP), nil
	}

	dnsServers := []string{"8.8.8.8:53", "8.8.4.4:53"} // 自定义DNS服务器地址列表

	timeoutSecond := 5 * time.Second
	ch := make(chan net.IP)
	for _, dnsServer := range dnsServers {
		go func(dnsServer string) {
			resolver := &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					dialer := &net.Dialer{
						Timeout: timeoutSecond,
					}
					conn, err := dialer.DialContext(ctx, "udp", dnsServer)
					if err != nil {
						return nil, err
					}
					return conn, nil
				},
			}
			ips, err := resolver.LookupIPAddr(ctx, name)
			if err != nil {
				fmt.Printf("lookup %s on %s failed: %v\n", name, dnsServer, err)
				return
			}
			select {
			case ch <- ips[0].IP:
				fmt.Printf("lookup %s on %s success: %v\n", name, dnsServer, ips[0].IP)
			default:
				fmt.Printf("lookup %s on %s failed 2: %v\n", name, dnsServer, ips[0].IP)
			}
		}(dnsServer)
	}

	select {
	case ip := <-ch:
		d.cache.Set(name, ip, cache.DefaultExpiration)
		return ctx, ip, nil
	case <-time.After(timeoutSecond): // 超时时间可根据需要进行修改
		fmt.Printf("lookup timeout\n")
	}

	return ctx, nil, fmt.Errorf("lookup %s failed", name)
}
