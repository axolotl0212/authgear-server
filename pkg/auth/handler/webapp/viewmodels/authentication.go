package viewmodels

import (
	"github.com/authgear/authgear-server/pkg/auth/dependency/identity"
)

// Ideally we should use type alias to present
// LoginPageTextLoginIDVariant and LoginPageTextLoginIDInputType
// But they may be passed to localize which does not support type alias of builtin types.

const (
	LoginPageTextLoginIDVariantNone            = "none"
	LoginPageTextLoginIDVariantEamilOrUsername = "email_or_username"
	LoginPageTextLoginIDVariantEmail           = "email"
	LoginPageTextLoginIDVariantUsername        = "username"
)

const (
	LoginPageTextLoginIDInputTypeText  = "text"
	LoginPageTextLoginIDInputTypeEmail = "email"
)

type AuthenticationViewModel struct {
	IdentityCandidates            []identity.Candidate
	LoginPageLoginIDHasPhone      bool
	LoginPageTextLoginIDVariant   string
	LoginPageTextLoginIDInputType string
}

// func (m *AuthenticationViewModeler) ViewModel(r *http.Request) AuthenticationViewModel {
// 	userID := ""
// 	if sess := authn.GetSession(r.Context()); sess != nil {
// 		userID = sess.AuthnAttrs().UserID
// 	}
//
// 	identityCandidates, err := m.Identity.ListCandidates(userID)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	hasEmail := false
// 	hasUsername := false
// 	hasPhone := false
// 	for _, c := range identityCandidates {
// 		if c[identity.CandidateKeyType] == string(authn.IdentityTypeLoginID) {
// 			if c[identity.CandidateKeyLoginIDType] == "phone" {
// 				c["login_id_input_type"] = "phone"
// 				hasPhone = true
// 			} else if c[identity.CandidateKeyLoginIDType] == "email" {
// 				c["login_id_input_type"] = "email"
// 				hasEmail = true
// 			} else {
// 				c["login_id_input_type"] = "text"
// 				hasUsername = true
// 			}
// 		}
// 	}
//
// 	var loginPageTextLoginIDVariant string
// 	var loginPageTextLoginIDInputType string
// 	if hasEmail {
// 		if hasUsername {
// 			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantEamilOrUsername
// 			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeText
// 		} else {
// 			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantEmail
// 			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeEmail
// 		}
// 	} else {
// 		if hasUsername {
// 			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantUsername
// 			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeText
// 		} else {
// 			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantNone
// 			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeText
// 		}
// 	}
//
// 	return AuthenticationViewModel{
// 		IdentityCandidates:            identityCandidates,
// 		LoginPageLoginIDHasPhone:      hasPhone,
// 		LoginPageTextLoginIDVariant:   loginPageTextLoginIDVariant,
// 		LoginPageTextLoginIDInputType: loginPageTextLoginIDInputType,
// 	}
// }
