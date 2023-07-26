package telegram

type Conversation struct {
	typeOfConversation typeOfConversation
	step               int
}

type typeOfConversation string

const (
	rigesterType typeOfConversation = "register"
)
