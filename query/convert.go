package query

import (
	"encoding/base64"
	"encoding/hex"
)

var (
	hexSlice = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
)

// HashToUID converts the client.Base64Hash to the client UID
// This is useful to find the client an avatar belongs to
func (a Agent) HashToUID(hash string) (string, error) {
	b, err := hex.DecodeString(convertTsHexToHex(hash))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// UIDToHash converts a client.UID to a Base64Hash
func (a Agent) UIDToHash(uid string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(uid)
	if err != nil {
		return "", err
	}
	return convertHexToTsHex(hex.EncodeToString(b)), nil
}

func convertHexToTsHex(h string) string {
	result := ""
	for _, v := range h {
		result += string(hexSlice[findIndex(false, v)-16])
	}
	return result
}

func convertTsHexToHex(h string) string {
	result := ""
	for _, v := range h {
		result += string(hexSlice[findIndex(true, v)+16])
	}
	return result
}

func findIndex(direction bool, r rune) int {
	if direction {
		for i := 0; i < len(hexSlice); i++ {
			if hexSlice[i] == r {
				return i
			}
		}
	} else {
		for i := len(hexSlice) - 1; i > 0; i-- {
			if hexSlice[i] == r {
				return i
			}
		}
	}
	return -1
}
