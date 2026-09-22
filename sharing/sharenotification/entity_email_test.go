package sharenotification

import (
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/stretchr/testify/require"
)

func TestEntityEmailNullRemarshals(t *testing.T) {
	t.Parallel()
	var e Entity
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"name":"x","email":null}`), &e))
	require.True(t, e.Email.IsNull())

	b, err := jsonv2.Marshal(&e)
	require.NoError(t, err)
	require.JSONEq(t, `{"name":"x","email":null}`, string(b))
}

func TestEntityEmailEmptyDistinctFromNull(t *testing.T) {
	t.Parallel()
	var e Entity
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"name":"x","email":""}`), &e))
	require.False(t, e.Email.IsNull())
	v, ok := e.Email.Value()
	require.True(t, ok)
	require.Equal(t, "", v)
}

func TestEntityPrincipalIDNull(t *testing.T) {
	t.Parallel()
	var e Entity
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"principalId":null}`), &e))
	require.True(t, e.PrincipalID.IsNull())

	b, err := jsonv2.Marshal(&e)
	require.NoError(t, err)
	require.JSONEq(t, `{"principalId":null}`, string(b))
}

func TestShareNotificationRightsNullAndEmpty(t *testing.T) {
	t.Parallel()
	var n ShareNotification
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"oldRights":null,"newRights":{}}`), &n))
	require.True(t, n.OldRights.IsNull())
	rights, ok := n.NewRights.Value()
	require.True(t, ok)
	require.Empty(t, rights)

	b, err := jsonv2.Marshal(&n)
	require.NoError(t, err)
	require.JSONEq(t, `{"oldRights":null,"newRights":{}}`, string(b))
}

func TestFilterConditionNullDates(t *testing.T) {
	t.Parallel()
	var f FilterCondition
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"after":null,"before":null}`), &f))
	require.True(t, f.After.IsNull())
	require.True(t, f.Before.IsNull())

	b, err := jsonv2.Marshal(&f)
	require.NoError(t, err)
	require.JSONEq(t, `{"after":null,"before":null}`, string(b))
}
