package structs

type Response struct {
	Sdp      string `json:"sdp"`
	RoomName string `json:"roomName"`
	Uid      string `json:"uid"`
}
