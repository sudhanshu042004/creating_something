package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sudhanshu042004/sandbox/internal/storage"
	"github.com/sudhanshu042004/sandbox/internal/types"
	"github.com/sudhanshu042004/sandbox/internal/utils/hashing"
	"github.com/sudhanshu042004/sandbox/internal/utils/response"
	"github.com/sudhanshu042004/sandbox/internal/utils/token"
)

func SignUp(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.User

		//decoding in json
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("body is empty")))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		//validation
		if err := validator.New().Struct(user); err != nil {
			validationErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationErr))
			return
		}

		//bcrypt
		hashedPassword, err := hashing.HashPassword(user.Password)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		//save into db
		lastId, err := storage.CreateUser(user.Name, user.Email, hashedPassword)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("user successfully created", slog.String("userId", fmt.Sprint(lastId)))

		//token
		tokenString, err := token.CreateToken(lastId, user.Email)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		token.SetAuthCookie(w, tokenString)
		response.WriteJson(w, http.StatusOK, map[string]string{"message": "successs"})
	}
}

func Login(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.User

		//decoding in json
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("body is empty")))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//validation
		if err := validator.New().Struct(user); err != nil {
			validationErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationErr))
			return
		}
		//verifying
		existingUser, err := storage.GetUser(user.Email)
		if err != nil {
			if errors.Is(err, response.ErrUserNotFound) {
				response.WriteJson(w, http.StatusNotFound, err)
				return
			}
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		isValid := hashing.CheckPassword(existingUser.Password, user.Password)
		if !isValid {
			response.WriteJson(w, http.StatusForbidden, errors.New("invalid credentials"))
			return
		}

		//token
		tokenString, err := token.CreateToken(existingUser.Id, existingUser.Email)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		token.SetAuthCookie(w, tokenString)
		response.WriteJson(w, http.StatusOK, map[string]string{"message": "successs"})
	}
}
