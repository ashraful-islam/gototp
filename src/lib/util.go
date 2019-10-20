package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// handle errors
func CheckErr(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg, err)
		os.Exit(1)
	}
}

// Fixes provided secret by removing non-base32 characters
func FixKey(key string) string {
	// remove preceeding/trailing whitespace and uppercase
	cleanKey := strings.ToUpper(key)
	// non base32 chars and whitespaces pattern
	rexp := regexp.MustCompile(`[^A-Z0-9=]`)
	// cleanup
	cleanKey = rexp.ReplaceAllString(cleanKey, "")
	return strings.TrimSpace(cleanKey)
}

// Appends '=' character suffix at the end
func PadBase32(key string) string {

	n := (8 - len(key)) % 8
	// when -ve we need to fix to positive number(i.e. padding should be always positive)
	if n < 0 {
		n += 8
	}
	suffix := strings.Repeat("=", n)
	return (key + suffix)
}

// Tries to read one line from Stdin as secret(totp seed)
func ReadStdinForSecret() string {
	r := bufio.NewReader(os.Stdin)
	secret := ""
	for {
		// expected only single read
		line, err := r.ReadString('\n')

		if len(strings.TrimSpace(line)) > 0 {
			secret = strings.TrimSpace(line)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			CheckErr(err, "Failed to read key from stdin")
		}
	}

	return secret
}
