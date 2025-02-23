package presenter

type RoomResponse struct {
	Id string `json:"id"`
}

type RoomRequest struct {
	Title      string `json:"title"`
	MaxPlayers int    `json:"maxPlayers"`
}

type JoinRoomRequest struct {
	Id string `json:"id"`
}

type LeaveRoomRequest struct {
	Id string `json:"id"`
}
