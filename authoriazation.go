package tenuki

import "net/http"

// authorizationHelper

type TokenAuthorization struct {
	Token string
}

func (a *TokenAuthorization) Apply(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+a.Token)
}
