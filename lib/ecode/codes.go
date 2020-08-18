package ecode

// FIXME currently one err code can be used in more than one place
// in the future try to use one code for one line of code
// this will allow us to track errors from logs to single lines of code

// all codes have implied fail or error
// it makes no sense to write this in error package ¯\_(ツ)_/¯

///go:generate stringer -type=ErrorCode
//type ErrorCode int

const (
	NoError = iota
	/// master
	// generic
	Unknown // FIXME not my code, not sure, please check it
	JsonDecode
	JsonEncode
	DbSave
	DbError
	OidError
	// auth/user
	WrongUsernameOrPassword
	WrongPassword
	WrongEmail
	PasswordsDoNotMatch
	EmailsDoNotMatch
	DuplicateUsername
	JwtGenerate
	GroupNotFound
	// host
	HostNotFound
	HostData
	HostNotTurnedOn
	// gs
	GsData
	GsUpdate
	GsStart
	GsStop
	GsCmd
	GsRestart
	GsInstall
	GsNotFound
	GsRemove
	GsNoLogs
	GsFileList
	// games
	GameNotFound
	///agent
	///wrapper
	TemplateGenerationFailed
)
