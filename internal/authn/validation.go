package authn

import "regexp"

var ethereumAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

func IsValidEthereumAddress(address string) bool {
	return ethereumAddressRegex.MatchString(address)
}