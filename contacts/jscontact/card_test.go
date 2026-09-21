package jscontact

import (
	jsonv2 "encoding/json/v2"
	"testing"

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
