package plugins

import (
	"context"
	"io"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
)

type Context struct {
	Out     io.Writer
	Pather  devctlpath.Pather
	Context context.Context
}
