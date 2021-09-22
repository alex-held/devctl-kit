package log_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alex-held/devctl-kit/pkg/log"
)

func TestLevel_Colorize(t *testing.T) {
	tcs := []struct {
		level     log.Level
		colorCode string
		expected  string
	}{
		{
			log.Debug,
			"1;32",
			log.DebugPrefix,
		},
		{
			log.Info,
			"37",
			log.InfoPrefix,
		},
		{
			log.Warn,
			"1;33",
			log.WarnPrefix,
		},
		{
			log.Error,
			"31",
			log.ErrorPrefix,
		},
		{
			log.FATAL,
			"35",
			log.FatalPrefix,
		},
	}

	for _, tt := range tcs {
		t.Run(tt.level.String(), func(t *testing.T) {
			actual := tt.level.Colorize()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func setup(setupFn func(config *log.Config)) (l log.Logger, out *bytes.Buffer) {
	out = &bytes.Buffer{}
	cfg := &log.Config{
		Color: false,
		FatalFunc: func() {
			out.Write([]byte(fmt.Sprint("fatalFunc()")))
		},
	}
	if setupFn != nil {
		setupFn(cfg)
	}
	if cfg.Out != nil {
		cfg.Out = io.MultiWriter(cfg.Out, out)
	} else {
		cfg.Out = out
	}

	return log.New(cfg), out
}

func TestLogger_FatalFunc(t *testing.T) {
	fatalFuncCalled := false
	logger, _ := setup(func(config *log.Config) {
		config.FatalFunc = func() {
			fatalFuncCalled = true
		}
	})

	logger.Fatalf("bye")

	assert.True(t, fatalFuncCalled)
}

func TestLogger_Debugf(t *testing.T) {
	logger, out := setup(nil)
	logger.Debugf("hello %s", "world")

	assert.Equal(t, "[DEBUG]	  hello world\n", out.String())
}

func TestLogger_Infof(t *testing.T) {
	logger, out := setup(nil)
	logger.Infof("hello %s", "world")

	assert.Equal(t, "[INFO]	  hello world\n", out.String())
}

func TestLogger_Warnf(t *testing.T) {
	logger, out := setup(nil)
	logger.Warnf("hello %s", "world")

	assert.Equal(t, "[WARN]	  hello world\n", out.String())
}

func TestLogger_Errorf(t *testing.T) {
	logger, out := setup(nil)
	logger.Errorf("hello %s", "world")

	assert.Equal(t, "[ERROR]	  hello world\n", out.String())
}

func TestLogger_Fatalf(t *testing.T) {
	logger, out := setup(func(config *log.Config) {
		config.FatalFunc = func() {}
	})
	logger.Fatalf("hello %s", "world")
	assert.Equal(t, "[FATAL]	  hello world\n", out.String())
}

func TestLogger_WithColor2(t *testing.T) {
	tcs := []struct {
		level     log.Level
		colorCode string
		expected  string
	}{
		{
			log.Debug,
			"1;32",
			"\033[1;30m[\033[1;32mDEBUG\033[1;30m]\033[0m\t  hello world!\n",
		},
		{
			log.Info,
			"1;37",
			"\033[1;30m[\033[1;37mINFO\033[1;30m]\033[0m\t  hello world!\n",
		},
		{
			log.Warn,
			"1;33",
			"\033[1;30m[\033[1;33mWARN\033[1;30m]\033[0m\t  hello world!\n",
		},
		{
			log.Error,
			"1;31",
			"\033[1;30m[\033[1;31mERROR\033[1;30m]\033[0m\t  hello world!\n",
		},
		{
			log.FATAL,
			"1;35",
			"\033[1;30m[\033[1;35mFATAL\033[1;30m]\033[0m\t  hello world!\n",
		},
	}

	for _, tt := range tcs {
		t.Run(tt.level.String(), func(t *testing.T) {
			logger, out := setup(func(config *log.Config) {
				config.Color = true
				config.FatalFunc = func() {}
			})

			logger.Logf(tt.level, "hello %s!", "world")
			actual := out.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
