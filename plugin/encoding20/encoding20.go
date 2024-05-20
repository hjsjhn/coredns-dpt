package encoding20

import (
	"net"
	"unicode"
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// speific id for 0x20 encoding response
type Sid struct {
	Next plugin.Handler
	Data string
}

func (h Sid) Name() string { return "20encoding" }

func (m Sid) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// 检查DNS请求中的域名是否为大小写混杂的请求
	if !isLower(state.Name()) {
		// 如果是大小写混杂的域名请求，返回特殊的结果IP
		responseIP := m.Data
		return m.returnSpecialIP(w, r, responseIP)
	}

	// 如果不是大小写混杂的域名请求，继续传递请求给下一个插件
	return plugin.NextOrFailure(m.Name(), m.Next, ctx, w, r)
}

func isLower(domain string) bool {
	for _,r := range domain {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// 返回特殊IP的函数
func (m Sid) returnSpecialIP(w dns.ResponseWriter, r *dns.Msg, ip string) (int, error) {
	// 创建特殊IP的DNS响应
	response := new(dns.Msg)
	ip_net := net.ParseIP(ip)
	if ip_net == nil {
		// create a server failure err
		err := dns.ErrConnEmpty
		return dns.RcodeServerFailure, err
	}
	// set the IP field in response to ip
	response.Answer = append(response.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   ip_net,
	})

	// 发送DNS响应
	err := w.WriteMsg(response)
	if err != nil {
		return dns.RcodeServerFailure, err
	}

	return dns.RcodeSuccess, nil
}