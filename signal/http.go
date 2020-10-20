package signal

import (
	"flag"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"net/http"
	"socketio_webrtc/structs"
	"strconv"
)

var (
	SocketServer *socketio.Server
	answer       string
)

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer() (offerOut chan string, answerIn chan string) {
	port := flag.Int("port", 8000, "http server port")
	flag.Parse()

	offerOut = make(chan string)
	answerIn = make(chan string)
	SocketServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	SocketServer.OnEvent("/", "createConnection", func(s socketio.Conn, roomName string) string {
		if SocketServer.RoomLen("/", roomName) == 0 {
			SocketServer.JoinRoom("/", roomName, s)
			responsePkt := structs.Response{
				RoomName: roomName,
				Uid:      s.ID(),
			}
			s.Emit("created", responsePkt)
		} else {
			SocketServer.JoinRoom("/", roomName, s)
			responsePkt := structs.Response{
				RoomName: roomName,
				Uid:      s.ID(),
			}
			s.Emit("joined", responsePkt)
		}
		return "recv " + roomName
	})

	SocketServer.OnEvent("/", "ready", func(s socketio.Conn, roomName string) string {
		booll := SocketServer.BroadcastToRoom("/", roomName, "ready", roomName)
		fmt.Print("\nready socket : ", booll)
		return "recv " + roomName
	})

	SocketServer.OnEvent("/", "offer", func(s socketio.Conn, data structs.Offer) string {
		offerOut <- data.Sdp
		answer = <-answerIn
		responsePkt := structs.Response{
			Sdp:      answer,
			RoomName: data.RoomName,
			Uid:      s.ID(),
		}
		s.Emit("answer", responsePkt)
		return "done"
	})

	SocketServer.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	SocketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go func() {
		defer SocketServer.Close()
		http.Handle("/socket.io/", SocketServer)
		http.Handle("/", http.FileServer(http.Dir("./public")))
		err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if err != nil {
			panic(err)
		}
	}()

	return
}
