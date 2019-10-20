package main

import (
	"flag"
	"fmt"
	"gototp/src/lib"
	"gototp/src/totp"
)

func main() {

	// input flags
	fStep := flag.Int("step", 30, "Time step(seconds) after which new key is generated (default: 30)")
	fDigits := flag.Int("digits", 6, "Number of digits for the generated keys (default: 6)")
	fDigest := flag.String("digest", "none", "Digest to use, valid options(one of): sha1, sha256, sha512 (default: auto detects by key size)")
	fKey := flag.String("key", "", "Base32 encoded key to use (optional, if key is piped in)")
	flag.Parse()

	// collect the secret
	secret := *fKey
	// when secret is blank try to read from stdin
	if secret == "" {
		secret = lib.ReadStdinForSecret()
	}
	//other params
	step := *fStep
	digits := *fDigits
	digest := *fDigest
	// generate new instance
	totp := totp.TotpToken{
		Secret: secret,
		Step:   step,
		Digits: digits,
		Digest: digest,
	}

	token := totp.Generate()
	// report
	fmt.Printf("\nTOKEN: %v\n\n", token)

}
