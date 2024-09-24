package controllers

import (
	"net/http"
)
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Takaful Umum Back End"))
}