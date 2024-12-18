package go_ws_sh

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
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

	Host string `json:"host"`
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
	defer os.Exit(0)
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
	if configdata.Servers.Ca != "" {
		configureWebSocketTLSCA(x, configdata)
	}

	x.EnableCompression = true
	header := http.Header{}
	username := configdata.Credentials.Username
	password := configdata.Credentials.Password
	fmt.Println("username:", username, "password:", password)
	authStr := username + ":" + password
	authBytes := []byte(authStr)
	encodedAuth := base64.StdEncoding.EncodeToString(authBytes)

	// 将编码后的字符串设置到Authorization字段
	authHeader := "Basic " + encodedAuth
	// fmt.Println("Authorization:", authHeader)
	header.Set("Authorization", authHeader)
	conn, response, err := x.Dial(url, header)
	log.Println("Response Status:", response.Status)
	fmt.Println("Response Headers:")
	for k, v := range response.Header {
		fmt.Println(k, ":", strings.Join(v, ","))
	}
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	defer conn.Close()
	var in_queue = NewQueue()
	var err_queue = NewQueue()
	var out_queue = NewQueue()
	defer out_queue.Close()
	defer err_queue.Close()
	defer in_queue.Close()
	go func() {
		io.Copy(os.Stdout, out_queue)
	}()
	go func() {
		io.Copy(os.Stderr, err_queue)
	}()
	go func() {
		io.Copy(in_queue, os.Stdin)
	}()

	go func() {

		for {
			data := in_queue.Dequeue()
			if data == nil {
				break
			}

			encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
				"type": "stdin",
				"body": data,
			})
			if err != nil {
				log.Println("encode:", err)
				continue
			}

			err = conn.WriteMessage(websocket.BinaryMessage, encoded)

			if err != nil {
				log.Println("write:", err)

			}
		}
	}()
	for {
		mt, message, err := conn.ReadMessage()
		if err, ok := err.(*websocket.CloseError); ok {

			log.Println("close:", err)

			os.Exit(0)
			break
		}
		if err != nil {
			log.Println("read:", err)
			return
		}
		if mt == websocket.TextMessage {
			var data TextMessage
			//parse json data

			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Println("read:", err)
				//return
				log.Printf("ignored recv text: %s", message)
				return
			}

			if data.Type == "rejected" {
				log.Println("rejected:", data.Body)
				os.Exit(0)
				return
				//break
			} else if data.Type == "resolved" {
				log.Println("resolved:", data.Body)
			} else {
				log.Printf("ignored unknown recv text:%v", data)
			}
		} else {
			decoded, undecoded, err := codec.NativeFromBinary(message)
			if len(undecoded) > 0 {

				log.Println("undecoded:", undecoded)

			}
			if err != nil {
				log.Println("decode:", err)

			} else {
				// log.Printf("recv binary: %s", decoded)
				var md = decoded.(map[string]interface{})
				if md["type"] == "stderr" {
					var body = md["body"].([]byte)

					err_queue.Enqueue(body)
				} else if md["type"] == "stdout" {
					var body = md["body"].([]byte)

					out_queue.Enqueue(body)
				} else {
					log.Println("ignored unknown type:", md["type"])
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

func configureWebSocketTLSCA(x *websocket.Dialer, configdata ConfigClient) {
	x.TLSClientConfig = &tls.Config{
		RootCAs: x509.NewCertPool(),
	}
	caCert, err := os.ReadFile(configdata.Servers.Ca)
	if err != nil {
		log.Fatal(err)
	}
	ok := x.TLSClientConfig.RootCAs.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("Failed to append CA certificate")
	}
}
