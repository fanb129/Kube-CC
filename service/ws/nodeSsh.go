package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

func NodeWsSsh(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("Sec-WebSocket-Key")
	fmt.Println(id)
	webssh := NewWebSSH()
	// term 可以使用 ansi, linux, vt100, xterm, dumb，除了 dumb外其他都有颜色显示, 默认 xterm
	webssh.SetTerm(TermXterm)
	webssh.SetBuffSize(8192)
	webssh.SetId(id)
	webssh.SetConnTimeOut(5 * time.Second)
	webssh.SetLogger(log.New(os.Stderr, "[webssh] ", log.Ltime|log.Ldate))

	// 是否启用 sz 与 rz
	//webssh.DisableSZ()
	//webssh.DisableRZ()

	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		//Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")},
		Subprotocols:    []string{"webssh"},
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
	}

	ws, err := upGrader.Upgrade(w, r, nil)

	if err != nil {
		zap.S().Errorln(err)
	}

	//ws.SetCompressionLevel(4)
	//ws.EnableWriteCompression(true)

	webssh.AddWebsocket(ws)
}
