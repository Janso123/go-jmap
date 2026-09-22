package jscontact

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardRoundTripMinimal(t *testing.T) {
	const input = `{
		"@type":"Card",
		"version":"1.0",
		"uid":"urn:uuid:test",
		"name":{"full":"Ada Lovelace"}
	}`

	var card Card
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &card))

	data, err := jsonv2.Marshal(&card)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestCardRoundTripContactMaps(t *testing.T) {
	const input = `{
		"@type":"Card",
		"version":"1.0",
		"uid":"urn:uuid:test",
		"emails":{
			"e1":{
				"address":"ada@example.com",
				"contexts":{"work":true},
				"pref":1,
				"label":"primary"
			}
		},
		"phones":{
			"p1":{
				"number":"tel:+12025550123",
				"features":{"voice":true,"mobile":true},
				"contexts":{"work":true},
				"pref":1,
				"label":"desk"
			}
		},
		"addresses":{
			"a1":{
				"full":"123 Example Street\nLondon",
				"countryCode":"GB",
				"components":[
					{"kind":"number","value":"123"},
					{"kind":"name","value":"Example Street"},
					{"kind":"locality","value":"London"}
				],
				"isOrdered":true,
				"contexts":{"work":true},
				"pref":1
			}
		}
	}`

	var card Card
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &card))

	data, err := jsonv2.Marshal(&card)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestCardRoundTripBroaderProperties(t *testing.T) {
	const input = `{
		"@type":"Card",
		"version":"1.0",
		"uid":"urn:uuid:grace-hopper",
		"created":"2024-01-02T03:04:05Z",
		"updated":"2024-02-03T04:05:06Z",
		"kind":"individual",
		"language":"en",
		"prodId":"go-jmap test",
		"relatedTo":{
			"urn:uuid:friend":{"relation":{"friend":true}}
		},
		"name":{
			"@type":"Name",
			"full":"Grace Hopper",
			"components":[
				{"kind":"given","value":"Grace"},
				{"kind":"surname","value":"Hopper"}
			],
			"isOrdered":true
		},
		"nicknames":{
			"n1":{"name":"Amazing Grace","contexts":{"private":true},"pref":1}
		},
		"organizations":{
			"o1":{
				"name":"US Navy",
				"units":[{"name":"Research"}],
				"sortAs":"Navy",
				"contexts":{"work":true}
			}
		},
		"speakToAs":{
			"grammaticalGender":"feminine",
			"pronouns":{
				"p1":{"pronouns":"she/her","contexts":{"work":true},"pref":1}
			}
		},
		"titles":{
			"t1":{"name":"Rear Admiral","kind":"title","organizationId":"o1"}
		},
		"onlineServices":{
			"s1":{
				"service":"Mastodon",
				"uri":"https://example.com/@grace",
				"user":"@grace",
				"contexts":{"work":true},
				"pref":1,
				"label":"profile"
			}
		},
		"preferredLanguages":{
			"l1":{"language":"en","contexts":{"work":true},"pref":1}
		},
		"calendars":{
			"c1":{"kind":"calendar","uri":"webcal://example.com/grace.ics","pref":1}
		},
		"schedulingAddresses":{
			"sa1":{"uri":"mailto:grace@example.com","contexts":{"work":true},"pref":1,"label":"calendar"}
		},
		"cryptoKeys":{
			"k1":{
				"uri":"https://example.com/grace.asc",
				"mediaType":"application/pgp-keys",
				"contexts":{"work":true},
				"pref":1,
				"label":"pgp"
			}
		},
		"directories":{
			"d1":{"kind":"entry","uri":"https://dir.example.com/grace","pref":1,"listAs":1}
		},
		"links":{
			"link1":{"kind":"contact","uri":"mailto:contact@example.com","pref":1}
		},
		"media":{
			"m1":{
				"kind":"photo",
				"uri":"https://example.com/grace.jpg",
				"mediaType":"image/jpeg",
				"contexts":{"work":true},
				"pref":1,
				"label":"portrait"
			}
		},
		"localizations":{
			"fr":{"name/full":"Grace Hopper"}
		},
		"anniversaries":{
			"a1":{
				"kind":"birth",
				"date":{"year":1906,"month":12,"day":9},
				"place":{"full":"New York, USA"}
			}
		},
		"keywords":{"compiler":true,"navy":true},
		"notes":{
			"note1":{
				"note":"Invented COBOL",
				"created":"2024-03-04T05:06:07Z",
				"author":{"name":"Historian","uri":"mailto:hist@example.com"}
			}
		},
		"personalInfo":{
			"pi1":{
				"kind":"expertise",
				"value":"programming languages",
				"level":"high",
				"listAs":1,
				"label":"specialty"
			}
		}
	}`

	var card Card
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &card))

	data, err := jsonv2.Marshal(&card)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestCardRoundTripVendorExtensions(t *testing.T) {
	const input = `{
		"@type":"Card",
		"version":"1.0",
		"uid":"urn:uuid:vendor-extension",
		"example.com:foo":"bar",
		"name":{
			"full":"Ada Lovelace",
			"example.com:phoneticHint":"AH-duh"
		}
	}`

	var card Card
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &card))

	data, err := jsonv2.Marshal(&card)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestCardRoundTripGroupMembers(t *testing.T) {
	const input = `{
		"@type":"Card",
		"version":"1.0",
		"uid":"urn:uuid:doe-family",
		"kind":"group",
		"name":{"full":"The Doe Family"},
		"members":{
			"urn:uuid:member-1":true,
			"urn:uuid:member-2":true
		}
	}`

	var card Card
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &card))

	data, err := jsonv2.Marshal(&card)
	require.NoError(t, err)
	assert.JSONEq(t, input, string(data))
}

func TestCardExtraRoundTrip(t *testing.T) {
	const input = `{"@type":"Card","uid":"urn:uuid:test","foo":1}`
	var card Card
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &card))
	require.Equal(t, "urn:uuid:test", card.UID)
	raw, ok := card.Extra["foo"]
	require.True(t, ok)
	require.Equal(t, "1", string(raw))
	_, hasUID := card.Extra["uid"]
	require.False(t, hasUID)

	out, err := jsonv2.Marshal(card)
	require.NoError(t, err)
	require.JSONEq(t, input, string(out))
}

func TestMediaBlobIDNotOnlyInExtra(t *testing.T) {
	t.Parallel()
	var m Media
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"uri":"https://ex/photo.jpg","blobId":"b1"}`), &m))
	_, inMediaExtra := m.Extra["blobId"]
	_, inResourceExtra := m.Resource.Extra["blobId"]
	require.False(t, inMediaExtra)
	require.False(t, inResourceExtra)
	require.Equal(t, jmap.ID("b1"), m.BlobID)
	b, err := jsonv2.Marshal(&m)
	require.NoError(t, err)
	require.Contains(t, string(b), `"blobId":"b1"`)
}

func TestLinkHasNoBlobID(t *testing.T) {
	t.Parallel()
	// Link has no BlobID field. blobId round-trips through the single Extra map.
	type linkWire struct {
		Resource
	}
	var l Link
	_ = linkWire(l)

	const input = `{"uri":"https://ex/","blobId":"b2"}`
	require.NoError(t, jsonv2.Unmarshal([]byte(input), &l))
	raw, ok := l.Extra["blobId"]
	require.True(t, ok)
	require.JSONEq(t, `"b2"`, string(raw))
	out, err := jsonv2.Marshal(&l)
	require.NoError(t, err)
	require.JSONEq(t, input, string(out))
}

func TestSeparatorComponentKeepsEmptyValue(t *testing.T) {
	t.Parallel()
	const in = `{"kind":"separator","value":""}`
	for _, name := range []string{"name", "address"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var out []byte
			var err error
			if name == "name" {
				var c NameComponent
				require.NoError(t, jsonv2.Unmarshal([]byte(in), &c))
				out, err = jsonv2.Marshal(&c)
			} else {
				var c AddressComponent
				require.NoError(t, jsonv2.Unmarshal([]byte(in), &c))
				out, err = jsonv2.Marshal(&c)
			}
			require.NoError(t, err)
			require.Contains(t, string(out), `"value":""`)
		})
	}
}

func TestPartialDateYearZeroKept(t *testing.T) {
	t.Parallel()
	const in = `{"year":0,"month":1,"day":1}`
	var date PartialDate
	require.NoError(t, jsonv2.Unmarshal([]byte(in), &date))
	out, err := jsonv2.Marshal(&date)
	require.NoError(t, err)
	require.Contains(t, string(out), `"year":0`)
	require.Contains(t, string(out), `"month":1`)
	require.Contains(t, string(out), `"day":1`)
}

func TestAnniversaryEmptyOmitsDate(t *testing.T) {
	t.Parallel()
	out, err := jsonv2.Marshal(Anniversary{})
	require.NoError(t, err)
	require.NotContains(t, string(out), "date")
}

func TestAnniversaryDateBothNilOmitsDate(t *testing.T) {
	t.Parallel()
	out, err := jsonv2.Marshal(Anniversary{Date: &AnniversaryDate{}})
	require.NoError(t, err)
	require.NotContains(t, string(out), "date")
}

func TestDefaultSeparatorNull(t *testing.T) {
	t.Parallel()
	var name Name
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"defaultSeparator":null}`), &name))
	require.True(t, name.DefaultSeparator.IsNull())
	out, err := jsonv2.Marshal(&name)
	require.NoError(t, err)
	require.Contains(t, string(out), `"defaultSeparator":null`)

	var addr Address
	require.NoError(t, jsonv2.Unmarshal([]byte(`{"defaultSeparator":null}`), &addr))
	require.True(t, addr.DefaultSeparator.IsNull())
	out, err = jsonv2.Marshal(&addr)
	require.NoError(t, err)
	require.Contains(t, string(out), `"defaultSeparator":null`)
}

func TestEmbeddedResourceHasSingleExtra(t *testing.T) {
	t.Parallel()
	raw := jsontext.Value(`1`)

	cal := Calendar{
		Extra: map[string]jsontext.Value{"x": raw}}
	require.Equal(t, raw, cal.Resource.Extra["x"])

	key := CryptoKey{
		Extra: map[string]jsontext.Value{"x": raw}}
	require.Equal(t, raw, key.Resource.Extra["x"])

	dir := Directory{
		Extra: map[string]jsontext.Value{"x": raw}}
	require.Equal(t, raw, dir.Resource.Extra["x"])

	link := Link{
		Extra: map[string]jsontext.Value{"x": raw}}
	require.Equal(t, raw, link.Resource.Extra["x"])

	media := Media{
		Extra: map[string]jsontext.Value{"x": raw}}
	require.Equal(t, raw, media.Resource.Extra["x"])
}
