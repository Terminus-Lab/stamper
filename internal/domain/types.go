package domain

type Turn struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
}

type Conversation struct {
	ConversationID string         `json:"conversation_id"`
	Turns          []Turn         `json:"turns"`
	Annotation     string         `json:"human_annotation,omitempty"`
	Extra          map[string]any `json:"-"` //TODO: handle passthrough latter
}
