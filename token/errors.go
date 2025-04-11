package token

import "errors"

var (
	ErrNotValidPublicAddr = &DIDTokenError{errors.New("invalid public address format")}
	ErrNotValidProof      = &DIDTokenError{errors.New("signature mismatch between 'proof' and 'claim'. Please " +
		"generate a new token with an intended issuer")}
	ErrExpired    = &DIDTokenError{errors.New("given DID token has expired. Please generate a new one")}
	ErrNbfExpired = &DIDTokenError{errors.New("given DID token cannot be used at this time. Please " +
		"check the 'nbf' field and regenerate a new token with a suitable value")}
	ErrAudMismatch = &DIDTokenError{errors.New("audience does not match client ID. Please ensure your " +
		"secret key matches the application which generated the DID token")}
)

type DIDTokenError struct {
	err error
}

func (e *DIDTokenError) Error() string {
	return e.err.Error()
}
