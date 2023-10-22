package self

import (
	"net/http"
)

func (p *OAuthProxy) SetUpstreamProxy(proxy http.Handler) {
	p.upstreamProxy = proxy
}
