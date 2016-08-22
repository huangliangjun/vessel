package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

var (
	uuidLock *sync.Mutex
	lastNum  int64
)

func init() {
	uuidLock = new(sync.Mutex)
}

// UUID generate unique values
func UUID() string {
	uuidLock.Lock()
	result := time.Now().UnixNano()

	for lastNum == result {
		result = time.Now().UnixNano()
	}
	lastNum = result
	uuidLock.Unlock()
	return MD5String(strconv.Itoa(int(lastNum)))
}

// MD5String MD5 string
func MD5String(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
