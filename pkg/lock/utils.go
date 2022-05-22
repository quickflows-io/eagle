package lock

import (
	"github.com/google/uuid"
)

// genToken generate token
func genToken() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
