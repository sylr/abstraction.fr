package accesscontrol

import (
	"fmt"
	"net"
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"

	"go.uber.org/zap"
)

// IPWhitelistMiddleware ...
type IPWhitelistMiddleware struct {
	Config *config.Config
	Logger *zap.Logger
	next   http.Handler
}

// NewIPWhitelistMiddleware returns a *WhitelistMiddleware. next is an http.Handler
// which will be used if the request IP is not whitelisted. If next is nil then
// the middleware will return http.StatusUnauthorized (401)
func NewIPWhitelistMiddleware(conf *config.Config, logger *zap.Logger, next http.Handler) *IPWhitelistMiddleware {
	mdw := IPWhitelistMiddleware{
		Config: conf,
		Logger: logger,
		next:   next,
	}

	return &mdw
}

// Middleware ...
func (mdw *IPWhitelistMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := trace.StartRegion(r.Context(), "IPWhitelistMiddleware")
		defer tr.End()

		// If no configuration set for trusted subnets we continue to the next middleware
		if mdw.Config.Internal == nil || mdw.Config.Internal.TrustedSubnets == nil {
			next.ServeHTTP(w, r)
			return
		}

		// Logger fields
		zfields := zap.Fields(
			zap.String("_func", "WhitelistMiddleware.Middleware"),
			zap.String("_path", "r.URL.Path"),
		)

		// Get client IP from HTTP headers set by the reverse proxy
		ip := r.Header.Get("X-Real-Ip")

		// If no client IP found we default to the one from the tcp connection
		if len(ip) == 0 {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}

		netIP := net.ParseIP(ip)
		for _, sub := range mdw.Config.Internal.TrustedSubnets {
			_, cidr, err := net.ParseCIDR(sub)

			if err != nil {
				mdw.Logger.WithOptions(zfields).Error(fmt.Sprintf("%s is not a valid CIDR", sub), zap.Error(err))
				continue
			}

			if cidr.Contains(netIP) {
				mdw.Logger.WithOptions(zfields).Debug(fmt.Sprintf("%v is whitelisted by %v", netIP, cidr))
				next.ServeHTTP(w, r)
				return
			}
		}

		mdw.Logger.Debug(fmt.Sprintf("%v is not whitelisted", netIP))

		if mdw.next != nil {
			mdw.Logger.Debug(fmt.Sprintf("next handler is not nil, using it"))
			mdw.next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
