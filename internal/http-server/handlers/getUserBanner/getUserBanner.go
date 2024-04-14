package getUserBanner

import (
	"context"
	"net/http"
	"strconv"

	"avito_backend/internal/http-server/responces"
	tokens "avito_backend/internal/lib"
	apierr "avito_backend/internal/lib/errors"
	models "avito_backend/internal/models"
	"avito_backend/internal/storage"
)

type Response struct {
	Cnt     *models.Content `json:"content,omitempty"`
	Message string          `json:"message"`
}

type UserBannerGetter interface {
	GetUserBanner(ctx context.Context, tag, fid int, last bool) (*models.Content, error)
}

func New(bannergetter UserBannerGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := tokens.GetToken(w, r)
		if err == apierr.ErrNoAuth {
			resp := Response{Message: "User unathorized"}
			responces.JSONResponse(w, r, 401, resp)
		} else if token == tokens.CasualToken {
			ctx := r.Context()
			tag, err := strconv.Atoi(r.URL.Query().Get("tid"))
			if err != nil {
				resp := Response{
					Message: "Некорректные данные",
				}
				responces.JSONResponse(w, r, 400, resp)
				return
			}
			feature, err := strconv.Atoi(r.URL.Query().Get("fid"))
			if err != nil {
				resp := Response{
					Message: "Некорректные данные",
				}
				responces.JSONResponse(w, r, 400, resp)
				return
			}
			last, _ := strconv.ParseBool(r.URL.Query().Get("last_revision"))

			banner, err := bannergetter.GetUserBanner(ctx, tag, feature, last)
			if err != nil {
				if err == storage.ErrBannerNotFound {
					resp := Response{
						Message: "Баннер не найден",
					}
					responces.JSONResponse(w, r, 404, resp)
					return
				} else {
					resp := Response{
						Message: "Внутренняя ошибка сервера"}
					responces.JSONResponse(w, r, 500, resp)

					return
				}

			}
			resp := Response{
				Cnt:     banner,
				Message: "OK",
			}
			responces.JSONResponse(w, r, 200, resp)
		} else {
			resp := Response{
				Message: "Доступ запрещен"}
			responces.JSONResponse(w, r, 403, resp)
		}
	}
}
