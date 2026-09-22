package jmap

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json/jsontext"
)

// UnsignedInt is an RFC 8620 §1.2 UnsignedInt: an integer in [0, 2^53-1].
type UnsignedInt uint64

// Int is an RFC 8620 §1.2 Int: an integer in [-(2^53-1), 2^53-1].
type Int int64

const (
	MaxUnsignedInt UnsignedInt = 1<<53 - 1
	MaxInt         Int         = 1<<53 - 1
	MinInt         Int         = -MaxInt
)

func (u UnsignedInt) MarshalJSONTo(enc *jsontext.Encoder) error {
	if u > MaxUnsignedInt {
		return fmt.Errorf("jmap: UnsignedInt %d exceeds 2^53-1 (RFC 8620 §1.3)", uint64(u))
	}
	return enc.WriteToken(jsontext.Uint(uint64(u)))
}

func (u *UnsignedInt) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '0' {
		return fmt.Errorf("jmap: UnsignedInt: expected number, got %v", tok.Kind())
	}
	s := tok.String()
	if strings.ContainsAny(s, ".eE-") {
		return fmt.Errorf("jmap: UnsignedInt: %q is not a non-negative integer", s)
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil || v > uint64(MaxUnsignedInt) {
		return fmt.Errorf("jmap: UnsignedInt: %q out of range", s)
	}
	*u = UnsignedInt(v)
	return nil
}

func (i Int) MarshalJSONTo(enc *jsontext.Encoder) error {
	if i > MaxInt || i < MinInt {
		return fmt.Errorf("jmap: Int %d outside [-(2^53-1), 2^53-1] (RFC 8620 §1.3)", int64(i))
	}
	return enc.WriteToken(jsontext.Int(int64(i)))
}

func (i *Int) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '0' {
		return fmt.Errorf("jmap: Int: expected number, got %v", tok.Kind())
	}
	s := tok.String()
	if strings.ContainsAny(s, ".eE") {
		return fmt.Errorf("jmap: Int: %q is not an integer", s)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || Int(v) > MaxInt || Int(v) < MinInt {
		return fmt.Errorf("jmap: Int: %q out of range", s)
	}
	*i = Int(v)
	return nil
}
