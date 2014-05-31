package webutil

import (
	"github.com/hoisie/web"
	"github.com/saintfish/httpauth-go"
	"time"
)

type Digest httpauth.Digest

func NewDigest(realm, htdigest string, duration time.Duration) (*Digest, error) {
	digest, err := httpauth.NewDigest(realm, httpauth.OpenHtdigest(htdigest), false, nil)
	if err != nil {
		return nil, err
	}
	digest.ClientCacheResidence = duration
	return (*Digest)(digest), nil
}

func HandleAuth(digest *Digest, ctx *web.Context) bool {
	d := (*httpauth.Digest)(digest)
	if username := d.Authorize(ctx.Request); username == "" {
		d.NotifyAuthRequired(ctx.ResponseWriter, ctx.Request)
		return false
	}
	return true
}

func Logout(digest *Digest, ctx *web.Context) {
	d := (*httpauth.Digest)(digest)
	d.Logout(ctx.Request)
}
