#!/bin/bash

# Security Utilities for PGBridge-Go
# This script helps with common security tasks

show_help() {
    echo "PGBridge-Go Security Utilities"
    echo ""
    echo "Usage:"
    echo "  ./security-utils.sh hash-password [password]    - Generate bcrypt hash for password"
    echo "  ./security-utils.sh generate-key                - Generate 32-byte encryption key"
    echo "  ./security-utils.sh check-deps                  - Check for dependency vulnerabilities"
    echo "  ./security-utils.sh help                        - Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./security-utils.sh hash-password mypassword"
    echo "  ./security-utils.sh generate-key"
    echo ""
}

hash_password() {
    if [ -z "$1" ]; then
        echo "Error: Please provide a password to hash"
        echo "Usage: ./security-utils.sh hash-password [password]"
        exit 1
    fi
    
    echo "Generating bcrypt hash for password..."
    go run -c "
    package main
    import (
        \"fmt\"
        \"golang.org/x/crypto/bcrypt\"
        \"os\"
    )
    func main() {
        password := os.Args[1]
        hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            fmt.Println(\"Error:\", err)
            os.Exit(1)
        }
        fmt.Println(\"Bcrypt hash:\", string(hash))
    }" "$1"
}

generate_key() {
    echo "Generating 32-byte encryption key..."
    if command -v openssl >/dev/null 2>&1; then
        key=$(openssl rand -hex 32)
        echo "Generated key: $key"
        echo ""
        echo "Add this to your .env file:"
        echo "MASTER_KEY=$key"
    else
        echo "Error: openssl not found. Please install openssl or generate manually."
        echo "Alternative: use 'dd if=/dev/urandom bs=32 count=1 2>/dev/null | xxd -p -c 32'"
    fi
}

check_deps() {
    echo "Checking for dependency vulnerabilities..."
    cd src 2>/dev/null || cd .
    
    if command -v govulncheck >/dev/null 2>&1; then
        echo "Running govulncheck..."
        govulncheck ./...
    else
        echo "govulncheck not installed. Installing..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
        if command -v govulncheck >/dev/null 2>&1; then
            govulncheck ./...
        else
            echo "Failed to install govulncheck. Please install manually:"
            echo "go install golang.org/x/vuln/cmd/govulncheck@latest"
        fi
    fi
}

case "$1" in
    "hash-password")
        hash_password "$2"
        ;;
    "generate-key")
        generate_key
        ;;
    "check-deps")
        check_deps
        ;;
    "help"|"--help"|"-h"|"")
        show_help
        ;;
    *)
        echo "Error: Unknown command '$1'"
        echo ""
        show_help
        exit 1
        ;;
esac