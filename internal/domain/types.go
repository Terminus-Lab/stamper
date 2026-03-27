package domain

import "encoding/json"

type Turn struct {
	UserQuery string `json:"user_query"`
	Answer    string `json:"answer"`
}

type Conversation struct {
	ConversationID string                     `json:"conversation_id"`
	Turns          []Turn                     `json:"turns"`
	Annotation     string                     `json:"human_annotation,omitempty"`
	Reason         string                     `json:"human_reason,omitempty"`
	Extra          map[string]json.RawMessage `json:"-"`
}

func (c *Conversation) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// keep full passthrough copy first
	c.Extra = make(map[string]json.RawMessage, len(raw))
	for k, v := range raw {
		c.Extra[k] = v
	}

	// extract typed fields if present
	if v, ok := raw["conversation_id"]; ok {
		if err := json.Unmarshal(v, &c.ConversationID); err != nil {
			return err
		}
		delete(c.Extra, "conversation_id")
	}

	if v, ok := raw["turns"]; ok {
		if err := json.Unmarshal(v, &c.Turns); err != nil {
			return nil
		}

		delete(c.Extra, "turns")
	}

	if v, ok := raw["human_annotation"]; ok {
		if err := json.Unmarshal(v, &c.Annotation); err != nil {
			return nil
		}
		delete(c.Extra, "human_annotation")
	}

	if v, ok := raw["human_reason"]; ok {
		if err := json.Unmarshal(v, &c.Reason); err != nil {
			return nil
		}
		delete(c.Extra, "human_reason")
	}

	return nil

}

func (c *Conversation) MarshalJSON() ([]byte, error) {
	out := make(map[string]json.RawMessage, len(c.Extra)+3)

	for k, v := range c.Extra {
		out[k] = v
	}

	if b, err := json.Marshal(c.ConversationID); err != nil {
		return nil, err
	} else {
		out["conversation_id"] = b
	}

	if b, err := json.Marshal(c.Turns); err != nil {
		return nil, err
	} else {
		out["turns"] = b
	}

	if c.Annotation != "" {
		if b, err := json.Marshal(c.Annotation); err != nil {
			return nil, err
		} else {
			out["human_annotation"] = b
		}
	}

	if c.Reason != "" {
		if b, err := json.Marshal(c.Reason); err != nil {
			return nil, err
		} else {
			out["human_reason"] = b
		}
	}

	return json.Marshal(out)
}
