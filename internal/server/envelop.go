package server

type Envelop struct {
	Type    EnvelopType `json:"type"`
	Payload interface{} `json:"payload"`
}

type EnvelopType string

const (
	Offer         EnvelopType = "offer"
	Answer        EnvelopType = "answer"
	ICE_Candidate EnvelopType = "ice_candidate"
	JoinRoom      EnvelopType = "join_room"
	LeftRoom      EnvelopType = "left_room"
	KickRoom      EnvelopType = "kick_room"
)
