package web3util

import (
	"fmt"
	"strings"
)

type DID struct {
	Method DIDMethod
	Namespace Namespace
	Address string
}

func NewDID(method DIDMethod, namespace Namespace, address string) *DID {
	return &DID{
		Method: method,
		Namespace: namespace,
		Address: address,
	}
}

func (d *DID) String() string {
	return fmt.Sprintf("did:%s:%s:%s", d.Method, d.Namespace, d.Address)
}

func DIDFromString(rawDID string) (*DID, error) {
	didElements := strings.Split(rawDID, ":")
	if len(didElements) != 4 {
		return nil, fmt.Errorf("failed converting did string to struct, invalid amount of path elements")
	}

	method := DIDMethod(didElements[1])
	if method != DIDMethodPKH {
		return nil, fmt.Errorf("failed converting did string to struct, invalid method")
	}

	namespace := Namespace(didElements[2])
	if namespace != NamespaceSolana {
		return nil, fmt.Errorf("failed converting did string to struct, invalid namespace")
	}

	address := didElements[3]
	if err := VerifySolanaPublicKey(address); err != nil {
		return nil, fmt.Errorf("failed converting did string to struct, invalid address")
	}

	return &DID{
		Method: method,
		Namespace: namespace,
		Address: address,
	}, nil
}
