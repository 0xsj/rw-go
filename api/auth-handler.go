package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/0xsj/rw-go/db/sqlc"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

func (req *UserRegisterRequest) bind(r *http.Request, p *db.CreateUserParams) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	p.Username = req.User.Username
	p.Email = req.User.Email
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Password = string(hashed)
	return nil
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var (
		req UserRegisterRequest
		p   db.CreateUserParams
	)
	if err := req.bind(r, &p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := s.store.CreateUser(r.Context(), p)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			log.Printf("error code: %s", pgErr.Code)
			log.Printf("error message: %s", pgErr.Message)
			log.Printf("error detail: %s", pgErr.Detail)
			log.Printf("error hint: %s", pgErr.Hint)
			log.Printf("error position: %s", pgErr.Position)
			log.Printf("error internal position: %s", pgErr.InternalPosition)
			log.Printf("error internal query: %s", pgErr.InternalQuery)
			log.Printf("error where: %s", pgErr.Where)
			log.Printf("error schema name: %s", pgErr.Schema)
			log.Printf("error table name: %s", pgErr.Table)
			log.Printf("error column name: %s", pgErr.Column)
			log.Printf("error data type name: %s", pgErr.DataTypeName)
			log.Printf("error constraint name: %s", pgErr.Constraint)
			log.Printf("error file: %s", pgErr.File)
			log.Printf("error line: %s", pgErr.Line)
			log.Printf("error routine: %s", pgErr.Routine)

			if pgErr.Code == "23505" {
				http.Error(w, "user already exists", http.StatusUnprocessableEntity)
				return
			}
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user})

}
