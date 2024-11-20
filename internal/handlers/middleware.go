package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/db"
	"github.com/cryptkeeperhq/cryptkeeper/internal/metrics"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/dgrijalva/jwt-go"
)

const requestIDKey = "REQUEST_ID"

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenStr = tokenStr[len("Bearer "):]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return h.JWTKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var identity utils.Identity
		switch claims.AuthType {
		case "user":
			identity = utils.UserIdentity{
				UserID:   claims.UserID,
				Username: claims.Username,
			}
		case "approle":
			identity = utils.AppRoleIdentity{
				AppRoleID: claims.UserID,
				RoleName:  claims.UserID,
			}
		default:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store identity in context
		ctx := context.WithValue(r.Context(), "identity", identity)

		defer func() {
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}

			metrics.RequestCount.WithLabelValues(r.URL.Path).Inc()
			metrics.RequestDuration.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())

			h.Config.Logger.Info(
				"Request",
				"request_id", requestID,
				"method", r.Method,
				"uri", r.RequestURI,
				"path", r.URL.Path,
				"addr", r.RemoteAddr,
				"ua", r.UserAgent(),
				"duration", time.Since(start),
			)
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) CheckPermission(permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		key := r.URL.Query().Get("key")

		identity, ok := r.Context().Value("identity").(utils.Identity)
		if !ok {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		isAllowed := h.isAllowed(identity, path, key, permission)

		if !isAllowed {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) isAllowed(identity utils.Identity, path string, key string, permission string) bool {
	// Retrieve path details once to avoid repeated DB calls
	pathDetails, err := db.GetPathByName(path)
	if err != nil {
		return false // Path not found, immediately deny access
	}

	fmt.Println("HERE ", identity.GetAuthType(), identity.GetID(), identity.GetUsername(), path, key, permission)

	// Fetch user groups and policy information once
	var userGroups []string
	err = h.DB.Model((*models.Group)(nil)).
		Distinct().
		Join("JOIN user_groups ON id = user_groups.group_id").
		Where("user_groups.user_id = ?", identity.GetID()). // Filter groups by user ID
		Column("name").
		Select(&userGroups)

	if err != nil {
		return false
	}

	var policy models.Policy
	err = h.DB.Model(&policy).Where("path_id = ?", pathDetails.ID).First()
	if err != nil {
		return false
	}

	// Check if any policy paths allow the requested permission
	isAllowed := h.checkPolicyPaths(policy.Paths, identity, userGroups, permission)

	// Check if any secret paths deny the requested permission
	if key != "" && isAllowed {
		isAllowed = h.checkSecretPaths(policy.Secrets, pathDetails.ID, identity, userGroups, key, permission)
	}

	return isAllowed
}

// Helper to check policy paths for user/group permissions
func (h *Handler) checkPolicyPaths(paths []models.PolicyPath, identity utils.Identity, userGroups []string, permission string) bool {
	for _, policyPath := range paths {
		// Check if user or any group has permission for this path
		if utils.Contains(policyPath.Users, identity.GetUsername()) || h.hasGroupAccess(policyPath.Groups, userGroups) {

			if utils.Contains(policyPath.DenyPermissions, permission) {
				return false
			}

			if utils.Contains(policyPath.Permissions, permission) {
				return true
			}
		}

		// Check if app or any group has permission for this path
		if utils.Contains(policyPath.Apps, identity.GetUsername()) {
			if utils.Contains(policyPath.DenyPermissions, permission) {
				return false
			}

			if utils.Contains(policyPath.Permissions, permission) {
				return true
			}
		}
	}
	return false
}

// Helper to check if any secret paths explicitly deny the permission
func (h *Handler) checkSecretPaths(secrets []models.PolicySecret, pathID string, identity utils.Identity, userGroups []string, key string, permission string) bool {

	for _, p := range secrets {
		// Check is user is given explicit deny on the secret
		if p.Name == key && utils.Contains(p.DenyPermissions, permission) && utils.Contains(*p.DenyUsers, identity.GetUsername()) {
			return false
		}

		// Check is user group is given explicit deny on the secret
		if p.Name == key && utils.Contains(p.DenyPermissions, permission) && h.hasGroupAccess(*p.DenyGroups, userGroups) {
			return false
		}

		//TODO: Check for apps and certificates
	}
	return true
}

// Helper to check if any group in the user's group list matches the allowed groups
func (h *Handler) hasGroupAccess(allowedGroups, userGroups []string) bool {

	for _, group := range allowedGroups {
		if utils.Contains(userGroups, group) {
			return true
		}
	}
	return false
}
