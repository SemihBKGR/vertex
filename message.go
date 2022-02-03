package main

const (
	actionJoin    = "join"
	actionLeave   = "leave"
	actionJoined  = "joined"
	actionLeft    = "left"
	actionMatched = "matched"
)

type message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

func joinedMessage() *message {
	return &message{
		Action: actionJoined,
	}
}

func leftMessage() *message {
	return &message{
		Action: actionLeft,
	}
}

func matchedMessage() *message {
	return &message{
		Action: actionMatched,
	}
}
