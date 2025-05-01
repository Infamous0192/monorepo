package fiber

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

type Server struct {
	*fiber.App
	Config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	server := &Server{
		Config: config,
	}
	server.setupPaths()

	server.App = fiber.New(fiber.Config{
		Concurrency:  256,
		ServerHeader: server.Config.Name,
		BodyLimit:    server.Config.UploadLimit,
		ErrorHandler: server.Config.ErrorHandler,
	})

	server.setupMiddlewares()
	server.setupStatic()

	if server.Config.SwaggerConfig != nil {
		server.setupSwagger()
	}

	return server
}

func (s *Server) setupMiddlewares() {
	s.Use(recover.New())
	s.Use(logger.New())

	if s.Config.Cors != nil {
		s.Use(cors.New(cors.Config{
			AllowOrigins:     s.Config.Cors.AllowOrigins,
			AllowMethods:     s.Config.Cors.AllowMethods,
			AllowHeaders:     s.Config.Cors.AllowHeaders,
			AllowCredentials: s.Config.Cors.AllowCredentials,
			ExposeHeaders:    s.Config.Cors.ExposeHeaders,
			MaxAge:           s.Config.Cors.MaxAge,
		}))
	}
}

func (s *Server) setupStatic() {
	// Upload Path
	if s.Config.UploadPath != "" {
		s.Static("/uploads", s.Config.UploadPath, fiber.Static{
			Compress:      true,
			ByteRange:     true,
			CacheDuration: 24 * time.Hour,
		})
	}

	// Public Path
	if s.Config.PublicPath != "" {
		s.Static("/", s.Config.PublicPath)

		// Redirect when access non-index route
		s.Use(func(c *fiber.Ctx) error {
			// Skip if paths starting with "/api" or "/docs"
			if strings.HasPrefix(c.Path(), "/api") || strings.HasPrefix(c.Path(), "/docs") {
				return c.Next()
			}

			// Serve index.html for all other non-file routes
			if _, err := os.Stat(filepath.Join(s.Config.PublicPath, c.Path())); os.IsNotExist(err) {
				return c.SendFile(filepath.Join(s.Config.PublicPath, "index.html"))
			}
			return c.Next()
		})
	}
}

func (s *Server) setupSwagger() {
	s.Get("/docs/*", swagger.New(*s.Config.SwaggerConfig))
}

func (s *Server) setupPaths() {
	if s.Config.Url == "" {
		s.Config.Url = fmt.Sprintf("http://%s:%s", s.Config.Host, s.Config.Port)
	}

	path, _ := os.Getwd()
	if s.Config.ExecPath {
		path = getExecutablePath()
	}

	s.Config.Path = path
	if s.Config.UploadPath != "" {
		s.Config.UploadPath = makeDir(filepath.Join(path, s.Config.UploadPath))
	}
	if s.Config.PublicPath != "" {
		s.Config.PublicPath = makeDir(filepath.Join(path, s.Config.PublicPath))
	}
	if s.Config.UploadLimit > 0 {
		s.Config.UploadLimit = s.Config.UploadLimit * 1024 * 1024 // Convert to bytes
	}
}

func (s *Server) Serve() error {
	return s.Listen(fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port))
}

func getExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func makeDir(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	return path
}
