package presenter

type RoomResponse struct {
	Id string `json:"id"`
}

type RoomRequest struct {
	Title string `json:"title"`
}

type JoinRoomRequest struct {
	Id string `json:"id"`
}
