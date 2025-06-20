package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/kelvinzer0/secure-email-validator/internal/checker"
	"github.com/kelvinzer0/secure-email-validator/internal/config"
)

func main() {
	var (
		email      = flag.String("email", "", "Email address to validate")
		verbose    = flag.Bool("verbose", false, "Enable verbose output")
		timeout    = flag.Int("timeout", 10, "SMTP connection timeout in seconds")
		jsonOutput = flag.Bool("json", false, "Output result in JSON format")
		server     = flag.Bool("server", false, "Run as HTTP server")
		port       = flag.String("port", "8587", "Server port (only used with -server)")
		help       = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Server mode
	if *server {
		startServer(*port)
		return
	}

	// CLI mode
	if *email == "" {
		fmt.Println("Error: Email address is required")
		showHelp()
		os.Exit(1)
	}

	cfg := &config.Config{
		Timeout: *timeout,
		Verbose: *verbose,
	}

	emailChecker := checker.NewEmailChecker(cfg)
	result := emailChecker.ValidateEmail(*email)

	if *jsonOutput {
		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonResult))
	} else {
		printHumanReadableResult(result, *email, *verbose)
	}

	if !result.Valid {
		os.Exit(1)
	}
}

func printHumanReadableResult(result *checker.ValidationResult, originalEmail string, verbose bool) {
	if result.Valid {
		fmt.Printf("✅ Email '%s' is valid and secure\n", originalEmail)
		if result.NormalizedEmail != originalEmail {
			fmt.Printf("📧 Normalized: %s\n", result.NormalizedEmail)
		}
		fmt.Printf("✨ Reason: %s\n", result.Reason)
	} else {
		fmt.Printf("❌ Email '%s' is invalid\n", originalEmail)
		if result.NormalizedEmail != originalEmail {
			fmt.Printf("📧 Normalized: %s\n", result.NormalizedEmail)
		}
		fmt.Printf("⚠️  Reason: %s\n", result.Reason)
	}

	if verbose {
		fmt.Println("\n--- Detailed Information ---")
		fmt.Printf("Domain: %s\n", result.Domain)
		fmt.Printf("Has MX Record: %t\n", result.HasMXRecord)
		fmt.Printf("Has DNSSEC: %t\n", result.HasDNSSEC)
		fmt.Printf("Primary MX Server: %s\n", result.PrimaryMXServer)
		fmt.Printf("Supports STARTTLS: %t\n", result.SupportsSTARTTLS)
	}
}

func startServer(port string) {
	http.HandleFunc("/validate", handleValidation)
	http.HandleFunc("/health", handleHealth)
	
	fmt.Printf("🚀 Secure Email Validator Server starting on port %s\n", port)
	fmt.Printf("📍 Validation endpoint: http://localhost:%s/validate?email=test@example.com\n", port)
	fmt.Printf("💚 Health check: http://localhost:%s/health\n", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}
}

func handleValidation(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, `{"error": "Email parameter required"}`, http.StatusBadRequest)
		return
	}

	timeoutStr := r.URL.Query().Get("timeout")
	timeout := 10
	if timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		}
	}

	verbose := r.URL.Query().Get("verbose") == "true"

	cfg := &config.Config{
		Timeout: timeout,
		Verbose: verbose,
	}

	emailChecker := checker.NewEmailChecker(cfg)
	result := emailChecker.ValidateEmail(email)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "secure-email-validator",
	})
}

func showHelp() {
	fmt.Println("Secure Email Validator")
	fmt.Println("A professional tool for email security validation")
	fmt.Println("Usage: secure-email-validator [options]")
	fmt.Println("\nModes:")
	fmt.Println("  CLI Mode (default): Validate single email")
	fmt.Println("  Server Mode: Run as HTTP API server")
	fmt.Println("\nCLI Options:")
	fmt.Println("  -email string     Email address to validate (required for CLI mode)")
	fmt.Println("  -verbose         Enable verbose output")
	fmt.Println("  -timeout int     SMTP connection timeout in seconds (default 10)")
	fmt.Println("  -json           Output result in JSON format")
	fmt.Println("  -help           Show this help message")
	fmt.Println("\nServer Options:")
	fmt.Println("  -server         Run as HTTP server")
	fmt.Println("  -port string    Server port (default 8587)")
	fmt.Println("\nExamples:")
	fmt.Println("  # CLI mode")
	fmt.Println("  email-checker -email john.doe@gmail.com -verbose")
	fmt.Println("  email-checker -email test@example.com -json")
	fmt.Println("")
	fmt.Println("  # Server mode")
	fmt.Println("  email-checker -server -port 3000")
	fmt.Println("  curl 'http://localhost:8587/validate?email=test@gmail.com&verbose=true'")
}
