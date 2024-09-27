package structs

type TokenData struct {
	UserID     string `json:"userId"`
	AcccountID string `json:"accountId"`
	Exp        string `json:"exp"`
	Iat        string `json:"Iat"`
}

type CreateRoomRequest struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RoomJoinRequest struct {
	RoomID   string      `json:"roomId"`
	UserID   string      `json:"userId"`
	Password string      `json:"password"`
	Offer    interface{} `json:"offer"`
}

type CreateRoomResponse struct {
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
}

type LeftRoomRequest struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}

type KickRoomEvent struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}
