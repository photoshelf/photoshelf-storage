package photo

import (
	"crypto/md5"
	"fmt"
	"time"
)

type Identifier struct {
	value string
}

func NewIdentifier(data []byte) *Identifier {
	dataHash := fmt.Sprintf("%x", md5.Sum(data))
	filename := fmt.Sprintf("%x", md5.Sum([]byte(dataHash+time.Now().String())))
	return &Identifier{filename}
}

func IdentifierOf(value string) *Identifier {
	return &Identifier{value}
}

func (id *Identifier) Value() string {
	return id.value
}
