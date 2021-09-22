package devctlpath

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	_ "github.com/onsi/gomega/matchers"

	"github.com/alex-held/devctl-kit/pkg/system"
)

func TestHomeFinder(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	var lazyFinder Pather

	customUserHomeFn := func() string {
		return "/h/o/m/e/user"
	}

	customCacheFn := func() string {
		return "/c/a/c/h/e/"
	}

	const testAppPrefix = "test_devctl"
	var testAppPrefixWithLeadingDot = fmt.Sprintf(".%s", testAppPrefix)
	var customConfigRoot = fmt.Sprintf("/h/o/m/e/user/%s", testAppPrefix)

	/* ConfigRoot */
	g.Describe("ConfigRoot", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix)
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefixWithLeadingDot)
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := customConfigRoot
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := customConfigRoot
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.ConfigRoot()
			Expect(actual).To(Equal(expected))
		})
	})

	/* CacheDir */
	g.Describe("Cache", func() {
		g.It("WHEN no pathFn set", func() {
			if runtime.GOOS == "linux" {
				_ = os.Setenv("XDG_CACHE_HOME", "/tmp/cache")
			}

			cacheDir, _ := os.UserCacheDir()
			expected := filepath.Join(cacheDir, "io.alexheld.test_devctl")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN cacheFn set", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithCachePathFn(customCacheFn))
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing path elems", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithCachePathFn(customCacheFn), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.Cache("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customCacheFn(), "io.alexheld.test_devctl")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithCachePathFn(customCacheFn), WithConfigRootFn(func() string {
				return expected
			}))
			actual := lazyFinder.Cache()
			Expect(actual).To(Equal(expected))
		})
	})

	/* Bin */
	g.Describe("Bin", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "bin")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Bin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "/bin/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Bin("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	/* SDK */
	g.Describe("SDK", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "sdks")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.SDK()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "sdks/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.SDK("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	/* PLUGINS */
	g.Describe("PLUGINS", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "plugins")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Plugin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "plugins")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.Plugin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "plugins")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Plugin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "plugins")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Plugin()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "plugins/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Plugin("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	/* Config */
	g.Describe("Config", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "config")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Config()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "config/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Config("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})

	/* Download */
	g.Describe("Download", func() {
		g.It("WHEN no pathFn set", func() {
			userHome, _ := os.UserHomeDir()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN userHomeFn set", func() {
			userHome := customUserHomeFn()
			expected := resolveConfigSubDir(userHome, testAppPrefix, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithUserHomeFn(customUserHomeFn))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN configRoot set", func() {
			expected := filepath.Join(customConfigRoot, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN app prefix starts with a '.'", func() {
			expected := filepath.Join(customConfigRoot, "downloads")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefixWithLeadingDot), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Download()
			Expect(actual).To(Equal(expected))
		})

		g.It("WHEN providing sub directories parameter", func() {
			expected := filepath.Join(customConfigRoot, "downloads/some/sub/dir")
			lazyFinder = NewPather(WithAppPrefix(testAppPrefix), WithConfigRootFn(func() string {
				return customConfigRoot
			}))
			actual := lazyFinder.Download("some", "sub", "dir")
			Expect(actual).To(Equal(expected))
		})
	})
}

type ExpectedPaths struct {
	user    string
	prefix  string
	cache   string
	cfgFile string
}

func (e ExpectedPaths) ConfigRoot(elem ...string) string {
	return filepath.Join(filepath.Join(e.user, e.prefix), filepath.Join(elem...))
}

func (e ExpectedPaths) Config(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "config", filepath.Join(elem...))
}

func (e ExpectedPaths) Bin(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "bin", filepath.Join(elem...))
}

func (e ExpectedPaths) Download(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "downloads", filepath.Join(elem...))
}

func (e ExpectedPaths) SDK(elem ...string) string {
	return filepath.Join(e.ConfigRoot(), "sdks", filepath.Join(elem...))
}

func (e ExpectedPaths) Cache(elem ...string) string {
	return filepath.Join(e.cache, fmt.Sprintf("io.alexheld%s", e.prefix), filepath.Join(elem...))
}

func (e ExpectedPaths) ConfigFilePath() string {
	return e.ConfigRoot(e.cfgFile)
}

func resolveConfigSubDir(home, prefix string, elem ...string) (path string) {

	prefix = "." + strings.TrimPrefix(prefix, ".")
	var cfgRoot string
	switch {
	case system.OSRuntimeInfoGetter{}.Get().OS == "linux":
		cfgRoot = filepath.Join(home, ".config", prefix)
	case system.OSRuntimeInfoGetter{}.Get().OS == "darwin":
		cfgRoot = filepath.Join(home, prefix)
	default:
		// windows not yet supported
		panic(fmt.Errorf("the current os is not yet supported; os=%s", system.OSRuntimeInfoGetter{}.Get().OS))
	}
	return filepath.Join(cfgRoot, filepath.Join(elem...))
}
