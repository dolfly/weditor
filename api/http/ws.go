package http

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/dolfly/weditor/web"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Method string      `json:"method"`
	Value  interface{} `json:"value"`
}

func prepare(callback func(tty *os.File)) (tty *os.File, pycmd *exec.Cmd, err error) {
	pycmd = exec.Command("/usr/local/bin/python", "-u", web.TempScript())
	pycmd.Env = append(os.Environ(),
		"PYTHONIOENCODING=utf-8")
	tty, err = pty.Start(pycmd)
	if err != nil {
		return
	}
	go callback(tty)
	return
}

func Python(c *gin.Context) {
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	var callback = func(tty *os.File) {
		if tty == nil {
			return
		}
		scanner := bufio.NewScanner(tty)
		for scanner.Scan() {
			text := scanner.Text()
			arr := strings.Split(text, ":")
			if len(arr) < 2 {
				continue
			}
			cmdx, value := arr[0], strings.Join(arr[1:], ":")
			switch cmdx {
			case "LNO":
				ws.WriteJSON(Message{
					Method: "gotoLine",
					Value:  value,
				})
			case "DBG":
				// ws.WriteJSON(Message{
				// 	Method: "output",
				// 	Value:  "-" + value + "\n",
				// })
			case "WRT":
				// ws.WriteJSON(Message{
				// 	Method: "output",
				// 	Value:  "> " + value + "\n",
				// })
			case "EOF":
				ws.WriteJSON(Message{
					Method: "finish",
					Value:  value,
				})
			default:
				ws.WriteJSON(Message{
					Method: "output",
					Value:  "\n",
				})
			}
		}
	}
	tty, pycmd, err := prepare(callback)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		pycmd.Process.Kill()
		pycmd.Process.Wait()
		ws.Close()
		tty.Close()
	}()
	adjust := func(code interface{}) []byte {
		str := code.(string)
		return []byte(str + "\n")
	}
	for {
		var (
			msg Message
			err error
		)
		msg = Message{}
		ws.ReadJSON(&msg)
		if err != nil {
			break
		}
		switch strings.ToLower(msg.Method) {
		case "input":
			tty.Write(adjust(msg.Value))
		case "keyboardinterrupt":
			if pycmd.Process != nil {
				pycmd.Process.Signal(syscall.SIGINT)
			}
		case "restartkernel":
			tty.Close()
			if pycmd.Process != nil {
				err = pycmd.Process.Kill()
			}
			if err != nil {
				ws.WriteJSON(Message{
					Method: "output",
					Value:  "tty close:" + err.Error(),
				})
			}
			time.Sleep(3 * time.Second)
			tty, pycmd, err = prepare(callback)
			if err == nil {
				ws.WriteJSON(Message{
					Method: "restarted",
					Value:  "success",
				})
			} else {
				ws.WriteJSON(Message{
					Method: "output",
					Value:  err.Error(),
				})
			}
			// if err != nil {
			// 	ws.WriteJSON(Message{
			// 		Method: "output",
			// 		Value:  "tty close:" + err.Error(),
			// 	})
			// }
			// } else {
			// 	tty, err = pty.Start(pycmd)
			// 	if err != nil {
			// 		ws.WriteJSON(Message{
			// 			Method: "output",
			// 			Value:  "pty:" + err.Error(),
			// 		})
			// 	}
			// }
			// if pycmd.Process != nil {
			// 	err := pycmd.Process.Kill()
			// 	if err == nil {
			// 		ws.WriteJSON(Message{
			// 			Method: "restarted",
			// 			Value:  "success",
			// 		})
			// 	} else {
			// 		ws.WriteJSON(Message{
			// 			Method: "restarted",
			// 			Value:  err.Error(),
			// 		})
			// 	}
			// } else {
			// 	ws.WriteJSON(Message{
			// 		Method: "restarted",
			// 		Value:  "not exec",
			// 	})
			// }
		default:
			ws.WriteJSON(Message{
				Method: "default",
				Value:  nil,
			})
		}
	}
}
