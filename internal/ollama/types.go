package ollama

type ChatRole string

const (
	RoleUser   ChatRole = "user"
	RoleSystem ChatRole = "system"
)

type ChatMessage struct {
	Role ChatRole
	Text string
}

type ChatReply struct {
	Err  error
	Text string
}
