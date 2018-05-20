package prepaid

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

func UUID() string {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	return id.String()
}
