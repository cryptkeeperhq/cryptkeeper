package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	KeycloakURL = "http://localhost:9999"
	ClientID    = "apiClient"
	ClientSecet = "4NZ9L23PRPEo51yZnclWCbXEyWU70rgU"
	Realm       = "myrealm"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.Config.Auth.SSOEnabled {
		err := h.KAuth.CreateUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}

	user.ID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	_, err = h.DB.Model(&user).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type KeycloakAuthRequest struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	Token        string `json:"token"`
}

func (h *Handler) AuthKeycloak(w http.ResponseWriter, r *http.Request) {

	var creds KeycloakAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	info, err := h.KAuth.GetUserInfo(creds.Token)

	fmt.Println(info)
	if err != nil {
		http.Error(w, "invalid token or User info", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: *info.PreferredUsername,
		UserID:   *info.Sub,
		AuthType: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(h.JWTKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    tokenString,
		"username": info.PreferredUsername,
		"name":     info.GivenName,
		"roles":    info.Profile,
	})

}

func (h *Handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User

	if h.Config.Auth.SSOEnabled {
		tokens, err := h.KAuth.Login(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Invalid username or password: "+err.Error(), http.StatusUnauthorized)
			return
		}
		info, err := h.KAuth.GetUserInfo(tokens.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		userGroups, _ := h.KAuth.GetGroups(*info.Sub)
		gR := h.KAuth.GetUserGroupsAndRoles(userGroups)
		fmt.Println(gR)

		user.Username = *info.PreferredUsername
		user.Name = *info.Name
		user.ID = *info.Sub
	} else {

		err := h.DB.Model(&user).Where("username = ?", creds.Username).Select()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid username", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: user.Username,
		UserID:   user.ID,
		AuthType: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JWTKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Generated token: %s", tokenString)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    tokenString,
		"username": user.Username,
		"name":     user.Name,
		"id":       user.ID,
		"roles":    []string{},
		// "roles":       gR,
	})

}

func (h *Handler) AuthenticateAppRole(w http.ResponseWriter, r *http.Request) {

	var req models.AppRole

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var appRole models.AppRole
	err := h.DB.Model(&appRole).Where("role_id = ? AND secret_id = ?", req.RoleID, req.SecretID).Select()
	if err != nil {
		http.Error(w, "Invalid app role ID or secret", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token for the app role
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		AuthType: "approle",
		UserID:   appRole.ID,
		Username: appRole.RoleID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	fmt.Println(claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JWTKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
	})
}
