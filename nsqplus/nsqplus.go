package nsqplus

import (
	"sync"
	"net"
	"internal/util"
	"os"
	"fmt"
	"internal/version"
	"crypto/tls"
	"internal/http_api"
	"strings"
)

type NSQPlus struct {
	sync.RWMutex
	opts                *Options
	httpListener        net.Listener
	waitGroup           util.WaitGroupWrapper
	httpClientTLSConfig *tls.Config
}

func New(opts *Options) *NSQPlus {
	n := &NSQPlus{
		opts:          opts,
	}
	n.logf("nsqdHttpAddresses :%s", strings.Join(opts.NSQDHTTPAddresses, ","))
	if len(opts.NSQDHTTPAddresses) == 0 {
		n.logf("--nsqd-http-address  required.")
		os.Exit(1)
	}

	// verify that the supplied address is valid
	verifyAddress := func(arg string, address string) *net.TCPAddr {
		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			n.logf("FATAL: failed to resolve %s address (%s) - %s", arg, address, err)
			os.Exit(1)
		}
		return addr
	}
	for _, address := range opts.NSQDHTTPAddresses {
		verifyAddress("--nsqd-http-address", address)
	}

	n.logf(version.String("nsq-plus producer"))
	return n
}

func (n *NSQPlus) logf(f string, args ...interface{}) {
	if n.opts.Logger == nil {
		return
	}
	n.opts.Logger.Output(2, fmt.Sprintf(f, args...))
}


func (n *NSQPlus)Main() {
	httpListener, err := net.Listen("tcp", n.opts.HTTPAddress)
	if err != nil {
		n.logf("FATAL: listen (%s) failed - %s", n.opts.HTTPAddress, err)
		os.Exit(1)
	}
	n.Lock()
	n.httpListener = httpListener
	n.Unlock()

	httpServer := NewHTTPServer(&Context{n})
	n.waitGroup.Wrap(func() {
		http_api.Serve(n.httpListener, http_api.CompressHandler(httpServer), "HTTP", n.opts.Logger)
	})
}


func (n *NSQPlus) Exit() {
	n.httpListener.Close()
	n.waitGroup.Wait()
}
