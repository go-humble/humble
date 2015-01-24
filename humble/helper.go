package humble

import (
	"github.com/dchest/uniuri"
	"strconv"
	"time"
)

type Identifier struct {
	id string
}

func (i Identifier) Id() string {
	if i.id == "" {
		i.id = generateRandomId()
	}
	return i.id
}

// generateRandomId generates a random string that is more or less
// garunteed to be unique. Used as ids for records where an id is
// not otherwise provided.
func generateRandomId() string {
	timeInt := time.Now().Unix()
	timeString := strconv.FormatInt(timeInt, 36)
	randomString := uniuri.NewLen(16)
	return randomString + timeString
}
