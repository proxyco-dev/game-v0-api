package presenter

type RoomResponse struct {
	Id string `json:"id"`
}

type RoomRequest struct {
	Title   string `json:"title"`
	Private bool   `json:"private"`
	Code    string `json:"code"`
}

type JoinRoomRequest struct {
	Id string `json:"id"`
}
