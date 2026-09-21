package jmap

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"time"
)

// UTCDate is an RFC 8620 UTCDate: always marshaled in UTC with a Z suffix.
type UTCDate time.Time

func (d UTCDate) MarshalJSONTo(enc *jsontext.Encoder) error {
	t := time.Time(d).UTC()
	return jsonv2.MarshalEncode(enc, t.Format("2006-01-02T15:04:05Z"))
}

func (d *UTCDate) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s string
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*d = UTCDate(t.UTC())
	return nil
}

// Date is a local calendar date per RFC 8620 §1.4; marshal as YYYY-MM-DD.
type Date time.Time

func (d Date) MarshalJSONTo(enc *jsontext.Encoder) error {
	t := time.Time(d)
	return jsonv2.MarshalEncode(enc, t.Format("2006-01-02"))
}

func (d *Date) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s string
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}
