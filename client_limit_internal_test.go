package jmap

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/stretchr/testify/require"
)

func TestResponseLimit(t *testing.T) {
	for name, tc := range map[string]struct {
		session *Session
		want    int64
	}{
		"nil session":        {nil, maxJSONBody},
		"no core capability": {&Session{}, maxJSONBody},
		"nil core capability": {
			&Session{Capabilities: map[URI]Capability{CoreURI: (*Core)(nil)}},
			maxJSONBody,
		},
		"below default does not lower": {
			&Session{Capabilities: map[URI]Capability{CoreURI: &Core{MaxSizeRequest: 1 << 20}}},
			maxJSONBody,
		},
		"above default raises": {
			&Session{Capabilities: map[URI]Capability{CoreURI: &Core{MaxSizeRequest: 50000000}}},
			50000000,
		},
		"huge value is clamped": {
			&Session{Capabilities: map[URI]Capability{CoreURI: &Core{MaxSizeRequest: MaxUnsignedInt}}},
			maxSessionJSONBody,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := &Client{Session: tc.session}
			require.Equal(t, tc.want, c.responseLimit())
		})
	}
}

func TestResponseLimitExplicitCapWins(t *testing.T) {
	c := &Client{
		maxResponseBytes: 1 << 30,
		Session:          &Session{Capabilities: map[URI]Capability{CoreURI: &Core{MaxSizeRequest: 1 << 20}}},
	}
	require.Equal(t, int64(1<<30), c.responseLimit())
}

func TestSessionDecodesCoreCapability(t *testing.T) {
	var s Session
	err := jsonv2.Unmarshal([]byte(`{"capabilities":{"urn:ietf:params:jmap:core":{"maxSizeRequest":50000000}}}`), &s)
	require.NoError(t, err)
	core, ok := s.Capabilities[CoreURI].(*Core)
	require.True(t, ok)
	require.Equal(t, UnsignedInt(50000000), core.MaxSizeRequest)
}

func TestSessionRejectsMalformedCoreCapability(t *testing.T) {
	for _, v := range []string{`"50000000"`, `5e7`, `-1`} {
		var s Session
		err := jsonv2.Unmarshal([]byte(`{"capabilities":{"urn:ietf:params:jmap:core":{"maxSizeRequest":`+v+`}}}`), &s)
		require.Error(t, err, v)
	}
}
