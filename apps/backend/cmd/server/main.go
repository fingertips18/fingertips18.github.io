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

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/server"
	flagUtils "github.com/fingertips18/fingertips18.github.io/backend/pkg/utils"
	"github.com/joho/godotenv"
)

// Define constants for flags to improve manageability
const (
	FlagEnv                  = "env"
	FlagClientURL            = "client-url"
	FlagPort                 = "port"
	FlagAuthToken            = "auth-token"
	FlagEmailJSServiceID     = "emailjs-service-id"
	FlagEmailJSTemplateID    = "emailjs-template-id"
	FlagEmailJSPublicKey     = "emailjs-public-key"
	FlagEmailJSPrivateKey    = "emailjs-private-key"
	FlagGoogleMeasurementID  = "google-measurement-id"
	FlagGoogleAPISecret      = "google-api-secret" // #nosec
	FlagDatabaseURL          = "database-url"
	FlagUsername             = "username"
	FlagPassword             = "password"
	FlagUploadthingSecretKey = "uploadthing-secret-key"
)

// @title Portfolio Backend API
// @version 1.0
// @description This API powers the portfolio platform, offering endpoints to manage projects, users, and supporting data for both public and admin use.
//
// @license.name MIT
// @license.url https://opensource.org/license/mit
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	var (
		flagEnvironment          = flag.String(FlagEnv, "local", "Environment")
		flagClientURL            = flag.String(FlagClientURL, "http://localhost:5378", "Client URL")
		flagPort                 = flag.String(FlagPort, "8080", "Port server")
		flagAuthToken            = flag.String(FlagAuthToken, "", "Basic token auth")
		flagEmailJSServiceID     = flag.String(FlagEmailJSServiceID, "", "EmailJS Service ID")
		flagEmailJSTemplateID    = flag.String(FlagEmailJSTemplateID, "", "EmailJS Template ID")
		flagEmailJSPublicKey     = flag.String(FlagEmailJSPublicKey, "", "EmailJS Public Key")
		flagEmailJSPrivateKey    = flag.String(FlagEmailJSPrivateKey, "", "EmailJS Private Key")
		flagGoogleMeasurementID  = flag.String(FlagGoogleMeasurementID, "", "Google Measurement ID")
		flagGoogleAPISecret      = flag.String(FlagGoogleAPISecret, "", "Google API Secret")
		flagDatabaseURL          = flag.String(FlagDatabaseURL, "", "Postgres Database URL")
		flagUsername             = flag.String(FlagUsername, "", "Backend Username Access")
		flagPassword             = flag.String(FlagPassword, "", "Backend Password Access")
		flagUploadthingSecretKey = flag.String(FlagUploadthingSecretKey, "", "Uploadthing Secret Key")
	)

	flag.Parse()

	flagUtils.Require(
		FlagAuthToken,
		FlagEmailJSServiceID,
		FlagEmailJSTemplateID,
		FlagEmailJSPublicKey,
		FlagGoogleMeasurementID,
		FlagGoogleAPISecret,
		FlagDatabaseURL,
		FlagUsername,
		FlagPassword,
		FlagUploadthingSecretKey,
	)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get secret token values
	port := *flagPort
	clientURL := *flagClientURL
	authToken := *flagAuthToken
	emailJSServiceID := *flagEmailJSServiceID
	emailJSTemplateID := *flagEmailJSTemplateID
	emailJSPublicKey := *flagEmailJSPublicKey
	emailJSPrivateKey := *flagEmailJSPrivateKey
	googleMeasurementID := *flagGoogleMeasurementID
	googleAPISecret := *flagGoogleAPISecret
	databaseURL := *flagDatabaseURL
	username := *flagUsername
	password := *flagPassword
	uploadthingSecretKey := *flagUploadthingSecretKey
	if *flagEnvironment != "local" {
		data, err := os.ReadFile(*flagPort)
		if err != nil {
			log.Printf("Failed to read port from file, using flag value: %v", *flagPort)
		} else {
			port = string(data)
		}

		data, err = os.ReadFile(*flagClientURL)
		if err != nil {
			log.Printf("Failed to read client url from file, using flag value: %v", *flagClientURL)
		} else {
			clientURL = string(data)
		}

		data, err = os.ReadFile(*flagAuthToken)
		if err != nil {
			log.Printf("Failed to read auth token from file, using flag value: %v", *flagAuthToken)
		} else {
			authToken = string(data)
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

		data, err = os.ReadFile(*flagGoogleAPISecret)
		if err != nil {
			log.Printf("Failed to read google API secret from file, using flag value: %v", *flagGoogleAPISecret)
		} else {
			googleAPISecret = string(data)
		}

		data, err = os.ReadFile(*flagDatabaseURL)
		if err != nil {
			log.Printf("Failed to read database url string from file, using flag value: %v", *flagDatabaseURL)
		} else {
			databaseURL = string(data)
		}

		data, err = os.ReadFile(*flagUsername)
		if err != nil {
			log.Printf("Failed to read username from file, using flag value: %v", *flagUsername)
		} else {
			username = string(data)
		}

		data, err = os.ReadFile(*flagPassword)
		if err != nil {
			log.Printf("Failed to read password from file, using flag value: %v", *flagPassword)
		} else {
			password = string(data)
		}

		data, err = os.ReadFile(*flagUploadthingSecretKey)
		if err != nil {
			log.Printf("Failed to read uploadthing secret key from file, using flag value: %v", *flagUploadthingSecretKey)
		} else {
			uploadthingSecretKey = string(data)
		}
	} else {
		if port == "" {
			port = os.Getenv("PORT")
		}

		if clientURL == "" {
			clientURL = os.Getenv("CLIENT_URL")
		}

		if authToken == "" {
			authToken = os.Getenv("AUTH_TOKEN")
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

		if googleAPISecret == "" {
			googleAPISecret = os.Getenv("GOOGLE_API_SECRET")
		}

		if databaseURL == "" {
			databaseURL = os.Getenv("DATABASE_URL")
		}

		if username == "" {
			username = os.Getenv("USERNAME")
		}

		if password == "" {
			password = os.Getenv("PASSWORD")
		}

		if uploadthingSecretKey == "" {
			uploadthingSecretKey = os.Getenv("UPLOADTHING_SECRET_KEY")
		}
	}

	// Setup database
	database := database.NewDatabase(databaseURL)

	// Setup server
	s := server.New(
		server.Config{
			Environment:          *flagEnvironment,
			ClientURL:            clientURL,
			Port:                 port,
			AuthToken:            authToken,
			EmailJSServiceID:     emailJSServiceID,
			EmailJSTemplateID:    emailJSTemplateID,
			EmailJSPublicKey:     emailJSPublicKey,
			EmailJSPrivateKey:    emailJSPrivateKey,
			GoogleMeasurementID:  googleMeasurementID,
			GoogleAPISecret:      googleAPISecret,
			Username:             username,
			Password:             password,
			UploadthingSecretKey: uploadthingSecretKey,
			DatabaseAPI:          database,
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

	// Close database pool after graceful shutdown
	defer database.Close()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting...")
}
