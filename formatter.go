package ayolog

import "time"

const (
	defaultTimeStampFormat = time.RFC3339
	FieldKeyMsg            = "msg"
	FieldKeyLevel          = "level"
	FieldKeyTime           = "time"
	FieldKeyFile           = "file"
	FieldKeyFunc           = "func"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}
