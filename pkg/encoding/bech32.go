// Package encoding provides bech32 encoding and qage key parsing.
package encoding

import (
	"errors"
	"strings"
)

// Reference implementation adapted (trimmed) from BIP-0173.

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var charsetRev [128]byte

func init() {
	for i := range charsetRev {
		charsetRev[i] = 0xFF
	}
	for i, c := range charset {
		charsetRev[c] = byte(i)
	}
}

func polymod(values []byte) uint32 {
	chk := uint32(1)
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ uint32(v)
		if (b & 1) != 0 {
			chk ^= 0x3b6a57b2
		}
		if (b & 2) != 0 {
			chk ^= 0x26508e6d
		}
		if (b & 4) != 0 {
			chk ^= 0x1ea119fa
		}
		if (b & 8) != 0 {
			chk ^= 0x3d4233dd
		}
		if (b & 16) != 0 {
			chk ^= 0x2a1462b3
		}
	}
	return chk
}

func hrpExpand(hrp string) []byte {
	out := make([]byte, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		out = append(out, hrp[i]>>5)
	}
	out = append(out, 0)
	for i := 0; i < len(hrp); i++ {
		out = append(out, hrp[i]&31)
	}
	return out
}

func createChecksum(hrp string, data []byte) []byte {
	values := append(hrpExpand(hrp), data...)
	values = append(values, []byte{0, 0, 0, 0, 0, 0}...)
	pm := polymod(values) ^ 1
	cs := make([]byte, 6)
	for i := 0; i < 6; i++ {
		cs[i] = byte((pm >> uint(5*(5-i))) & 31)
	}
	return cs
}

func verifyChecksum(hrp string, data []byte) bool {
	return polymod(append(hrpExpand(hrp), data...)) == 1
}

// ConvertBits groups bits from inWidth to outWidth (no padding if pad=false).
func ConvertBits(data []byte, inWidth, outWidth uint, pad bool) ([]byte, error) {
	var acc uint = 0
	var bits uint = 0
	maxv := (1 << outWidth) - 1
	out := make([]byte, 0, len(data)*int(inWidth)/int(outWidth))
	for _, value := range data {
		if value>>inWidth != 0 {
			return nil, errors.New("bech32: invalid data range")
		}
		acc = (acc << inWidth) | uint(value)
		bits += inWidth
		for bits >= outWidth {
			bits -= outWidth
			out = append(out, byte((acc>>bits)&uint(maxv)))
		}
	}
	if pad {
		if bits > 0 {
			out = append(out, byte((acc<<(outWidth-bits))&uint(maxv)))
		}
	} else if bits >= inWidth || ((acc<<(outWidth-bits))&uint(maxv)) != 0 {
		return nil, errors.New("bech32: invalid padding")
	}
	return out, nil
}

// Encode encodes raw 8-bit data into bech32 string (with conversion) under hrp.
func Encode(hrp string, raw []byte) (string, error) {
	if len(hrp) < 1 || len(hrp) > 83 {
		return "", errors.New("bech32: invalid hrp length")
	}
	for _, c := range hrp {
		if c < 33 || c > 126 || (c >= 'A' && c <= 'Z') {
			return "", errors.New("bech32: invalid hrp char")
		}
	}
	five, err := ConvertBits(raw, 8, 5, true)
	if err != nil {
		return "", err
	}
	checksum := createChecksum(hrp, five)
	combined := make([]byte, 0, len(five)+len(checksum))
	combined = append(combined, five...)
	combined = append(combined, checksum...)
	var sb strings.Builder
	sb.Grow(len(hrp) + 1 + len(combined))
	sb.WriteString(hrp)
	sb.WriteByte('1')
	for _, p := range combined {
		sb.WriteByte(charset[p])
	}
	return sb.String(), nil
}

// Decode decodes string into hrp and 8-bit data.
func Decode(s string) (string, []byte, error) {
	if len(s) < 8 || len(s) > 6000 {
		return "", nil, errors.New("bech32: invalid length")
	}
	// Lowercase only
	for _, c := range s {
		if c < 33 || c > 126 || (c >= 'A' && c <= 'Z') {
			return "", nil, errors.New("bech32: mixed or invalid case")
		}
	}
	pos := strings.LastIndexByte(s, '1')
	if pos < 1 || pos+7 > len(s) {
		return "", nil, errors.New("bech32: invalid separator position")
	}
	hrp := s[:pos]
	dataPart := s[pos+1:]
	data := make([]byte, len(dataPart))
	for i := range dataPart {
		c := dataPart[i]
		if c > 127 || charsetRev[c] == 0xFF {
			return "", nil, errors.New("bech32: invalid charset")
		}
		data[i] = charsetRev[c]
	}
	if !verifyChecksum(hrp, data) {
		return "", nil, errors.New("bech32: bad checksum")
	}
	payload := data[:len(data)-6]
	eight, err := ConvertBits(payload, 5, 8, false)
	if err != nil {
		return "", nil, err
	}
	return hrp, eight, nil
}
