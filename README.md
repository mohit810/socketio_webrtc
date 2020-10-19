# Signalling Pion Webrtc server via `Socket.io`

I have build  this small signalling upgradation for anyone intrested in building the thier own server to establish 1:N Broadcasting connection with the users as described `broadcast is a Pion WebRTC application that demonstrates how to broadcast a video to many peers, while only requiring the broadcaster to upload once.` in the [Pion/webrtc/example/broadcast/](https://github.com/pion/webrtc/tree/master/examples/broadcast).

Before using this solution you should set-up have pion/webrtc/v3 ([Go Modules](https://blog.golang.org/using-go-modules) are mandatory for using Pion WebRTC. So make sure you set export GO111MODULE=on, and explicitly specify /v3 when importing.).

### Open broadcast example page
[localhost:8000](http://http://localhost:8000/) You should see two buttons 'Publish a Broadcast' and 'Join a Broadcast'. `Keep in Mind that Publisher first needs to start the stream and enter a room name and after the successful conection the users can join the room by just entering the room name that the broadcaster used`

### Run Application
#### Linux/macOS/windows
Run `main.go`

## Big Thanks to the following 

* [Sean Der](https://github.com/Sean-Der) at [Poin/webrtc](https://github.com/pion/webrtc)
* [go-socket.io](https://github.com/googollee/go-socket.io)
