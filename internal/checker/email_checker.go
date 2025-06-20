package checker

import (
	"fmt"
	"net"
	"net/smtp"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/kelvinzer0/secure-email-validator/internal/config"
)

// EmailChecker handles email validation
type EmailChecker struct {
	config *config.Config
}

// ValidationResult contains the result of email validation
type ValidationResult struct {
	Valid             bool   `json:"valid"`
	Reason            string `json:"reason"`
	NormalizedEmail   string `json:"normalized_email"`
	Domain            string `json:"domain"`
	HasMXRecord       bool   `json:"has_mx_record"`
	HasDNSSEC         bool   `json:"has_dnssec"`
	PrimaryMXServer   string `json:"primary_mx_server"`
	SupportsSTARTTLS  bool   `json:"supports_starttls"`
}

// NewEmailChecker creates a new EmailChecker instance
func NewEmailChecker(cfg *config.Config) *EmailChecker {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	return &EmailChecker{
		config: cfg,
	}
}

// ValidateEmail performs comprehensive email validation
func (ec *EmailChecker) ValidateEmail(email string) *ValidationResult {
	result := &ValidationResult{}
	
	// Basic email format validation
	if !ec.isValidEmailFormat(email) {
		result.Valid = false
		result.Reason = "Invalid email format"
		result.NormalizedEmail = email
		return result
	}

	// Normalize email (especially Gmail)
	result.NormalizedEmail = ec.normalizeEmail(email)
	result.Domain = ec.extractDomain(result.NormalizedEmail)

	// Check MX record
	result.HasMXRecord = ec.hasMXRecord(result.Domain)
	if !result.HasMXRecord {
		result.Valid = false
		result.Reason = "Domain doesn't have MX record"
		return result
	}

	// Check DNSSEC
	result.HasDNSSEC = ec.hasDNSSEC(result.Domain)
	if !result.HasDNSSEC {
		result.Valid = false
		result.Reason = "Domain doesn't support DNSSEC"
		return result
	}

	// Get primary MX server
	result.PrimaryMXServer = ec.getPrimaryMXServer(result.Domain)
	if result.PrimaryMXServer == "" {
		result.Valid = false
		result.Reason = "Failed to get MX server"
		return result
	}

	// Check STARTTLS support
	result.SupportsSTARTTLS = ec.smtpSupportsSTARTTLS(result.PrimaryMXServer)
	if !result.SupportsSTARTTLS {
		result.Valid = false
		result.Reason = "SMTP server doesn't support STARTTLS"
		return result
	}

	result.Valid = true
	result.Reason = "Email is valid and domain supports secure mail delivery"
	return result
}

func (ec *EmailChecker) isValidEmailFormat(email string) bool {
	// RFC 5322 compliant email regex (simplified)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return emailRegex.MatchString(email) && len(email) <= 254
}

func (ec *EmailChecker) normalizeEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local, domain := parts[0], parts[1]
	
	// Normalize Gmail and Googlemail domains
	if domain == "gmail.com" || domain == "googlemail.com" {
		// Remove dots from local part
		local = strings.ReplaceAll(local, ".", "")
		// Remove everything after + (Gmail alias feature)
		if plusIndex := strings.Index(local, "+"); plusIndex != -1 {
			local = local[:plusIndex]
		}
		domain = "gmail.com" // Normalize googlemail.com to gmail.com
	}

	return local + "@" + domain
}

func (ec *EmailChecker) extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(parts[1]))
}

func (ec *EmailChecker) hasMXRecord(domain string) bool {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		if ec.config.Verbose {
			fmt.Printf("MX lookup failed for %s: %v\n", domain, err)
		}
		return false
	}
	return len(mxRecords) > 0
}

func (ec *EmailChecker) hasDNSSEC(domain string) bool {
	// Use dig command to check DNSSEC
	cmd := exec.Command("dig", "+dnssec", "+short", "SOA", domain)
	output, err := cmd.Output()
	if err != nil {
		if ec.config.Verbose {
			fmt.Printf("DNSSEC check failed for %s: %v\n", domain, err)
		}
		return false
	}
	
	outputStr := string(output)
	// Look for RRSIG record which indicates DNSSEC is enabled
	return strings.Contains(outputStr, "RRSIG") || 
		   strings.Contains(strings.ToUpper(outputStr), "RRSIG")
}

func (ec *EmailChecker) getPrimaryMXServer(domain string) string {
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		if ec.config.Verbose {
			fmt.Printf("Failed to get MX records for %s: %v\n", domain, err)
		}
		return ""
	}

	// Find MX record with lowest priority number (highest priority)
	var primaryMX *net.MX
	for _, mx := range mxRecords {
		if primaryMX == nil || mx.Pref < primaryMX.Pref {
			primaryMX = mx
		}
	}

	if primaryMX != nil {
		// Remove trailing dot if present
		host := primaryMX.Host
		if strings.HasSuffix(host, ".") {
			host = host[:len(host)-1]
		}
		return host
	}

	return ""
}

func (ec *EmailChecker) smtpSupportsSTARTTLS(mxServer string) bool {
	timeout := time.Duration(ec.config.Timeout) * time.Second
	
	// Try to connect to SMTP server on port 25
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:25", mxServer), timeout)
	if err != nil {
		if ec.config.Verbose {
			fmt.Printf("Failed to connect to SMTP server %s: %v\n", mxServer, err)
		}
		return false
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, mxServer)
	if err != nil {
		if ec.config.Verbose {
			fmt.Printf("Failed to create SMTP client for %s: %v\n", mxServer, err)
		}
		return false
	}
	defer client.Quit()

	// Check if STARTTLS extension is supported
	ok, _ := client.Extension("STARTTLS")
	return ok
}
