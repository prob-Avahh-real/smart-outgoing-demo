package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// TLSCertConfig holds TLS certificate configuration
type TLSCertConfig struct {
	CertFile string
	KeyFile  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	HTTPPort  string
	HTTPSPort string
	TLS       *TLSCertConfig
	EnableTLS bool
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		HTTPPort:  "8080",
		HTTPSPort: "8443",
		EnableTLS: false,
		TLS: &TLSCertConfig{
			CertFile: "./certs/server.crt",
			KeyFile:  "./certs/server.key",
		},
	}
}

// ServerConfigFromEnv loads server configuration from environment variables
func ServerConfigFromEnv() *ServerConfig {
	config := DefaultServerConfig()

	if port := os.Getenv("HTTP_PORT"); port != "" {
		config.HTTPPort = port
	}
	if port := os.Getenv("HTTPS_PORT"); port != "" {
		config.HTTPSPort = port
	}
	if cert := os.Getenv("TLS_CERT_FILE"); cert != "" {
		config.TLS.CertFile = cert
	}
	if key := os.Getenv("TLS_KEY_FILE"); key != "" {
		config.TLS.KeyFile = key
	}
	if enableTLS := os.Getenv("ENABLE_TLS"); enableTLS == "true" || enableTLS == "1" {
		config.EnableTLS = true
	}

	return config
}

// RunServer runs the HTTP/HTTPS server
func RunServer(router *gin.Engine, config *ServerConfig) error {
	// Create HTTP to HTTPS redirect server
	if config.EnableTLS {
		go func() {
			redirectServer := &http.Server{
				Addr: ":" + config.HTTPPort,
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Extract host without port
					host := r.Host
					if idx := strings.Index(host, ":"); idx != -1 {
						host = host[:idx]
					}
					target := "https://" + host + ":" + config.HTTPSPort + r.URL.Path
					if len(r.URL.RawQuery) > 0 {
						target += "?" + r.URL.RawQuery
					}
					log.Printf("Redirecting HTTP request to HTTPS: %s -> %s", r.URL.Path, target)
					http.Redirect(w, r, target, http.StatusPermanentRedirect)
				}),
			}
			log.Printf("HTTP redirect server listening on port %s", config.HTTPPort)
			if err := redirectServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP redirect server error: %v", err)
			}
		}()
	}

	// Configure TLS
	if config.EnableTLS {
		// Verify certificate files exist
		if _, err := os.Stat(config.TLS.CertFile); os.IsNotExist(err) {
			log.Printf("TLS certificate file not found: %s", config.TLS.CertFile)
			log.Printf("Please generate certificates using: ./scripts/generate-certs.sh")
			return err
		}
		if _, err := os.Stat(config.TLS.KeyFile); os.IsNotExist(err) {
			log.Printf("TLS key file not found: %s", config.TLS.KeyFile)
			log.Printf("Please generate certificates using: ./scripts/generate-certs.sh")
			return err
		}

		// Configure TLS with modern security settings and HTTP/2 support
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // Required for HTTP/2
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256, // Required for HTTP/2
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		server := &http.Server{
			Addr:      ":" + config.HTTPSPort,
			Handler:   router,
			TLSConfig: tlsConfig,
		}

		log.Printf("HTTPS server listening on port %s", config.HTTPSPort)
		log.Printf("TLS Certificate: %s", config.TLS.CertFile)
		log.Printf("TLS Key: %s", config.TLS.KeyFile)

		return server.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile)
	}

	// Run HTTP server only
	log.Printf("HTTP server listening on port %s", config.HTTPPort)
	return router.Run("0.0.0.0:" + config.HTTPPort)
}

// RunDualServer runs both HTTP and HTTPS servers
func RunDualServer(router *gin.Engine, config *ServerConfig) error {
	if config.EnableTLS {
		// Run HTTPS server
		go func() {
			if err := RunServer(router, config); err != nil {
				log.Printf("HTTPS server error: %v", err)
			}
		}()

		// Run HTTP server on different port for development
		httpRouter := gin.New()
		httpRouter.Use(func(c *gin.Context) {
			// Copy the main router handlers
			router.HandleContext(c)
		})

		log.Printf("HTTP server also listening on port %s (for development)", config.HTTPPort)
		return httpRouter.Run("0.0.0.0:" + config.HTTPPort)
	}

	return RunServer(router, config)
}
