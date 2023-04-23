package files

import "io"

type FilterFunc func(r io.Reader) bool
