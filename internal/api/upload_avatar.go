package api

import (
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/supabase/auth/internal/api/apierrors"
	"github.com/supabase/auth/internal/conf"
)

type Oss struct {
	mini *minio.Client
	c    *conf.MinioConfiguration
}

func NewOss(c *conf.MinioConfiguration) *Oss {
	client, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Id, c.Secret, ""),
		Secure: c.Secure,
		Region: c.Region,
	})
	if err != nil {
		panic(err)
	}
	o := &Oss{mini: client, c: c}
	return o
}
func (a *API) GetUploadAvatarURL(w http.ResponseWriter, r *http.Request) error {
	if a.oss == nil {
		return apierrors.NewInternalServerError("oss is disable")
	}
	ctx := r.Context()
	claims := getClaims(ctx)

	url, err := a.oss.mini.PresignedPutObject(r.Context(), a.oss.c.Bucket, claims.Subject, time.Hour)
	if err != nil {
		return apierrors.NewInternalServerError("Could not get url")
	}
	return sendJSON(w, http.StatusOK, map[string]string{"upload_avatar_url": url.String()})
}
