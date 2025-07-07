package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

// ServiceRegistry holds the URLs for the different microservices.
type ServiceRegistry struct {
	services map[string]*url.URL
	logger   *zap.Logger
}

// NewServiceRegistry initializes a new ServiceRegistry with predefined services.
// Each service is identified by a unique name and has a corresponding URL.
func NewServiceRegistry(logger *zap.Logger) *ServiceRegistry {
	services := map[string]*url.URL{
		"leaderboard": mustParseURL("http://leaderboard-service:3001"),
		// New services go under here [Pun intended :)]
	}
	return &ServiceRegistry{
		services: services,
		logger:   logger,
	}
}

// proxyHandler creates a Fiber handler that proxies requests to the specified service.
func (sr *ServiceRegistry) proxyHandler(serviceName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sr.logger.Info("Proxy handler invoked",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("serviceName", serviceName),
		)

		serviceURL, exists := sr.services[serviceName]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ // Use Fiber's status constants
				"error": "Service not available",
			})
		}

		// Construct the target URL for the internal service.
		path := c.Path()
		cleanPath := strings.TrimPrefix(path, "/api/v1/"+serviceName)

		targetURL := serviceURL.String() + cleanPath
		if len(c.Request().URI().QueryString()) > 0 {
			targetURL += "?" + string(c.Request().URI().QueryString())
		}

		sr.logger.Info("Path details for proxying",
			zap.String("originalPath", path),
			zap.String("prefixUsed", "/api/v1/"+serviceName),
			zap.String("cleanedPath", cleanPath),
			zap.String("targetURL", targetURL),
		)
	
		req, err := http.NewRequest(c.Method(), targetURL, strings.NewReader(string(c.Body())))
		if err != nil {
			sr.logger.Error("Failed to create proxy request", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Proxy error"}) // Use Fiber's status constants
		}


		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			sr.logger.Error("Service request failed",
				zap.String("service", serviceName),
				zap.Error(err),
			)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{ // Use Fiber's status constants
				"error": "Service unavailable",
			})
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			sr.logger.Error("Failed to read response body from proxied service",
				zap.String("service", serviceName),
				zap.Error(err),
			)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read response from service"})
		}

		c.Status(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		return c.Send(bodyBytes)
	}
}

// mustParseURL is a helper function to parse URLs and panic on error.
func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse URL %s: %v", rawURL, err))
	}
	return u
}

func main() {
	// Initialize a production-ready Zap logger.
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err) 
	}
	defer logger.Sync()

	logger.Info("Gateway service starting...") 

	// Initialize the service registry.
	registry := NewServiceRegistry(logger)

	// Create a new Fiber app.
	app := fiber.New()

	// Global middleware for CORS.
	app.Use(cors.New())

	app.Use(func(c *fiber.Ctx) error {
		logger.Info("Request received by Fiber app",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
		)
		return c.Next()
	})

	// Health check for the gateway itself.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "gateway"})
	})

	// Service routes: All requests to /api/v1/leaderboard/* will be proxied
	// to the 'leaderboard' service defined in the registry.
	app.All("/api/v1/leaderboard/*", registry.proxyHandler("leaderboard"))

	// Start the Fiber application and listen on port 80.
	// log.Fatal will log the error and exit if the server fails to start.
	log.Fatal(app.Listen(":80"))
}
