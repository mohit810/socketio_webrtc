package signal

import (
	"flag"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"practice/socketio_webrtc/structs"
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
	/*SocketServer = InitializeSockets()*/
	SocketServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	SocketServer.OnEvent("/", "createConnection", func(s socketio.Conn, roomName string) string {
		/*roomName = append(roomName, msg)*/
		if SocketServer.RoomLen("/", roomName) == 0 {
			SocketServer.JoinRoom("/", roomName, s)
			responsePkt := structs.Response{
				RoomName: roomName,
				Uid:      s.ID(),
			}
			s.Emit("created", responsePkt)
		} else if SocketServer.RoomLen("/", roomName) == 1 {
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

	/*SocketServer.OnEvent("/", "candidate", func(s socketio.Conn, data structs.Candidate) string {
		SocketServer.RoomLen("/", data.Room)
		coolean :=SocketServer.BroadcastToRoom("/", data.Room, "candidate", data)
		fmt.Print("\ncandidate socket : ", coolean)
		return "recv " + data.Room
	})*/

	SocketServer.OnEvent("/", "offer", func(s socketio.Conn, data structs.Offer) string {
		/*SocketServer.BroadcastToRoom("/", data.RoomName, "offer", data.Sdp)*/ //for one-one chatting
		offerOut <- data.Sdp
		answer = <-answerIn
		responsePkt := structs.Response{
			Sdp:      answer,
			RoomName: data.RoomName,
			Uid:      s.ID(),
		}
		SocketServer.BroadcastToRoom("/", responsePkt.RoomName, "answer", responsePkt)
		return "done"
	})

	//use the below socket for one-one chatting
	/*SocketServer.OnEvent("/", "answer", func(s socketio.Conn, data structs.Offer) string {
		SocketServer.BroadcastToRoom("/", data.RoomName, "answer", data.Sdp)
		return "done"
	})*/

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

func InitializeSockets() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	return server
}
