package humble

import (
	"math/rand"
	"strconv"
	"time"
)

type Identifier struct {
	Id string
}

func (i *Identifier) GetId() string {
	if i.Id == "" {
		i.Id = generateRandomId()
	}
	return i.Id
}

// generateRandomId generates a random string that is more or less
// garunteed to be unique. Used as ids for records where an id is
// not otherwise provided.
func generateRandomId() string {
	timeInt := time.Now().Unix()
	timeString := strconv.FormatInt(timeInt, 36)
	randomString := generateRandomAlphanum(16)
	return randomString + timeString
}

func generateRandomAlphanum(size int) string {
	var alphanum []byte = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	var slice []byte = make([]byte, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(slice)
}
