package getBanner

import (
	"context"
	"net/http"
	"strconv"

	"avito_backend/internal/http-server/responces"
	tokens "avito_backend/internal/lib"
	apierr "avito_backend/internal/lib/errors"
	models "avito_backend/internal/models"
)

type Response struct {
	Banner  []models.Banner `json:"banners,omitempty"`
	Message string          `json:"error,omitempty"`
}

type BannerGetter interface {
	GetBanner(ctx context.Context, tag, fid, limit, offset *int) ([]models.Banner, error)
}

func New(bannergetter BannerGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tid, fid, limit, offset *int
		token, err := tokens.GetToken(w, r)
		if err == apierr.ErrNoAuth {
			resp := Response{
				Message: "User unathorized"}
			responces.JSONResponse(w, r, 401, resp)
		} else if token == tokens.AdminToken {
			ctx := r.Context()
			if r.URL.Query().Get("tid") == "" {
				tid = nil
			} else {
				tid_n, _ := strconv.Atoi(r.URL.Query().Get("tid"))
				tid = &tid_n
			}
			if r.URL.Query().Get("fid") == "" {
				fid = nil
			} else {
				fid_n, _ := strconv.Atoi(r.URL.Query().Get("fid"))
				fid = &fid_n
			}
			if r.URL.Query().Get("limit") == "" {
				limit_n := 50 //значение лимита по умолчанию
				limit = &limit_n
			} else {
				limit_n, _ := strconv.Atoi(r.URL.Query().Get("limit"))
				limit = &limit_n
			}
			if r.URL.Query().Get("offset") == "" {
				offset = nil
			} else {
				offset_n, _ := strconv.Atoi(r.URL.Query().Get("offset"))
				offset = &offset_n
			}

			banner, err := bannergetter.GetBanner(ctx, tid, fid, limit, offset)
			if err != nil {
				resp := Response{
					Message: "failed to get user banners"}
				responces.JSONResponse(w, r, 500, resp)
				return
			}
			resp := Response{Banner: banner}
			responces.JSONResponse(w, r, 200, resp)
		} else {
			resp := Response{
				Message: "user has no access"}
			responces.JSONResponse(w, r, 403, resp)
			return
		}
	}
}
