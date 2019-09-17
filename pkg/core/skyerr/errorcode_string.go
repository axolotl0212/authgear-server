// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT.

package skyerr

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NotAuthenticated-101]
	_ = x[PermissionDenied-102]
	_ = x[AccessKeyNotAccepted-103]
	_ = x[AccessTokenNotAccepted-104]
	_ = x[InvalidCredentials-105]
	_ = x[BadRequest-106]
	_ = x[InvalidArgument-107]
	_ = x[Duplicated-108]
	_ = x[ResourceNotFound-109]
	_ = x[UndefinedOperation-110]
	_ = x[PasswordPolicyViolated-111]
	_ = x[UserDisabled-112]
	_ = x[VerificationRequired-113]
	_ = x[WebHookTimeOut-114]
	_ = x[WebHookFailed-115]
	_ = x[CurrentIdentityBeingDeleted-116]
	_ = x[AuthenticationSession-117]
	_ = x[InvalidAuthenticationSession-118]
	_ = x[UnexpectedError-10000]
	_ = x[UnexpectedAuthInfoNotFound-10001]
}

const (
	_ErrorCode_name_0 = "NotAuthenticatedPermissionDeniedAccessKeyNotAcceptedAccessTokenNotAcceptedInvalidCredentialsBadRequestInvalidArgumentDuplicatedResourceNotFoundUndefinedOperationPasswordPolicyViolatedUserDisabledVerificationRequiredWebHookTimeOutWebHookFailedCurrentIdentityBeingDeletedAuthenticationSessionInvalidAuthenticationSession"
	_ErrorCode_name_1 = "UnexpectedErrorUnexpectedAuthInfoNotFound"
)

var (
	_ErrorCode_index_0 = [...]uint16{0, 16, 32, 52, 74, 92, 102, 117, 127, 143, 161, 183, 195, 215, 229, 242, 269, 290, 318}
	_ErrorCode_index_1 = [...]uint8{0, 15, 41}
)

func (i ErrorCode) String() string {
	switch {
	case 101 <= i && i <= 118:
		i -= 101
		return _ErrorCode_name_0[_ErrorCode_index_0[i]:_ErrorCode_index_0[i+1]]
	case 10000 <= i && i <= 10001:
		i -= 10000
		return _ErrorCode_name_1[_ErrorCode_index_1[i]:_ErrorCode_index_1[i+1]]
	default:
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
