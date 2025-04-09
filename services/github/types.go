package github

type Event struct {
	Type      string                 `json:"type"`
	Repo      EventRepo              `json:"repo"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt string                 `json:"created_at"`
}

type EventRepo struct {
	Name string `json:"name"`
}
