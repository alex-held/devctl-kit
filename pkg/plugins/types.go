package plugins

import (
	"context"
	"io"
)

type Context struct {
	Out     io.Writer
	Pather  string
	Context context.Context
}
