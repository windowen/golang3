package errs

var (
	ErrArgs             = NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission     = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrRemotePlaceError = NewCodeError(RemotePlaceError, "RemotePlaceError")
	ErrDatabase         = NewCodeError(DatabaseError, "DatabaseError")
	ErrInternalServer   = NewCodeError(InternalSystemError, "InternalSystemError")
	ErrNetwork          = NewCodeError(NetworkError, "NetworkError")

	ErrData   = NewCodeError(DataError, "DataError")
	ErrUser   = NewCodeError(UserError, "UserError")
	ErrUpload = NewCodeError(UploadError, "UploadError")
	ErrFollow = NewCodeError(FollowErrExist, "UploadError")
)
