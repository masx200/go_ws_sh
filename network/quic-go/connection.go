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

package quic

import (
	"context"
	"net"

	"github.com/cloudwego/hertz/pkg/network"
	quicgo "github.com/quic-go/quic-go"
)

type conn struct {
	rawConn interface{}
	quicgo.EarlyConnection
}

// Context implements network.StreamConn.
// Subtle: this method shadows the method (EarlyConnection).Context of conn.EarlyConnection.
func (c *conn) Context() context.Context {

	return c.EarlyConnection.Context()
}

// HandshakeComplete implements network.StreamConn.
// Subtle: this method shadows the method (EarlyConnection).HandshakeComplete of conn.EarlyConnection.
func (c *conn) HandshakeComplete() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-c.EarlyConnection.HandshakeComplete():
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx
}

// LocalAddr implements network.StreamConn.
// Subtle: this method shadows the method (EarlyConnection).LocalAddr of conn.EarlyConnection.
func (c *conn) LocalAddr() net.Addr {

	return c.EarlyConnection.LocalAddr()
}

// RemoteAddr implements network.StreamConn.
// Subtle: this method shadows the method (EarlyConnection).RemoteAddr of conn.EarlyConnection.
func (c *conn) RemoteAddr() net.Addr {
	return c.EarlyConnection.RemoteAddr()
}

type quicgoVersionNumber = uint32
type versioner interface {
	GetVersion() quicgoVersionNumber
}

func (c *conn) GetVersion() uint32 {
	cc, ok := c.EarlyConnection.(versioner)
	if !ok {
		return 0
	}
	return uint32(cc.GetVersion())
}

func (c *conn) GetRawConnection() interface{} {
	return c.rawConn
}

func (c *conn) AcceptStream(ctx context.Context) (network.Stream, error) {
	stream, err := c.EarlyConnection.AcceptStream(ctx)
	return newStream(stream), err
}

func (c *conn) AcceptUniStream(ctx context.Context) (network.ReceiveStream, error) {
	stream, err := c.EarlyConnection.AcceptUniStream(ctx)
	return newReadStream(stream), err
}

func (c *conn) OpenStream() (network.Stream, error) {
	stream, err := c.EarlyConnection.OpenStream()
	return newStream(stream), err
}

func (c *conn) OpenStreamSync(ctx context.Context) (network.Stream, error) {
	stream, err := c.EarlyConnection.OpenStreamSync(ctx)
	return newStream(stream), err
}

func (c *conn) OpenUniStream() (network.SendStream, error) {
	stream, err := c.EarlyConnection.OpenUniStream()
	return newWriteStream(stream), err
}

func (c *conn) OpenUniStreamSync(ctx context.Context) (network.SendStream, error) {
	stream, err := c.EarlyConnection.OpenUniStreamSync(ctx)
	return newWriteStream(stream), err
}

func (c *conn) CloseWithError(err network.ApplicationError, errMsg string) error {
	return c.EarlyConnection.CloseWithError(quicgo.ApplicationErrorCode(err.ErrCode()), errMsg)
}

func newStreamConn(qc quicgo.EarlyConnection) network.StreamConn {
	return &conn{qc, qc}
}
func init() {
	var _ network.StreamConn = &conn{}
}
