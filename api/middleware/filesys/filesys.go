package filesys

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/utils"
	"html"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// NewFileSys is a slightly stripped version of Fiber's
// original filesystem.New function.
// The function will return a given file for any route.
func NewFileSys(config ...filesystem.Config) fiber.Handler {
	cfg := filesystem.ConfigDefault
	if len(config) > 0 {
		cfg = config[0]
		if cfg.Index == "" {
			cfg.Index = filesystem.ConfigDefault.Index
		}
		if !strings.HasPrefix(cfg.Index, "/") {
			cfg.Index = "/" + cfg.Index
		}
		if cfg.NotFoundFile != "" && !strings.HasPrefix(cfg.NotFoundFile, "/") {
			cfg.NotFoundFile = "/" + cfg.NotFoundFile
		}
	}

	if cfg.Root == nil {
		panic("filesystem: Root cannot be nil")
	}

	if cfg.PathPrefix != "" && !strings.HasPrefix(cfg.PathPrefix, "/") {
		cfg.PathPrefix = "/" + cfg.PathPrefix
	}
	cacheControlStr := "public, max-age=" + strconv.Itoa(cfg.MaxAge)

	return func(c *fiber.Ctx) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		method := c.Method()

		// We only serve static assets on GET or HEAD methods
		if method != fiber.MethodGet && method != fiber.MethodHead {
			return c.Next()
		}
		path := ""
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		// Add PathPrefix
		if cfg.PathPrefix != "" {
			// PathPrefix already has a "/" prefix
			path = cfg.PathPrefix + path
		}

		var (
			file http.File
			stat os.FileInfo
		)

		if len(path) > 1 {
			path = utils.TrimRight(path, '/')
		}
		file, err = cfg.Root.Open(path)
		if err != nil && os.IsNotExist(err) && cfg.NotFoundFile != "" {
			file, err = cfg.Root.Open(cfg.NotFoundFile)
		}

		if err != nil {
			if os.IsNotExist(err) {
				return c.Status(fiber.StatusNotFound).Next()
			}
			return
		}

		if stat, err = file.Stat(); err != nil {
			return
		}

		// Serve index if path is directory
		if stat.IsDir() {
			indexPath := utils.TrimRight(path, '/') + cfg.Index
			index, err := cfg.Root.Open(indexPath)
			if err == nil {
				indexStat, err := index.Stat()
				if err == nil {
					file = index
					stat = indexStat
				}
			}
		}

		// Browse directory if no index found and browsing is enabled
		if stat.IsDir() {
			if cfg.Browse {
				return dirList(c, file)
			}
			return fiber.ErrForbidden
		}

		modTime := stat.ModTime()
		contentLength := int(stat.Size())

		// Set Content Type header
		c.Type(getFileExtension(stat.Name()))

		// Set Last Modified header
		if !modTime.IsZero() {
			c.Set(fiber.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
		}

		if method == fiber.MethodGet {
			if cfg.MaxAge > 0 {
				c.Set(fiber.HeaderCacheControl, cacheControlStr)
			}
			c.Response().SetBodyStream(file, contentLength)
			return nil
		}
		if method == fiber.MethodHead {
			c.Request().ResetBody()
			c.Response().SkipBody = true
			c.Response().Header.SetContentLength(contentLength)
			if err := file.Close(); err != nil {
				return err
			}
			return nil
		}

		return c.Next()
	}
}

func getFileExtension(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n < 0 {
		return ""
	}
	return path[n:]
}

func dirList(c *fiber.Ctx, f http.File) error {
	fileinfos, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	fm := make(map[string]os.FileInfo, len(fileinfos))
	filenames := make([]string, 0, len(fileinfos))
	for _, fi := range fileinfos {
		name := fi.Name()
		fm[name] = fi
		filenames = append(filenames, name)
	}

	basePathEscaped := html.EscapeString(c.Path())
	fmt.Fprintf(c, "<html><head><title>%s</title><style>.dir { font-weight: bold }</style></head><body>", basePathEscaped)
	fmt.Fprintf(c, "<h1>%s</h1>", basePathEscaped)
	fmt.Fprint(c, "<ul>")

	if len(basePathEscaped) > 1 {
		parentPathEscaped := html.EscapeString(utils.TrimRight(c.Path(), '/') + "/..")
		fmt.Fprintf(c, `<li><a href="%s" class="dir">..</a></li>`, parentPathEscaped)
	}

	sort.Strings(filenames)
	for _, name := range filenames {
		pathEscaped := html.EscapeString(path.Join(c.Path() + "/" + name))
		fi := fm[name]
		auxStr := "dir"
		className := "dir"
		if !fi.IsDir() {
			auxStr = fmt.Sprintf("file, %d bytes", fi.Size())
			className = "file"
		}
		fmt.Fprintf(c, `<li><a href="%s" class="%s">%s</a>, %s, last modified %s</li>`,
			pathEscaped, className, html.EscapeString(name), auxStr, fi.ModTime())
	}
	fmt.Fprint(c, "</ul></body></html>")

	c.Type("html")

	return nil
}
