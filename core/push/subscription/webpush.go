package subscription

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

var (
	errShort   = errors.New("subscription: short aes128gcm body")
	errPadding = errors.New("subscription: bad aes128gcm padding")
)

// GenerateKeys creates a P-256 subscription key and a 16-byte auth secret.
// Public and Auth are raw URL base64 (no padding).
func GenerateKeys() (Key, *ecdh.PrivateKey, error) {
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return Key{}, nil, err
	}
	auth := make([]byte, 16)
	if _, err := rand.Read(auth); err != nil {
		return Key{}, nil, err
	}
	return Key{
		Public: base64.RawURLEncoding.EncodeToString(priv.PublicKey().Bytes()),
		Auth:   base64.RawURLEncoding.EncodeToString(auth),
	}, priv, nil
}

// OpenPush decrypts an aes128gcm Web Push body (RFC 8291) and decodes the
// JMAP object by @type. priv/auth come from GenerateKeys. A body that is not
// encrypted is rejected when priv != nil (RFC 8620 §8.7).
func OpenPush(h http.Header, body []byte, priv *ecdh.PrivateKey, auth []byte) (any, error) {
	if priv != nil && h.Get("Content-Encoding") != "aes128gcm" {
		return nil, errors.New("subscription: unencrypted push rejected; keys were set")
	}
	plain := body
	if priv != nil {
		var err error
		if plain, err = openAES128GCM(body, priv, auth); err != nil {
			return nil, err
		}
	}
	var probe struct {
		Type string `json:"@type"`
	}
	if err := jsonv2.Unmarshal(plain, &probe); err != nil {
		return nil, err
	}
	switch probe.Type {
	case "PushVerification":
		v := &Verification{}
		return v, jsonv2.Unmarshal(plain, v)
	case "StateChange":
		s := &jmap.StateChange{}
		return s, jsonv2.Unmarshal(plain, s)
	}
	return nil, fmt.Errorf("subscription: unknown push @type %q", probe.Type)
}

func openAES128GCM(body []byte, priv *ecdh.PrivateKey, auth []byte) ([]byte, error) {
	// header: salt[16] | rs[4] big-endian | idlen[1] | keyid (as_public, 65 bytes uncompressed P-256)
	if len(body) < 21 {
		return nil, errShort
	}
	salt := body[:16]
	rs := binary.BigEndian.Uint32(body[16:20])
	idlen := int(body[20])
	if len(body) < 21+idlen {
		return nil, errShort
	}
	asPub, err := ecdh.P256().NewPublicKey(body[21 : 21+idlen])
	if err != nil {
		return nil, err
	}
	records := body[21+idlen:]

	shared, err := priv.ECDH(asPub) // ecdh_secret
	if err != nil {
		return nil, err
	}
	uaPub := priv.PublicKey().Bytes()
	// RFC 8291 §3.3–3.4
	keyInfo := append(append([]byte("WebPush: info\x00"), uaPub...), asPub.Bytes()...)
	ikm := hkdfExpand(hkdfExtract(auth, shared), keyInfo, 32)
	// RFC 8188 §2.2
	prk := hkdfExtract(salt, ikm)
	cek := hkdfExpand(prk, []byte("Content-Encoding: aes128gcm\x00"), 16)
	nonceBase := hkdfExpand(prk, []byte("Content-Encoding: nonce\x00"), 12)

	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	var out []byte
	for seq := uint64(0); len(records) > 0; seq++ {
		n := min(int(rs), len(records))
		if n == 0 {
			return nil, errShort
		}
		nonce := make([]byte, 12)
		copy(nonce, nonceBase)
		binary.BigEndian.PutUint64(nonce[4:], binary.BigEndian.Uint64(nonce[4:])^seq)
		pt, err := gcm.Open(nil, nonce, records[:n], nil)
		if err != nil {
			return nil, err
		}
		// strip padding: last non-zero byte is 0x02 (final record) or 0x01
		i := len(pt) - 1
		for i >= 0 && pt[i] == 0 {
			i--
		}
		if i < 0 || (pt[i] != 0x01 && pt[i] != 0x02) {
			return nil, errPadding
		}
		out = append(out, pt[:i]...)
		records = records[n:]
	}
	return out, nil
}

// sealPush encrypts plaintext as one aes128gcm record (RFC 8291 / RFC 8188).
// uaPub is the subscription public key, auth is the 16-byte auth secret, and
// rs is 4096. The record delimiter is 0x02.
func sealPush(plaintext, uaPub, auth []byte) ([]byte, error) {
	as, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	ua, err := ecdh.P256().NewPublicKey(uaPub)
	if err != nil {
		return nil, err
	}
	shared, err := as.ECDH(ua)
	if err != nil {
		return nil, err
	}
	asPub := as.PublicKey().Bytes()
	keyInfo := append(append([]byte("WebPush: info\x00"), uaPub...), asPub...)
	ikm := hkdfExpand(hkdfExtract(auth, shared), keyInfo, 32)

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	prk := hkdfExtract(salt, ikm)
	cek := hkdfExpand(prk, []byte("Content-Encoding: aes128gcm\x00"), 16)
	nonce := hkdfExpand(prk, []byte("Content-Encoding: nonce\x00"), 12)

	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	padded := append(append([]byte{}, plaintext...), 0x02)
	ct := gcm.Seal(nil, nonce, padded, nil)

	body := make([]byte, 0, 21+len(asPub)+len(ct))
	body = append(body, salt...)
	var rs [4]byte
	binary.BigEndian.PutUint32(rs[:], 4096)
	body = append(body, rs[:]...)
	body = append(body, byte(len(asPub)))
	body = append(body, asPub...)
	body = append(body, ct...)
	return body, nil
}

// hkdfExtract is HKDF-Extract with SHA-256: HMAC(salt, ikm).
func hkdfExtract(salt, ikm []byte) []byte {
	mac := hmac.New(sha256.New, salt)
	mac.Write(ikm)
	return mac.Sum(nil)
}

// hkdfExpand is one block of HKDF-Expand with SHA-256: HMAC(prk, info || 0x01)[:n].
// n must be at most 32.
func hkdfExpand(prk, info []byte, n int) []byte {
	mac := hmac.New(sha256.New, prk)
	mac.Write(info)
	mac.Write([]byte{0x01})
	return mac.Sum(nil)[:n]
}
