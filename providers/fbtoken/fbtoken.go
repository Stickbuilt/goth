// Package facebook implements the OAuth2 protocol for authenticating users through Facebook.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package fbtoken

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"github.com/stickbuilt/goth"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	authURL         string = "https://www.facebook.com/dialog/oauth"
	tokenURL        string = "https://graph.facebook.com/oauth/access_token"
	endpointProfile string = "https://graph.facebook.com/me?fields=email,first_name,last_name,link,bio,id,name,picture,location"
)

// New creates a new Facebook verifier, and sets up important connection details.
// You should always call `facebook.New` to get a new Verifier. Never try to create
// one manually.
func New(clientKey string, secret string, scopes ...string) *Verifier {
	p := &Verifier{
		ClientKey: clientKey,
		Secret:    secret,
	}
	p.config = newConfig(p, scopes)
	return p
}

// Verifier is the implementation of `goth.Verifier` for accessing Facebook.
type Verifier struct {
	ClientKey string
	Secret    string
	config    *oauth2.Config
}

// Name is the name used to retrieve this verifier later.
func (v *Verifier) Name() string {
	return "fbtoken"
}

// Debug is a no-op for the facebook package.
func (v *Verifier) Debug(debug bool) {}

func (v *Verifier) VerifyAuth(access_token string) (goth.User, error) {

	token := &oauth2.Token{AccessToken: access_token}
	client := v.config.Client(nil, token)
	res, err := client.Get(endpointProfile)

	user := goth.User{AccessToken: access_token}
	if err == nil && res.StatusCode == http.StatusOK {
		defer res.Body.Close()
		bits, err := ioutil.ReadAll(res.Body)
		if err == nil {
			err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
			if err == nil {
				err = userFromReader(bytes.NewReader(bits), &user)
			}
		}
	}
	return user, err
}

func userFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Bio     string `json:"bio"`
		Name    string `json:"name"`
		Link    string `json:"link"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.NickName = u.Name
	user.Email = u.Email
	user.Description = u.Bio
	user.AvatarURL = u.Picture.Data.URL
	user.UserID = u.ID
	user.Location = u.Location.Name

	return err
}

func newConfig(verifier *Verifier, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     verifier.ClientKey,
		ClientSecret: verifier.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	for _, scope := range scopes {
		c.Scopes = append(c.Scopes, scope)
	}

	return c
}
