package handlers

import (
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"net/http"
)

const (
	MasterDBContextKey = "masterDB"
	UserIDContextKey   = "user_id"
	UserRoleContextKey = "role"
)

func MongoDB(r *http.Request) data.MasterDB {
	return r.Context().Value(MasterDBContextKey).(data.MasterDB)
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	return userID, ok
}

func GetUserRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(UserRoleContextKey).(string)
	return role, ok
}
