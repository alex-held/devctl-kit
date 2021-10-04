package env

import (
	"io"
	"os"
	"sync"

	"github.com/spf13/afero"

	"github.com/alex-held/devctl-kit/pkg/log"
	"github.com/alex-held/devctl-kit/pkg/system"

	"github.com/alex-held/devctl-kit/pkg/cli"
	"github.com/alex-held/devctl-kit/pkg/cli/util"
	"github.com/alex-held/devctl-kit/pkg/validation"
)

type factory struct {
	runtimeInfoGetter system.RuntimeInfoGetter
	logger            log.Logger

	// Caches OpenAPI document and parsed resources
	//	openAPIParser *openapi.CachedOpenAPIParser
	//	openAPIGetter *openapi.CachedOpenAPIGetter
	//	parser        sync.Once
	getter  sync.Once
	streams cli.IOStreams
	fs      afero.Fs
	paths   Paths
}

func (f *factory) RuntimeInfo() system.RuntimeInfo {
	return f.runtimeInfoGetter.Get()
}

func (f *factory) NewBuilder() *util.Builder {
	return &util.Builder{}
}

func (f *factory) Fs() afero.Fs {
	return f.fs
}

func (f *factory) Logger() log.Logger {
	return f.logger
}
func (f *factory) Paths() Paths {
	return f.paths
}

func (f *factory) Streams() cli.IOStreams {
	return f.streams
}

func (f *factory) Validator(validate bool) (validation.Schema, error) {
	return validation.NullSchema{}, nil
}

// Factory provides abstractions that allow the Devctl command to be extended across multiple types
// of resources and different API sets.
type Factory interface {
	RuntimeInfo() system.RuntimeInfo

	// NewBuilder returns an object that assists in loading objects from both disk and the server
	// and which implements the common patterns for CLI interactions with generic resources.
	NewBuilder() *util.Builder

	Logger() log.Logger

	Paths() Paths

	Fs() afero.Fs

	Streams() cli.IOStreams

	// Returns a schema that can validate objects stored on disk.
	Validator(validate bool) (validation.Schema, error)

	// OpenAPISchema returns the parsed openapi schema definition
	//	OpenAPISchema() (openapi.Resources, error)
	// OpenAPIGetter returns a getter for the openapi schema document
	//	OpenAPIGetter() discovery.OpenAPISchemaInterface
}

type FactoryConfig struct {
	Paths             Paths
	LoggerConfig      *log.Config
	Streams           *cli.IOStreams
	RuntimeInfoGetter system.RuntimeInfoGetter
	Fs                afero.Fs
}

type FactoryOption func(*FactoryConfig) *FactoryConfig

func WithIO(in io.Reader, out, err io.Writer) FactoryOption {
	return func(c *FactoryConfig) *FactoryConfig {
		c.Streams = &cli.IOStreams{
			In:     in,
			Out:    out,
			ErrOut: err,
		}
		return c
	}
}

func NewFactory(opts ...FactoryOption) Factory {
	cfg := &FactoryConfig{
		LoggerConfig:      &log.DefaultConfig,
		Fs:                afero.NewOsFs(),
		RuntimeInfoGetter: system.OSRuntimeInfoGetter{},
		Paths:             MustGetPaths(),
	}
	defaults := []FactoryOption{
		WithIO(os.Stdin, os.Stdout, os.Stdout),
	}

	for _, opt := range append(defaults, opts...) {
		opt(cfg)
	}

	return &factory{
		runtimeInfoGetter: cfg.RuntimeInfoGetter,
		paths:             cfg.Paths,
		logger:            log.New(cfg.LoggerConfig),
		getter:            sync.Once{},
		streams:           *cfg.Streams,
		fs:                cfg.Fs,
	}
}
