package productErrors

import "github.com/pkg/errors"

var (
	ErrObjectIDTypeConversion = errors.New("object id type conversion")
)
