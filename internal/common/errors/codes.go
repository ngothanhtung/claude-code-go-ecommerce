package errors

// Internal error codes.
const (
	CodeInternal          = 5000
	CodeValidation        = 4000
	CodeUnauthorized      = 4001
	CodeForbidden         = 4003
	CodeNotFound          = 4004
	CodeConflict          = 4009
	CodeInvalidCredentials = 4010
	CodeTokenExpired      = 4011
	CodeTokenInvalid      = 4012
	CodeEmailExists       = 4013
	CodeRateLimited       = 4029
	CodeUploadTooLarge    = 4030
	CodeUploadUnsupported = 4031
)
