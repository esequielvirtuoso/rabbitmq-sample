package rabbit

type Message struct {
	Message string `json:"message" binding:"required"`
}
