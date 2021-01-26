package main

import "C"
import (
	"encoding/json"
	"github.com/privacybydesign/gabi"
	"github.com/privacybydesign/gabi/big"
	"gitlab.com/confiks/ctcl/common"
	"gitlab.com/confiks/ctcl/issuer"
)

// export GetIssuingNonce
func GenerateIssuerNonceB64() *C.Char {
	issuerNonceB64, err := json.Marshal(issuer.GenerateIssuerNonce())
	if err != nil {
		panic("Could not serialize issuer nonce")
	}

	return C.CString(issuerNonceB64)
}

// export Issue
func Issue(issuerPkXml, issuerSkXml, issuerNonceB64, commitmentsJson *C.Char) *C.Char {
	issuerNonce := new(big.Int)
	err := issuerNonce.UnmarshalJSON(C.GoString(issuerNonceB64))
	if err != nil {
		panic("Could not deserialize issuerNonce")
	}

	if issuerNonce.BitLen() != int(common.GabiSystemParameters.Lstatzk) {
		panic("Invalid length for issuerNonce")
	}

	// Commitments
	var commitments *gabi.IssueCommitmentMessage
	err = json.Unmarshal(commitmentsJson.GoString())
	if err != nil {
		panic("Could not deserialize commitments")
	}

	sig := issuer.Issue(C.GoString(issuerPkXml), C.GoString(issuerSkXml), issuerNonce, []string {"foo", "bar"}, commitments)
	return C.CString(json.Marshal(sig))
}
