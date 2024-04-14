package deleteBanner

import (
	"avito_backend/internal/http-server/responces"
	tokens "avito_backend/internal/lib"
	apierr "avito_backend/internal/lib/errors"
	"avito_backend/internal/storage"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type BannerDeleter interface {
	DeleteBanner(ctx context.Context, banner_id int) (int64, error)
}

type Response struct {
	Id      int    `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

func New(bannerdeleter BannerDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.Split(r.URL.Path, "/banner/")
		token, err := tokens.GetToken(w, r)
		if err == apierr.ErrNoAuth {
			resp := Response{
				Message: "user unathorized"}
			responces.JSONResponse(w, r, 401, resp)
			return
		} else if token == tokens.AdminToken {
			banner_id, err := strconv.Atoi(path[1])
			if err != nil {
				resp := Response{Message: "Некорректные данные"}
				responces.JSONResponse(w, r, 400, resp)
				return
			}
			ctx := r.Context()

			_, err = bannerdeleter.DeleteBanner(ctx, banner_id)
			if errors.Is(err, storage.ErrBannerNotFound) {

				resp := Response{
					Message: "banner not found"}
				responces.JSONResponse(w, r, 404, resp)
				return
			}
			if err != nil {
				resp := Response{
					Message: "Внутренняя ошибка сервера"}
				responces.JSONResponse(w, r, 500, resp)

				return
			}
			resp := Response{Message: "banner was succesfully deleted"}
			responces.JSONResponse(w, r, 204, resp)
		} else {
			resp := Response{
				Message: "user has no access",
			}
			responces.JSONResponse(w, r, 403, resp)
		}
	}
}
