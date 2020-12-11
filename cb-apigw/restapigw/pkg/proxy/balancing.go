// Package proxy -
package proxy

import (
	"context"
	"net/url"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/sd"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// newLoadBalancedChain - Load Balancer 기능이 적용된 Call Chain 생성
func newLoadBalancedChain(sb sd.Balancer) CallChain {
	return func(next ...Proxy) Proxy {
		if 1 < len(next) {
			panic(ErrTooManyProxies)
		}

		return func(ctx context.Context, req *Request) (*Response, error) {
			host, err := sb.Host()
			if nil != err {
				return nil, err
			}

			r := req.Clone()

			var b strings.Builder
			b.WriteString(host)
			b.WriteString(r.Path)
			r.URL, err = url.Parse(b.String())
			if nil != err {
				return nil, err
			}
			if 0 < len(r.Query) {
				r.URL.RawQuery += "&" + r.Query.Encode()
			}

			return next[0](ctx, &r)
		}
	}
}

// ===== [ Public Functions ] =====

// NewLoadBalancedChainWithSubscriber - 지정된 Subscriber를 활용하는 Loadbalacer Chain 구성
func NewLoadBalancedChainWithSubscriber(subscriber sd.Subscriber) CallChain {
	return newLoadBalancedChain(sd.NewBalancer(subscriber))
}
