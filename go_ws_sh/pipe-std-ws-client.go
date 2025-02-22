package go_ws_sh

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/linkedin/goavro/v2"
)

type ClientSession struct {
	Path string `json:"path"`
}

type ConfigClient struct {
	Credentials Credentials   `json:"credentials"`
	Sessions    ClientSession `json:"sessions"`
	Servers     ClientConfig  `json:"servers"`
}
type ClientConfig struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Ca       string `json:"ca"`
	Insecure bool   `json:"insecure"`
	Host     string `json:"host"`
}

func Client_start(config string, serverip string) {
	configFile, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	jsonDecoder := json.NewDecoder(configFile)
	var configdata ConfigClient
	err = jsonDecoder.Decode(&configdata)
	if err != nil {
		log.Fatal(err)
	}

	pipe_std_ws_client(configdata, serverip)
}

func pipe_std_ws_client(configdata ConfigClient, serverip string) {
	var binaryandtextchannel = NewSafeChannel[WebsocketMessage]()
	defer (binaryandtextchannel).Close()

	codec, err := create_msg_codec()
	if err != nil {
		log.Println(err)

		return
	}
	protocol_map := map[string]string{
		"http":  "ws",
		"https": "wss",
	}
	x1, ok := protocol_map[configdata.Servers.Protocol]
	if !ok {
		log.Fatal("unknown protocol:", configdata.Servers.Protocol)
	}
	url := x1 + "://" + configdata.Servers.Host + ":" + configdata.Servers.Port + "/" + configdata.Sessions.Path

	x := websocket.DefaultDialer
	var tlscfg = &tls.Config{}
	if configdata.Servers.Ca != "" || configdata.Servers.Insecure {
		tlscfg = configureWebSocketTLSCA(x, configdata)
	}
	if serverip != "" {

		var serverIP = serverip
		x = &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// 解析出原地址中的端口
				_, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				// 用指定的 IP 地址和原端口创建新地址
				newAddr := net.JoinHostPort(serverIP, port)
				// 创建 net.Dialer 实例
				dialer := &net.Dialer{}
				// 发起连接
				return dialer.DialContext(ctx, network, newAddr)
			},
			NetDialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {

				// 解析出原地址中的端口
				address, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				// 用指定的 IP 地址和原端口创建新地址
				newAddr := net.JoinHostPort(serverIP, port)
				// 创建 net.Dialer 实例
				dialer := &net.Dialer{}
				// 发起连接
				conn, err := dialer.DialContext(ctx, network, newAddr)
				if err != nil {
					return nil, err
				}
				tlsConfig :=
					tlscfg
				tlsConfig.ServerName = address
				// &tls.Config{
				// 	ServerName: address,
				// }
				// 创建 TLS 连接
				tlsConn := tls.Client(conn, tlsConfig)
				// 进行 TLS 握手
				err = tlsConn.HandshakeContext(ctx)
				if err != nil {
					conn.Close()
					return nil, err
				}
				return tlsConn, nil
			},
		}
	}
	if configdata.Servers.Ca != "" || configdata.Servers.Insecure {
		configureWebSocketTLSCA(x, configdata)
	}

	x.EnableCompression = true
	header := http.Header{}
	username := configdata.Credentials.Username
	password := configdata.Credentials.Password

	authStr := username + ":" + password
	authBytes := []byte(authStr)
	encodedAuth := base64.StdEncoding.EncodeToString(authBytes)

	authHeader := "Basic " + encodedAuth

	header.Set("Authorization", authHeader)
	conn, response, err := x.Dial(url, header)
	if err != nil {
		log.Println("Dial error:", err)

		return
	}
	defer func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in panic", r)
			}
		}()
		defer conn.WriteMessage(websocket.CloseMessage, []byte{})

	}()
	// defer func() { os.Exit(0) }()
	if response != nil {
		log.Println("Response Status:", response.Status)
		fmt.Println("Response Headers:")
		fmt.Println("{")
		for k, v := range response.Header {
			fmt.Println(k, ":", strings.Join(v, ","))
		}
		fmt.Println("}")
	}

	defer conn.Close()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in panic", r)
			}
		}()
		SendMessageToWebSocketLoop(conn, binaryandtextchannel)

	}()

	var onsizechange = func(cols int, rows int) {

		var msgsize = EncodeMessageSizeToStringArray(MessageSize{
			Cols: int64(cols),
			Rows: int64(rows),
			Type: "resize",
		})
		databuf, err := json.Marshal(msgsize)
		if err != nil {
			log.Println(err)
			return
		}
		binaryandtextchannel.Send(WebsocketMessage{
			Body: databuf,
			Type: websocket.TextMessage,
		})
	}
	closable, startable, cols, rows, err := TermboxPipe(func(p []byte) (n int, err error) {

		sendMessageToWebsocketStdin(p, codec, binaryandtextchannel)

		return n, nil
	}, func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in panic", r)
			}
		}()
		defer conn.WriteMessage(websocket.CloseMessage, []byte{})

		conn.Close()

		// defer os.Exit(0)
		return nil
	}, onsizechange)
	if err != nil {
		log.Println(err)
		return
	}
	var msgsize = EncodeMessageSizeToStringArray(MessageSize{
		Cols: int64(cols),
		Rows: int64(rows),
		Type: "resize",
	})
	databuf, err := json.Marshal(msgsize)
	if err != nil {
		log.Println(err)
		return
	}
	binaryandtextchannel.Send(WebsocketMessage{
		Body: databuf,
		Type: websocket.TextMessage,
	})

	defer func() { go closable() }()
	go func() {
		startable()
	}()

	for {
		mt, message, err := ReadMessageFromWebSocket(conn)
		if err, ok := err.(*websocket.CloseError); ok {

			log.Println("close:", err)

			// defer os.Exit(0)
			break
		}
		if err != nil {
			log.Println("read7:", err)
			return
		}
		if mt == websocket.TextMessage {
			var array []string

			err = json.Unmarshal(message, &array)
			if err != nil {
				log.Println("read8:", err)

				return
			}
			var data TextMessage
			err = DecodeTextMessageFromStringArray(array, &data)
			if err != nil {
				log.Println("read9:", err)

				return
			}

			if data.Type == "rejected" {
				log.Println("rejected:", data.Body)
				// defer os.Exit(0)
				return

			} else if data.Type == "resolved" {
				log.Println("resolved:", data.Body)
			} else {
				log.Printf("ignored unknown recv text:%v", data)
				return
			}
		} else {

			var result BinaryMessage
			undecoded, err := DecodeStructAvroBinary(codec, message, &result)
			if len(undecoded) > 0 {

				log.Println("undecoded:", undecoded)

			}
			if err != nil {
				log.Println("decode:", err)

			} else {

				var md = result
				if md.Type == "stdout" {
					var body = md.Body

					os.Stdout.Write(body)

				} else {
					log.Println("ignored unknown type:", md.Type)
					return
				}
			}
		}
	}

}

func sendMessageToWebsocketStdin(data []byte, codec *goavro.Codec, binaryandtextchannel *SafeChannel[WebsocketMessage]) error {

	var message = BinaryMessage{
		Type: "stdin",
		Body: data,
	}

	encoded, err := EncodeStructAvroBinary(codec, &message)
	if err != nil {
		log.Println("encode:", err)
		return err
	}
	binaryandtextchannel.Send(WebsocketMessage{
		Body: encoded,
		Type: websocket.BinaryMessage,
	})
	return nil

}

func configureWebSocketTLSCA(x *websocket.Dialer, configdata ConfigClient) *tls.Config {
	x.TLSClientConfig = &tls.Config{
		RootCAs: x509.NewCertPool(),
	}
	x.TLSClientConfig.InsecureSkipVerify = configdata.Servers.Insecure
	caCert, err := os.ReadFile(configdata.Servers.Ca)
	if err != nil {
		log.Fatal(err)
	}
	ok := x.TLSClientConfig.RootCAs.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("Failed to append CA certificate")
	}
	return x.TLSClientConfig
}
