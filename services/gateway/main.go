package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type ServiceRegistry struct {
	services map[string]*url.URL
	logger   *zap.Logger
}
// NewServiceRegistry initializes a new ServiceRegistry with predefined services.
// Each service is identified by a unique name and has a corresponding URL.
func NewServiceRegistry(logger *zap.Logger) *ServiceRegistry {
	services := map[string]*url.URL{
		"leaderboard": mustParseURL("http://localhost:3001"),
		// New services go under here [Pun intended :)]

	}
	return &ServiceRegistry{
		services: services,
		logger:   logger,
	}
}

func (sr *ServiceRegistry) proxyHandler(serviceName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		serviceURL, exists := sr.services[serviceName]
		if !exists {
			return c.Status(404).JSON(fiber.Map{
				"error": "Service not available",
			})
		}

		path := c.Path()
		cleanPath := strings.TrimPrefix(path, "/api/v1"+serviceName)

		targetURL := serviceURL.String() + cleanPath

		req, err := http.NewRequest(c.Method(), targetURL, strings.NewReader(string(c.Body())))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Proxy error"})
		}
		// Copy all headers from the incoming request and
		// add to that to the one to be forwarded
		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		// Take the parameters if there are any from the incoming req
		req.URL.RawQuery = string(c.Request().URI().QueryString())

		// Execute the http request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			sr.logger.Error("Service request failed",
				zap.String("service", serviceName),
				zap.Error(err),
			)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Service unavailable",
			})
		}
		defer resp.Body.Close()

		c.Status(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		return c.SendStream(resp.Body)
	}
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	registry := NewServiceRegistry(logger)

	app := fiber.New()

	// Global middleware
	app.Use(cors.New())

	// Health check for the gateway 	
	app.Get("/health", func (c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "gateway"})
	})


	// Service routes 
	app.All("/api/v1/leaderboard/*", registry.proxyHandler("leaderboard"))

	log.Fatal(app.Listen(":80"))
}
