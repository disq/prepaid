package prepaid

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

// UUID generates a ULID using crypto/rand as the source.
func UUID() string {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	return id.String()
}
