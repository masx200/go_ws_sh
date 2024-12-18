/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http3

import (
	"context"
	"net/http"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/hertz/pkg/protocol/suite"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type Option struct{}

type Server struct {
	*http3.Server
	Option Option
	logger hlog.FullLogger
}

type handler struct {
	ctxPool *sync.Pool
	core    suite.Core
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := h.ctxPool.Get().(*app.RequestContext)
	_ = adaptor.CopyToHertzRequest(request, &ctx.Request)
	h.core.ServeHTTP(context.Background(), ctx)
	ctx.Response.Header.VisitAll(func(key, value []byte) {
		writer.Header().Add(string(key), string(value)) // TODO: B2S
	})
	writer.WriteHeader(ctx.Response.StatusCode())
	writer.Write(ctx.Response.Body())
	ctx.Reset()
	h.ctxPool.Put(ctx)
}

func (s *Server) Serve(c context.Context, conn network.StreamConn) error {
	cc, ok := conn.GetRawConnection().(quic.Connection)
	if !ok {
		return errors.NewPublicf("network-go http3: cannot convert raw connection to network.Connection")
	}
	return s.ServeQUICConn(cc)
}

func New(core suite.Core, logger hlog.FullLogger) *Server {
	handler := &handler{core: core, ctxPool: core.GetCtxPool()}
	s := &Server{Server: &http3.Server{}, logger: logger}
	s.Handler = handler
	return s
}
