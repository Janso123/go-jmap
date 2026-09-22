package websocket

import (
	jsonv2 "encoding/json/v2"
	"fmt"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
)

// PushEnable is a WebSocketPushEnable client frame (RFC 8887 §4.3.5.2).
// DataTypes nil (JSON null) means all types; a non-nil empty slice means none.
type PushEnable struct {
	Type      string            `json:"@type"`
	DataTypes *[]jmap.EventType `json:"dataTypes"`
	PushState string            `json:"pushState,omitzero"`
}

// PushDisable is a WebSocketPushDisable client frame (RFC 8887 §4.3.5.3).
type PushDisable struct {
	Type string `json:"@type"`
}

// serverFrame is a decoded server→client WebSocket message.
type serverFrame struct {
	RequestID     string
	Response      *jmap.Response
	StateChange   *jmap.StateChange
	RequestError  *jmap.RequestError
	CalendarAlert *calendar.CalendarAlert
}

func marshalRequest(id string, req *jmap.Request) ([]byte, error) {
	type wsReq struct {
		Type       string              `json:"@type"`
		ID         string              `json:"id,omitzero"`
		Using      []jmap.URI          `json:"using"`
		Calls      []*jmap.Invocation  `json:"methodCalls"`
		CreatedIDs map[jmap.ID]jmap.ID `json:"createdIds,omitzero"`
	}
	if err := req.ValidateArgs(); err != nil {
		return nil, err
	}
	return jsonv2.Marshal(wsReq{
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
	if err := jsonv2.Unmarshal(data, &probe); err != nil {
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
		if err := jsonv2.Unmarshal(data, &wrap); err != nil {
			return nil, err
		}
		fr.RequestID = wrap.RequestID
		resp := wrap.Response
		fr.Response = &resp
		return fr, nil
	case "StateChange":
		sc := &jmap.StateChange{}
		if err := jsonv2.Unmarshal(data, sc); err != nil {
			return nil, err
		}
		fr.StateChange = sc
		return fr, nil
	case "RequestError":
		var re jmap.RequestError
		if err := jsonv2.Unmarshal(data, &re); err != nil {
			return nil, err
		}
		fr.RequestID = re.RequestID
		fr.RequestError = &re
		return fr, nil
	case "CalendarAlert":
		alert := &calendar.CalendarAlert{}
		if err := jsonv2.Unmarshal(data, alert); err != nil {
			return nil, err
		}
		fr.CalendarAlert = alert
		return fr, nil
	default:
		return nil, fmt.Errorf("websocket: unknown frame @type %q", probe.Type)
	}
}

func pushEnableFrame(dataTypes []jmap.EventType, pushState string) PushEnable {
	fr := PushEnable{Type: "WebSocketPushEnable", PushState: pushState}
	if dataTypes != nil {
		cp := append([]jmap.EventType(nil), dataTypes...)
		fr.DataTypes = &cp
	}
	return fr
}
