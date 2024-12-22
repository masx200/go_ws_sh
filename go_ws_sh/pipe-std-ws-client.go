package go_ws_sh

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/linkedin/goavro/v2"
	// "google.golang.org/grpc/authz/audit/stdout"
	// "golang.org/x/term"
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

func Client_start(config string) {
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
	// var httpServeMux = http.NewServeMux()

	pipe_std_ws_client(configdata /* httpServeMux, handler */)
}

// pipe_std_ws_client 创建一个WebSocket客户端，根据配置数据连接到服务器。
// 它处理与WebSocket服务器的通信，包括身份验证、消息编码和解码，以及与标准输入/输出的交互。
func pipe_std_ws_client(configdata ConfigClient) {
	var binaryandtextchannel = make(chan WebsocketMessage)
	defer close(binaryandtextchannel)
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	codec, err := create_msg_codec()
	if err != nil {
		log.Println(err)
		// r.SetStatusCode(http.StatusInternalServerError)

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
	url := x1 + "://" + configdata.Servers.Host + ":" + configdata.Servers.Port + "/" + configdata.Sessions.Path // 替换为你的WebSocket服务器地址

	// 创建一个新的WebSocket客户端连接
	x := websocket.DefaultDialer
	if configdata.Servers.Ca != "" || configdata.Servers.Insecure {
		configureWebSocketTLSCA(x, configdata)
	}

	x.EnableCompression = true
	header := http.Header{}
	username := configdata.Credentials.Username
	password := configdata.Credentials.Password
	// fmt.Println("username:", username, "password:", password)
	authStr := username + ":" + password
	authBytes := []byte(authStr)
	encodedAuth := base64.StdEncoding.EncodeToString(authBytes)

	// 将编码后的字符串设置到Authorization字段
	authHeader := "Basic " + encodedAuth
	// fmt.Println("Authorization:", authHeader)
	header.Set("Authorization", authHeader)
	conn, response, err := x.Dial(url, header)

	defer func() {
		defer conn.WriteMessage(websocket.CloseMessage, []byte{})
		// if err := defer conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		// 	log.Println(err)
		// }
	}()
	defer func() { os.Exit(0) }()
	if response != nil {
		log.Println("Response Status:", response.Status)
		fmt.Println("Response Headers:")
		fmt.Println("{")
		for k, v := range response.Header {
			fmt.Println(k, ":", strings.Join(v, ","))
		}
		fmt.Println("}")
	}

	if err != nil {
		log.Println("Dial error:", err)

		return
	}

	defer conn.Close()
	go func() {
		var err error
		for {
			//var encoded,ok <-  binaryandtextchannel
			encoded, ok := <-binaryandtextchannel
			if ok {
				// mubinary.Lock()
				// defer mubinary.Unlock()
				err = conn.WriteMessage(encoded.Type, encoded.Body)
				if err != nil {
					log.Println("write:", err)
					return
				}
			} else {
				break
			}
		}
	}()
	// var in_queue = make(chan []byte)
	// var err_queue = make(chan []byte)
	// var out_queue = make(chan []byte)
	// defer close(out_queue)
	// defer close(err_queue)
	// defer close(in_queue)
	// go func() {
	// 	CopyChanToWriter(os.Stdout, out_queue)
	// }()
	// go func() {
	// 	CopyChanToWriter(os.Stderr, err_queue)
	// }()
	// 使用termbox接管stdin了
	// go func() {
	// 	io.Copy(in_queue, os.Stdin)
	// }()
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
		binaryandtextchannel <- WebsocketMessage{
			Body: databuf,
			Type: websocket.TextMessage,
		}
	}
	closable, startable, cols, rows, err := TermboxPipe(func(p []byte) (n int, err error) {

		// log.Println("write to stdin length:", len(p))
		// go func() {

		// in_queue <- (p)
		sendMessageToWebsocketStdin(p /* conn, */, codec, binaryandtextchannel)

		// }()

		return n, nil
	}, func() error {

		/* 	err :=  */
		defer conn.WriteMessage(websocket.CloseMessage, []byte{})
		// if err != nil {
		// 	log.Println(err)
		// }
		conn.Close()
		// close(out_queue)
		// close(err_queue)
		// close(in_queue)
		defer os.Exit(0)
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
	binaryandtextchannel <- WebsocketMessage{
		Body: databuf,
		Type: websocket.TextMessage,
	}
	// go func() {
	// 	for {
	// 		var data, err = ReadFixedSizeFromReader(os.Stdin, 1024*1024)
	// 		if data == nil || nil != err {
	// 			if err != nil {
	// 				log.Println("encode:", err)
	// 				return
	// 			}
	// 		}
	// 		log.Printf("os.stdin recv Binary length: %v", len(data))
	// 	}

	// }()
	defer func() { go closable() }()
	go func() {
		startable()
	}()
	// go func() {

	// 	for {
	// 		data, ok := <-in_queue
	// 		if data == nil || !ok {
	// 			break
	// 		}
	// 		log.Println("stdin recv Binary length: ", len(data))
	// 		var message = BinaryMessage{
	// 			Type: "stdin",
	// 			Body: data,
	// 		}

	// 		encoded, err := EncodeStructAvroBinary(codec, &message)
	// 		if err != nil {
	// 			log.Println("encode:", err)
	// 			continue
	// 		}

	// 		err = conn.WriteMessage(websocket.BinaryMessage, encoded)

	// 		if err != nil {
	// 			log.Println("write:", err)

	// 		}
	// 	}
	// }()
	for {
		mt, message, err := conn.ReadMessage()
		if err, ok := err.(*websocket.CloseError); ok {

			log.Println("close:", err)

			defer os.Exit(0)
			break
		}
		if err != nil {
			log.Println("read:", err)
			return
		}
		if mt == websocket.TextMessage {
			var array []string
			//parse json data

			err = json.Unmarshal(message, &array)
			if err != nil {
				log.Println("read:", err)
				//return
				// log.Printf("ignored recv text: %s", message)
				return
			}
			var data TextMessage
			err = DecodeTextMessageFromStringArray(array, &data)
			if err != nil {
				log.Println("read:", err)
				//return
				// log.Printf("ignored recv text: %s", message)
				return
			}
			// log.Println("websocket recv text length: ", len(message))
			if data.Type == "rejected" {
				log.Println("rejected:", data.Body)
				defer os.Exit(0)
				return
				//break
			} else if data.Type == "resolved" {
				log.Println("resolved:", data.Body)
			} else {
				log.Printf("ignored unknown recv text:%v", data)
				return
			}
		} else {
			// log.Println("websocket recv binary length: ", len(message))
			var result BinaryMessage
			undecoded, err := DecodeStructAvroBinary(codec, message, &result)
			if len(undecoded) > 0 {

				log.Println("undecoded:", undecoded)

			}
			if err != nil {
				log.Println("decode:", err)

			} else {
				// log.Printf("recv binary: %s", decoded)
				var md = result
				/* 	if md.Type == "stderr" {
					var body = md.Body
					// go func() {

					// err_queue <- body
					os.Stderr.Write(body)
					// }()

				} else */if md.Type == "stdout" {
					var body = md.Body
					// go func() {
					os.Stdout.Write(body)
					// out_queue <- body
					// }()
				} else {
					log.Println("ignored unknown type:", md.Type)
					return
				}
			}
		}
	}
	// cmd := exec.Command("pwsh")

	// // 设置标准输入、输出和错误流
	// stdin, err := cmd.StdinPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // 启动命令
	// if err := cmd.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	// // 处理标准输出和错误流
	// go func() {
	// 	io.Copy(os.Stdout, stdout)
	// }()
	// go func() {
	// 	io.Copy(os.Stderr, stderr)
	// }()

	// // 处理标准输入流
	// io.Copy(stdin, os.Stdin)

	// // 等待命令执行完毕
	// if err := cmd.Wait(); err != nil {
	// 	log.Fatal(err)
	// }
}

func sendMessageToWebsocketStdin(data []byte /* conn *websocket.Conn, */, codec *goavro.Codec, binaryandtextchannel chan WebsocketMessage) error {
	// log.Println("stdin recv Binary length: ", len(data))
	var message = BinaryMessage{
		Type: "stdin",
		Body: data,
	}

	encoded, err := EncodeStructAvroBinary(codec, &message)
	if err != nil {
		log.Println("encode:", err)
		return err
	}
	binaryandtextchannel <- WebsocketMessage{
		Body: encoded,
		Type: websocket.BinaryMessage,
	}
	return nil
	// err = conn.WriteMessage(websocket.BinaryMessage, encoded)

	// if err != nil {
	// 	log.Println("write:", err)

	// }
}

// configureWebSocketTLSCA 配置 WebSocket 的 TLS 客户端证书认证。
// 该函数接收一个 websocket.Dialer 实例 x 和一个 ConfigClient 类型的配置数据 configdata。
// 它的主要作用是为 WebSocket 连接配置 TLS 客户端，以确保连接的安全性。
func configureWebSocketTLSCA(x *websocket.Dialer, configdata ConfigClient) {
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
}
