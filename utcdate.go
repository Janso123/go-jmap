package jmap

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"
	"strings"
	"time"
)

// UTCDate is an RFC 8620 §1.4 UTCDate: RFC 3339 date-time with time-offset Z.
// time-secfrac is omitted when zero and included when non-zero.
type UTCDate time.Time

func (d UTCDate) MarshalJSONTo(enc *jsontext.Encoder) error {
	t := time.Time(d).UTC()
	if t.Nanosecond() == 0 {
		return jsonv2.MarshalEncode(enc, t.Format("2006-01-02T15:04:05Z"))
	}
	return jsonv2.MarshalEncode(enc, t.Format("2006-01-02T15:04:05.999999999Z"))
}

func (d *UTCDate) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s string
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	if strings.ContainsAny(s, "tz") {
		return fmt.Errorf("jmap: UTCDate %q uses lowercase time letters", s)
	}
	if !strings.HasSuffix(s, "Z") {
		return fmt.Errorf("jmap: UTCDate %q is not Z-terminated", s)
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*d = UTCDate(t.UTC())
	return nil
}

// Date is RFC 8620 §1.4 Date: an RFC 3339 date-time with the original
// time-offset preserved (not forced to Z). time-secfrac is omitted when zero.
type Date time.Time

func (d Date) MarshalJSONTo(enc *jsontext.Encoder) error {
	t := time.Time(d)
	if t.Nanosecond() == 0 {
		return jsonv2.MarshalEncode(enc, t.Format("2006-01-02T15:04:05Z07:00"))
	}
	return jsonv2.MarshalEncode(enc, t.Format("2006-01-02T15:04:05.999999999Z07:00"))
}

func (d *Date) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s string
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	if strings.ContainsAny(s, "tz") {
		return fmt.Errorf("jmap: Date %q uses lowercase time letters", s)
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}
