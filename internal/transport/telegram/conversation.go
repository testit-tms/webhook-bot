package telegram

// Conversation represents a conversation with a user.
type Conversation struct {
	typeOfConversation typeOfConversation
	step               int
}

type typeOfConversation string

const (
	rigesterType typeOfConversation = "register"
)
