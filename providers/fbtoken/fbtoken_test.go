package fbtoken_test

import (
	// "fmt"
	"os"
	"testing"

	"github.com/stickbuilt/goth"
	"github.com/stickbuilt/goth/providers/fbtoken"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	provider := facebookProvider()
	a.Equal(provider.ClientKey, os.Getenv("FACEBOOK_KEY"))
	a.Equal(provider.Secret, os.Getenv("FACEBOOK_SECRET"))
}

func Test_Implements_Provider(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Implements((*goth.Verifier)(nil), facebookProvider())
}

/*
func Test_SessionFromJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := facebookProvider()

	s, err := provider.UnmarshalSession(`{"AuthURL":"http://facebook.com/auth_url","AccessToken":"1234567890"}`)
	a.NoError(err)
	session := s.(*facebook.Session)
	a.Equal(session.AuthURL, "http://facebook.com/auth_url")
	a.Equal(session.AccessToken, "1234567890")
}
*/
func facebookProvider() *fbtoken.Verifier {
	return fbtoken.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"))
}
