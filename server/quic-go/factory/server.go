/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package factory

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/suite"

	http3 "github.com/masx200/go_ws_sh/server/quic-go"
)

var _ suite.StreamServerFactory = &serverFactory{}

type serverFactory struct {
	option *http3.Option
}

// New is called by Hertz during engine.Run()
func (s *serverFactory) New(core suite.Core) (server protocol.StreamServer, err error) {
	serv := http3.New(core, hlog.SystemLogger())
	serv.Option = *s.option
	return serv, nil
}

func NewServerFactory(option *http3.Option) suite.StreamServerFactory {
	return &serverFactory{
		option: option,
	}
}
