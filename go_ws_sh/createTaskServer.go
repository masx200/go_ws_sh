package go_ws_sh

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol/suite"
	"github.com/hertz-contrib/http2"
	"github.com/hertz-contrib/http2/config"
	factoryh2 "github.com/hertz-contrib/http2/factory"
	"github.com/hertz-contrib/logger/accesslog"

	quic "github.com/masx200/go_ws_sh/network/quic-go"
	http3 "github.com/masx200/go_ws_sh/server/quic-go"
	factoryh3 "github.com/masx200/go_ws_sh/server/quic-go/factory"
)

func createTaskServer(serverconfig ServerConfig, handler func(w context.Context, r *app.RequestContext)) func() (interface{}, error) {
	if serverconfig.Alpn == "h2" {

		return func() (interface{}, error) {
			cert, err := tls.LoadX509KeyPair(serverconfig.Cert, serverconfig.Key)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			cfg := &tls.Config{
				// add certificate
				Certificates: []tls.Certificate{cert},
				MaxVersion:   tls.VersionTLS13,
				// enable client authentication
				// ClientAuth: tls.RequireAndVerifyClientCert,
				// ClientCAs:  caCertPool,
				// cipher suites supported
				// CipherSuites: []uint16{
				// 	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				// 	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				// 	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				// },
				// set application protocol http2
				NextProtos: []string{http2.NextProtoTLS},
			}
			// cfg.NextProtos = append(cfg.NextProtos, "h2")
			hertzapp := server.Default( /* server.WithAltTransport(netpoll.NewTransporter) */ server.WithALPN(true), server.WithTLS(cfg), server.WithHostPorts(":"+serverconfig.Port), /*  server.WithTransport(quic.NewTransporter), */
				server.WithAltTransport(standard.NewTransporter))

			hertzapp.AddProtocol("h2", factoryh2.NewServerFactory(config.WithDisableKeepAlive(false)))
			// config.WithReadTimeout(time.Minute),
			// config.WithDisableKeepAlive(false)))
			hertzapp.Use(accesslog.New())
			// hertzapp.AddProtocol(suite.HTTP3, factoryh3.NewServerFactory(&http3.Option{}))
			log.Println("Alpn == h2")
			log.Println("TLS enabled and " + "WebSocket server started at :" + serverconfig.Port)
			// h2s := &http2.Server{
			// 	// ...
			// }
			// h1s := &http.Server{
			// 	Addr:    ":" + serverconfig.Port,
			// 	Handler: h2c.NewHandler(http.HandlerFunc(handler), h2s),
			// }
			hertzapp.GET("/:name", func(c context.Context, ctx *app.RequestContext) {
				handler(c, ctx)
			})
			x := hertzapp.Run()
			if x != nil {
				log.Fatal(x)
				return "", x
			}

			return "", nil
		}
	}
	// hertzapp.GET("/:name", func(c context.Context, ctx *app.RequestContext) {

	// })
	if serverconfig.Protocol == "https" {
		return func() (interface{}, error) {
			cert, err := tls.LoadX509KeyPair(serverconfig.Cert, serverconfig.Key)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			cfg := &tls.Config{
				// add certificate
				Certificates: []tls.Certificate{cert},
				MaxVersion:   tls.VersionTLS13,
				// enable client authentication
				// ClientAuth: tls.RequireAndVerifyClientCert,
				// ClientCAs:  caCertPool,
				// cipher suites supported
				// CipherSuites: []uint16{
				// 	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				// 	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				// 	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				// },
				// set application protocol http2
				NextProtos: []string{http2.NextProtoTLS},
			}
			// cfg.NextProtos = append(cfg.NextProtos, "h2")
			hertzapp := server.Default( /* server.WithAltTransport(netpoll.NewTransporter) */ server.WithALPN(true), server.WithTLS(cfg), server.WithHostPorts(":"+serverconfig.Port), server.WithTransport(quic.NewTransporter))
			/* server.WithAltTransport(standard.NewTransporter) */ //)

			// hertzapp.AddProtocol("h2", factoryh2.NewServerFactory(config.WithDisableKeepAlive(false)))
			// config.WithReadTimeout(time.Minute),
			// config.WithDisableKeepAlive(false)))
			hertzapp.Use(accesslog.New())
			hertzapp.AddProtocol(suite.HTTP3, factoryh3.NewServerFactory(&http3.Option{}))
			log.Println("Alpn == h3")
			log.Println("TLS enabled and " + "WebSocket server started at :" + serverconfig.Port)
			// h2s := &http2.Server{
			// 	// ...
			// }
			// h1s := &http.Server{
			// 	Addr:    ":" + serverconfig.Port,
			// 	Handler: h2c.NewHandler(http.HandlerFunc(handler), h2s),
			// }
			hertzapp.Any("/:name", func(c context.Context, ctx *app.RequestContext) {
				handler(c, ctx)
			})
			x := hertzapp.Run()
			if x != nil {
				log.Fatal(x)
				return "", x
			}

			return "", nil
		}

	} else {
		return func() (interface{}, error) {
			hertzapp := server.Default(server.WithHostPorts(":" + serverconfig.Port))
			hertzapp.Use(accesslog.New())
			log.Println("TLS disabled and " + "WebSocket server started at :" + serverconfig.Port)
			hertzapp.Any("/:name", func(c context.Context, ctx *app.RequestContext) {
				handler(c, ctx)
			})
			x := hertzapp.Run()
			if x != nil {
				log.Fatal(x)
				return "", x
			}

			// if err := http.ListenAndServe(":"+serverconfig.Port, http.HandlerFunc(handler)); err != nil {
			// 	log.Fatal("ListenAndServe: ", err)
			// 	return "", err
			// }
			return "", nil
		}
	}

}
