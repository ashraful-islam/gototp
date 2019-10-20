package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"gototp/src/lib"
	"hash"
	"strings"
	"time"
)

type TotpToken struct {
	// public props
	Secret     string
	keyDecoded []byte
	Digits     int
	Step       int
	// internal props
	counter          uint64
	counterBigEndian []byte
	Digest           string
	hmac             []byte
	token            string
	binaryCode       uint64
}

// remove non-base32 characters
func (t *TotpToken) cleanSecret() {
	t.Secret = lib.FixKey(t.Secret)
}

// decodes provided base32 encoded key
func (t *TotpToken) decodeSecret() {
	key, err := base32.StdEncoding.DecodeString(lib.PadBase32(t.Secret))
	lib.CheckErr(err, "Base32 Decode Error:")
	t.keyDecoded = key
}

// Detect correct hmac encoding to use
// From RFC-6238:
// 		sha1 - secret - 20 bytes
// 		sha256 - secret - 32 bytes
// 		sha512 - secret - 64 bytes
func (t *TotpToken) detectDigest() {
	size := len(t.keyDecoded)

	if size > 64 {
		lib.CheckErr(errors.New("Invalid size detected, boundary condition void"), "Digest Error:")
	}

	if size < 21 {
		t.Digest = "sha1"
	} else if size < 33 {
		t.Digest = "sha256"
	} else {
		t.Digest = "sha512"
	}
}

// converts counter to big endian byte array
func (t *TotpToken) convertCounter() {
	// calculate time offset in bigendian byte array
	t.counterBigEndian = make([]byte, 8) // uint64
	binary.BigEndian.PutUint64(t.counterBigEndian, t.counter)
}

func (t *TotpToken) calculateHmac() {
	var h hash.Hash
	switch t.Digest {
	case "sha1":
		h = hmac.New(sha1.New, t.keyDecoded)
	case "sha256":
		h = hmac.New(sha256.New, t.keyDecoded)
	case "sha512":
		h = hmac.New(sha512.New, t.keyDecoded)
	default:
		// unknown
		lib.CheckErr(fmt.Errorf("Unknown digest: %v", t.Digest), "HMAC Error:")
	}
	// append counter
	h.Write(t.counterBigEndian)
	// generate sum
	t.hmac = h.Sum(nil)
}

func (t *TotpToken) calculateBinaryCode() {
	offset := t.hmac[len(t.hmac)-1] & 0x0f
	binaryCode := uint64(binary.BigEndian.Uint32(t.hmac[offset : offset+4]))
	t.binaryCode = binaryCode & 0x7fffffff
}

func (t *TotpToken) calculateToken() {
	token := make([]byte, t.Digits)
	sToken := fmt.Sprint(t.binaryCode)
	size := len(sToken)
	// pad using 0s when not enough digits
	if size < t.Digits {
		sToken = strings.Repeat("0", t.Digits-size) + sToken
	}
	// copy relevant digits
	startIndex := size - t.Digits
	copy(token, sToken[startIndex:])
	t.token = string(token)
}

// set some initial values and calculations
func (t *TotpToken) Generate() string {
	// calculate current counter
	t.counter = uint64(time.Now().Unix() / int64(t.Step))

	// clean up and decode secret
	t.cleanSecret()
	t.decodeSecret()
	// detect proper digest if none provided
	if t.Digest == "" {
		t.detectDigest()
	}
	// prepare counter
	t.convertCounter()
	// prepare hmac
	t.calculateHmac()
	// prepare binary code
	t.calculateBinaryCode()
	// generate token
	t.calculateToken()
	return t.token
}
