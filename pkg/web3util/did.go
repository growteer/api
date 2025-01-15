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

func (d DID) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *DID) UnmarshalJSON(data []byte) error {
	didElements := strings.Split(string(data), ":")
	if len(didElements) != 4 {
		return fmt.Errorf("failed unmarshalling did, invalid amount of path elements")
	}

	method := DIDMethod(didElements[1])
	if method != DIDMethodPKH {
		return fmt.Errorf("failed unmarshalling did, invalid method")
	}
	d.Method = method

	namespace := Namespace(didElements[2])
	if namespace != NamespaceSolana {
		return fmt.Errorf("failed unmarshalling did, invalid namespace")
	}
	d.Namespace = namespace

	address := didElements[3]
	if err := VerifySolanaPublicKey(address); err != nil {
		return fmt.Errorf("failed unmarshalling did, invalid address")
	}
	d.Address = address

	return nil
}
