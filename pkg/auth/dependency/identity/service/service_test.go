package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/authgear/authgear-server/pkg/auth/dependency/identity"
	"github.com/authgear/authgear-server/pkg/auth/dependency/identity/loginid"
	"github.com/authgear/authgear-server/pkg/auth/dependency/identity/oauth"
	"github.com/authgear/authgear-server/pkg/lib/authn"
	"github.com/authgear/authgear-server/pkg/lib/config"
)

func TestProviderListCandidates(t *testing.T) {
	Convey("Provider ListCandidates", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loginIDProvider := NewMockLoginIDIdentityProvider(ctrl)
		oauthProvider := NewMockOAuthIdentityProvider(ctrl)

		p := &Service{
			Authentication: &config.AuthenticationConfig{},
			Identity: &config.IdentityConfig{
				LoginID: &config.LoginIDConfig{},
				OAuth:   &config.OAuthSSOConfig{},
			},
			LoginID: loginIDProvider,
			OAuth:   oauthProvider,
		}

		Convey("no candidates", func() {
			actual, err := p.ListCandidates("")
			So(err, ShouldBeNil)
			So(actual, ShouldBeEmpty)
		})

		Convey("oauth", func() {
			p.Authentication.Identities = []authn.IdentityType{authn.IdentityTypeOAuth}
			p.Identity.OAuth.Providers = []config.OAuthSSOProviderConfig{
				{
					Alias: "google",
					Type:  "google",
				},
			}

			actual, err := p.ListCandidates("")
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, []identity.Candidate{
				{
					"identity_id":         "",
					"type":                "oauth",
					"display_id":          "",
					"provider_type":       "google",
					"provider_alias":      "google",
					"provider_subject_id": "",
				},
			})
		})

		Convey("loginid", func() {
			p.Authentication.Identities = []authn.IdentityType{authn.IdentityTypeLoginID}
			p.Identity.LoginID.Keys = []config.LoginIDKeyConfig{
				{
					Type: "email",
					Key:  "email",
				},
			}

			actual, err := p.ListCandidates("")
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, []identity.Candidate{
				{
					"identity_id":    "",
					"type":           "login_id",
					"display_id":     "",
					"login_id_type":  "email",
					"login_id_key":   "email",
					"login_id_value": "",
				},
			})
		})

		Convey("respect authentication", func() {
			p.Identity.OAuth.Providers = []config.OAuthSSOProviderConfig{
				{
					Alias: "google",
					Type:  "google",
				},
			}
			p.Identity.LoginID.Keys = []config.LoginIDKeyConfig{
				{
					Type: "email",
					Key:  "email",
				},
			}

			actual, err := p.ListCandidates("")
			So(err, ShouldBeNil)
			So(actual, ShouldBeEmpty)
		})

		Convey("associate login ID identity", func() {
			userID := "a"

			p.Authentication.Identities = []authn.IdentityType{authn.IdentityTypeLoginID}
			p.Identity.LoginID.Keys = []config.LoginIDKeyConfig{
				{
					Type: "email",
					Key:  "email",
				},
			}

			loginIDProvider.EXPECT().List(userID).Return([]*loginid.Identity{
				{
					LoginIDKey:      "email",
					LoginID:         "john.doe@example.com",
					OriginalLoginID: "john.doe@example.com",
					Claims: map[string]string{
						"email": "john.doe@example.com",
					},
				},
			}, nil)
			oauthProvider.EXPECT().List(userID).Return(nil, nil)

			actual, err := p.ListCandidates(userID)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, []identity.Candidate{
				{
					"identity_id":    "",
					"type":           "login_id",
					"display_id":     "john.doe@example.com",
					"login_id_type":  "email",
					"login_id_key":   "email",
					"login_id_value": "john.doe@example.com",
				},
			})
		})

		Convey("associate oauth identity", func() {
			userID := "a"

			p.Authentication.Identities = []authn.IdentityType{authn.IdentityTypeOAuth}
			p.Identity.OAuth.Providers = []config.OAuthSSOProviderConfig{
				{
					Alias: "google",
					Type:  "google",
				},
			}

			loginIDProvider.EXPECT().List(userID).Return(nil, nil)
			oauthProvider.EXPECT().List(userID).Return([]*oauth.Identity{
				{
					ProviderID: config.ProviderID{
						Type: "google",
						Keys: map[string]interface{}{},
					},
					ProviderSubjectID: "john.doe@gmail.com",
					Claims: map[string]interface{}{
						"email": "john.doe@gmail.com",
					},
				},
			}, nil)

			actual, err := p.ListCandidates(userID)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, []identity.Candidate{
				{
					"identity_id":         "",
					"type":                "oauth",
					"display_id":          "john.doe@gmail.com",
					"provider_type":       "google",
					"provider_alias":      "google",
					"provider_subject_id": "john.doe@gmail.com",
				},
			})
		})
	})
}
