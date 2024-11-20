package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

func ToInt64(i string) int64 {
	v, _ := strconv.ParseInt(i, 10, 64)
	return v
}

func ToInt(i string) int {
	v, _ := strconv.Atoi(i)
	return v
}

// CombineShares combines the shares to reconstruct the master key
func CombineShares(shares [][]byte) ([]byte, error) {
	masterKey, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}
	return masterKey, nil
}

// func Contains(slice []string, item string) bool {
// 	set := make(map[string]struct{}, len(slice))
// 	for _, s := range slice {
// 		set[s] = struct{}{}
// 	}

// 	_, ok := set[item]
// 	return ok
// }

func Contains(actions []string, action string) bool {
	for _, a := range actions {
		if strings.EqualFold(a, action) {
			return true
		}
	}
	return false
}

func MD5Checksum(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
