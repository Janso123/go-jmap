package emailsubmission_test

import (
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail/emailsubmission"
	"github.com/stretchr/testify/require"
)

func TestSubmissionParametersTyped(t *testing.T) {
	a := emailsubmission.Address{Parameters: jmap.Some(map[string]*string{"NOTIFY": nil})}
	b, err := jsonv2.Marshal(a)
	require.NoError(t, err)
	require.Contains(t, string(b), `"NOTIFY":null`)
}
