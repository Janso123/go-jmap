package websocket

import (
	"encoding/json"
	"fmt"

	"git.sr.ht/~rockorager/go-jmap"
)

// PushEnable is a WebSocketPushEnable client frame (RFC 8887 §4.3.5.2).
type PushEnable struct {
	Type      string           `json:"@type"`
	DataTypes []jmap.EventType `json:"dataTypes"`
	PushState string           `json:"pushState,omitempty"`
}

// PushDisable is a WebSocketPushDisable client frame (RFC 8887 §4.3.5.3).
type PushDisable struct {
	Type string `json:"@type"`
}

// serverFrame is a decoded server→client WebSocket message.
type serverFrame struct {
	RequestID    string
	Response     *jmap.Response
	StateChange  *jmap.StateChange
	RequestError *jmap.RequestError
}

func marshalRequest(id string, req *jmap.Request) ([]byte, error) {
	type wsReq struct {
		Type       string              `json:"@type"`
		ID         string              `json:"id,omitempty"`
		Using      []jmap.URI          `json:"using"`
		Calls      []*jmap.Invocation  `json:"methodCalls"`
		CreatedIDs map[jmap.ID]jmap.ID `json:"createdIds,omitempty"`
	}
	return json.Marshal(wsReq{
		Type:       "Request",
		ID:         id,
		Using:      req.Using,
		Calls:      req.Calls,
		CreatedIDs: req.CreatedIDs,
	})
}

func decodeServerFrame(data []byte) (*serverFrame, error) {
	var probe struct {
		Type string `json:"@type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	fr := &serverFrame{}
	switch probe.Type {
	case "Response":
		var wrap struct {
			Type      string `json:"@type"`
			RequestID string `json:"requestId"`
			jmap.Response
		}
		if err := json.Unmarshal(data, &wrap); err != nil {
			return nil, err
		}
		fr.RequestID = wrap.RequestID
		resp := wrap.Response
		fr.Response = &resp
		return fr, nil
	case "StateChange":
		sc := &jmap.StateChange{}
		if err := json.Unmarshal(data, sc); err != nil {
			return nil, err
		}
		fr.StateChange = sc
		return fr, nil
	case "RequestError":
		var wrap struct {
			AtType    string  `json:"@type"`
			RequestID string  `json:"requestId"`
			Type      string  `json:"type"`
			Status    int     `json:"status"`
			Detail    string  `json:"detail"`
			Limit     *string `json:"limit"`
		}
		if err := json.Unmarshal(data, &wrap); err != nil {
			return nil, err
		}
		fr.RequestID = wrap.RequestID
		fr.RequestError = &jmap.RequestError{
			Type:   wrap.Type,
			Status: wrap.Status,
			Detail: wrap.Detail,
			Limit:  wrap.Limit,
		}
		return fr, nil
	default:
		return nil, fmt.Errorf("websocket: unknown frame @type %q", probe.Type)
	}
}
