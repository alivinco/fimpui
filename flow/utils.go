package flow

import (
	"github.com/dchest/uniuri"
)

func GenerateId(len int) string {
	return uniuri.NewLen(len)
}
