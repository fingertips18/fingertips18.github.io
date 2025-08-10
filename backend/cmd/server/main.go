package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fingertips18/fingertips18.github.io/backend/internal/server"
	flagUtils "github.com/Fingertips18/fingertips18.github.io/backend/pkg/utils"
	"github.com/joho/godotenv"
)

// Define constants for flags to improve manageability
const (
	FlagEnv                 = "env"
	FlagPort                = "port"
	FlagEmailJSServiceID    = "emailjs-service-id"
	FlagEmailJSTemplateID   = "emailjs-template-id"
	FlagEmailJSPublicKey    = "emailjs-public-key"
	FlagEmailJSPrivateKey   = "emailjs-private-key"
	FlagGoogleMeasurementID = "google-measurement-id"
)

func main() {
	var (
		flagEnvironment         = flag.String(FlagEnv, "local", "Environment")
		flagPort                = flag.String(FlagPort, "8080", "Port server")
		flagEmailJSServiceID    = flag.String(FlagEmailJSServiceID, "", "EmailJS Service ID")
		flagEmailJSTemplateID   = flag.String(FlagEmailJSTemplateID, "", "EmailJS Template ID")
		flagEmailJSPublicKey    = flag.String(FlagEmailJSPublicKey, "", "EmailJS Public Key")
		flagEmailJSPrivateKey   = flag.String(FlagEmailJSPrivateKey, "", "EmailJS Private Key")
		flagGoogleMeasurementID = flag.String(FlagGoogleMeasurementID, "", "Google Measurement ID")
	)

	flag.Parse()

	flagUtils.Require(
		FlagEmailJSServiceID,
		FlagEmailJSTemplateID,
		FlagEmailJSPublicKey,
		FlagGoogleMeasurementID,
	)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get secret token values
	port := *flagPort
	emailJSServiceID := *flagEmailJSServiceID
	emailJSTemplateID := *flagEmailJSTemplateID
	emailJSPublicKey := *flagEmailJSPublicKey
	emailJSPrivateKey := *flagEmailJSPrivateKey
	googleMeasurementID := *flagGoogleMeasurementID
	if *flagEnvironment != "local" {
		data, err := os.ReadFile(*flagPort)
		if err != nil {
			log.Printf("Failed to read port from file, using flag value: %v", *flagPort)
		} else {
			port = string(data)
		}

		data, err = os.ReadFile(*flagEmailJSServiceID)
		if err != nil {
			log.Printf("Failed to read emailjs service ID from file, using flag value: %v", *flagEmailJSServiceID)
		} else {
			emailJSServiceID = string(data)
		}

		data, err = os.ReadFile(*flagEmailJSTemplateID)
		if err != nil {
			log.Printf("Failed to read emailjs template ID from file, using flag value: %v", *flagEmailJSTemplateID)
		} else {
			emailJSTemplateID = string(data)
		}

		data, err = os.ReadFile(*flagEmailJSPublicKey)
		if err != nil {
			log.Printf("Failed to read emailjs public key from file, using flag value: %v", *flagEmailJSPublicKey)
		} else {
			emailJSPublicKey = string(data)
		}

		data, err = os.ReadFile(*flagEmailJSPrivateKey)
		if err != nil {
			log.Printf("Failed to read emailjs private key from file, using flag value: %v", *flagEmailJSPrivateKey)
		} else {
			emailJSPrivateKey = string(data)
		}

		data, err = os.ReadFile(*flagGoogleMeasurementID)
		if err != nil {
			log.Printf("Failed to read google measurement ID from file, using flag value: %v", *flagGoogleMeasurementID)
		} else {
			googleMeasurementID = string(data)
		}
	} else {
		if port == "" {
			port = os.Getenv("PORT")
		}

		if emailJSServiceID == "" {
			emailJSServiceID = os.Getenv("EMAILJS_SERVICE_ID")
		}

		if emailJSTemplateID == "" {
			emailJSTemplateID = os.Getenv("EMAILJS_TEMPLATE_ID")
		}

		if emailJSPublicKey == "" {
			emailJSPublicKey = os.Getenv("EMAILJS_PUBLIC_KEY")
		}

		if emailJSPrivateKey == "" {
			emailJSPrivateKey = os.Getenv("EMAILJS_PRIVATE_KEY")
		}

		if googleMeasurementID == "" {
			googleMeasurementID = os.Getenv("GOOGLE_MEASUREMENT_ID")
		}
	}

	// Setup server
	s := server.New(server.Config{
		Environment:         *flagEnvironment,
		Port:                port,
		EmailJSServiceID:    emailJSServiceID,
		EmailJSTemplateID:   emailJSTemplateID,
		EmailJSPublicKey:    emailJSPublicKey,
		EmailJSPrivateKey:   emailJSPrivateKey,
		GoogleMeasurementID: googleMeasurementID,
	},
	)

	// Initialize the server in a goroutine so that it won't block the graceful shutdown handling
	go func() {
		if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Listen for the interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting...")
}
