package accesscontrol

import (
	"net"
	"net/http"
	"runtime/trace"

	"abstraction.fr/config"
	log "github.com/sirupsen/logrus"
)

// IPWhitelistMiddleware ...
type IPWhitelistMiddleware struct {
	Config *config.Config
	Logger *log.Entry
	next   http.Handler
}

// NewIPWhitelistMiddleware returns a *WhitelistMiddleware. next is an http.Handler
// which will be used if the request IP is not whitelisted. If next is nil then
// the middleware will return http.StatusUnauthorized (401)
func NewIPWhitelistMiddleware(conf *config.Config, logger *log.Entry, next http.Handler) *IPWhitelistMiddleware {
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

		// local logger
		llogger := mdw.Logger.WithFields(log.Fields{
			"_func": "WhitelistMiddleware.Middleware",
			"_path": r.URL.Path,
		})

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
				llogger.Errorf("%s is not a valid CIDR", sub)
				continue
			}

			if cidr.Contains(netIP) {
				llogger.Tracef("%v is whitelisted by %v", netIP, cidr)
				next.ServeHTTP(w, r)
				return
			}
		}

		llogger.Debugf("%v is not whitelisted", netIP)

		if mdw.next != nil {
			llogger.Debugf("next handler is not nil, using it")
			mdw.next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
