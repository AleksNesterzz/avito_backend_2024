package changeBanner

import (
	"avito_backend/internal/http-server/responces"
	tokens "avito_backend/internal/lib"
	apierr "avito_backend/internal/lib/errors"
	"avito_backend/internal/models"
	"avito_backend/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type BannerChanger interface {
	ChangeBanner(ctx context.Context, bid int, cnt models.Content, act *bool) (string, error)
}

type Request struct {
	Tags      []int          `json:"tid,omitempty"`
	Fid       int            `json:"fid,omitempty"`
	Content   models.Content `json:"content,omitempty"`
	Is_active *bool          `json:"is_active,omitempty"`
}

type Response struct {
	Message string `json:"message"`
	Err     string `json:"error,omitempty"`
}

func New(banchanger BannerChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := tokens.GetToken(w, r)
		if err == apierr.ErrNoAuth {
			resp := Response{Message: "Пользователь не авторизован"}
			responces.JSONResponse(w, r, 401, resp)
			return
		} else if token == tokens.AdminToken {
			path := strings.Split(r.URL.Path, "/banner/")
			id, _ := strconv.Atoi(path[1])
			var req Request
			ctx := r.Context()
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				resp := Response{
					Message: "Некорректные данные"}
				responces.JSONResponse(w, r, 400, resp)
				return
			}
			_, err = banchanger.ChangeBanner(ctx, id, req.Content, req.Is_active)
			if errors.Is(err, storage.ErrBannerNotExists) {
				resp := Response{
					Message: "Баннер не найден"}
				responces.JSONResponse(w, r, 404, resp)
				return
			}
			if err != nil {
				resp := Response{
					Message: "Внутренняя ошибка сервера",
					Err:     err.Error()}
				responces.JSONResponse(w, r, 500, resp)
				return
			}
			resp := Response{Message: "OK"}

			responces.JSONResponse(w, r, 200, resp)
		} else {
			resp := Response{
				Message: "Пользователь не имеет доступа"}
			responces.JSONResponse(w, r, 403, resp)
			return
		}
	}
}
