package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"backendtku/app/helpers"
	"backendtku/app/middleware"
	"backendtku/app/models"

	"github.com/gorilla/mux"
)

func (h *Handler) GetClaimSafariAll(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		PolicyId     string `json:"policy_id"`
		ProductName  string `json:"product_name"`
		DateReport   string `json:"date_report"`
		RegistrantId string `json:"registrant_id"`
		Status       string `json:"status"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			ClaimId      string `json:"claim_id"`
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			DateReport   string `json:"date_report"`
			DateAccident string `json:"date_accident"`
			Status       string `json:"status"`
			Image        string `json:"image"`
			Evidence     string `json:"evidence"`
			Detail       string `json:"detail"`
			PolicyPdf    string `json:"policy_pdf"`
			RegistrantId string `json:"registrant_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request body"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		Response.Status = false
		Response.Message = "Invalid email or password"
		helpers.ResponseJSON(w, http.StatusUnauthorized, Response)
		return
	}

	if user.Type != "admin" {
		Response.Status = false
		Response.Message = "unauthorized"
		helpers.ResponseJSON(w, http.StatusUnauthorized, Response)
		return
	}

	query := h.DB.Model(&models.ClaimSafari{})

	if Request.PolicyId != "" {
		query = h.DB.Where("policy_id = ?", Request.PolicyId)
	}
	if Request.ProductName != "" {
		query = query.Where("product_name = ?", Request.ProductName)
	}
	if Request.DateReport != "" {
		query = query.Where("date_report = ?", Request.DateReport)
	}
	if Request.RegistrantId != "" {
		query = query.Where("registrant_id = ?", Request.RegistrantId)
	}
	if Request.Status != "" {
		query = query.Where("status = ?", Request.Status)
	}

	var claims []models.ClaimSafari
	if err := query.Order("created_at DESC").Find(&claims).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving products"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	for _, claim := range claims {
		var product models.ProductSafari
		query2 := h.DB.Where("name = ? AND LENGTH(image) > 0", claim.ProductName)
		if err := query2.Find(&product).Error; err != nil {
			Response.Status = false
			Response.Message = "Error retrieving products"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
		Response.Data = append(Response.Data, struct {
			ClaimId      string `json:"claim_id"`
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			DateReport   string `json:"date_report"`
			DateAccident string `json:"date_accident"`
			Status       string `json:"status"`
			Image        string `json:"image"`
			Evidence     string `json:"evidence"`
			Detail       string `json:"detail"`
			PolicyPdf    string `json:"policy_pdf"`
			RegistrantId string `json:"registrant_id"`
		}{
			ClaimId:      claim.ClaimId,
			PolicyId:     claim.PolicyId,
			ProductName:  claim.ProductName,
			DateReport:   claim.DateReport,
			DateAccident: claim.DateAccident,
			Status:       claim.Status,
			Image:        h.Config.BaseUrl + "/upload/product/" + product.Image,
			Evidence:     h.Config.BaseUrl + "/upload/claim/" + claim.Evidence,
			Detail:       claim.Detail,
			PolicyPdf:    h.Config.BaseUrl + "/upload/policy/pdfs/" + claim.PolicyId + ".pdf",
			RegistrantId: claim.RegistrantId,
		})
	}

	Response.Status = true
	Response.Message = "Products retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) GetClaimAbrorAll(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		PolicyId   string `json:"policy_id"`
		DateReport string `json:"date_report"`
		RegistrantId string `json:"registrant_id"`
		Status     string `json:"status"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			ClaimId      string `json:"claim_id"`
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			DateReport   string `json:"date_report"`
			DateAccident string `json:"date_accident"`
			Status       string `json:"status"`
			Image        string `json:"image"`
			Evidence     string `json:"evidence"`
			Detail       string `json:"detail"`
			PolicyPdf    string `json:"policy_pdf"`
			RegistrantId string `json:"registrant_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request body"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		Response.Status = false
		Response.Message = "Invalid email or password"
		helpers.ResponseJSON(w, http.StatusUnauthorized, Response)
		return
	}

	if user.Type != "admin" {
		Response.Status = false
		Response.Message = "Invalid email or password"
		helpers.ResponseJSON(w, http.StatusUnauthorized, Response)
		return
	}

	query := h.DB.Model(&models.ClaimAbror{})

	if Request.PolicyId != "" {
		query = h.DB.Where("policy_id = ?", Request.PolicyId)
	}
	if Request.DateReport != "" {
		query = query.Where("date_report = ?", Request.DateReport)
	}
	if Request.RegistrantId != "" {
		query = query.Where("registrant_id = ?", Request.RegistrantId)
	}
	if Request.Status != "" {
		query = query.Where("date_report = ?", Request.Status)
	}

	var claims []models.ClaimAbror
	if err := query.Order("created_at DESC").Find(&claims).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving products"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	for _, claim := range claims {
		var product models.ProductAbror
		query2 := h.DB.Where("name = ? AND LENGTH(image) > 0", claim.ProductName)
		if err := query2.Find(&product).Error; err != nil {
			Response.Status = false
			Response.Message = "Error retrieving products"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
		Response.Data = append(Response.Data, struct {
			ClaimId      string `json:"claim_id"`
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			DateReport   string `json:"date_report"`
			DateAccident string `json:"date_accident"`
			Status       string `json:"status"`
			Image        string `json:"image"`
			Evidence     string `json:"evidence"`
			Detail       string `json:"detail"`
			PolicyPdf    string `json:"policy_pdf"`
			RegistrantId string `json:"registrant_id"`
		}{
			ClaimId:      claim.ClaimId,
			PolicyId:     claim.PolicyId,
			ProductName:  claim.ProductName,
			DateReport:   claim.DateReport,
			DateAccident: claim.DateAccident,
			Status:       claim.Status,
			Image:        h.Config.BaseUrl + "/upload/product/" + product.Image,
			Evidence:     h.Config.BaseUrl + "/upload/claim/" + claim.Evidence,
			Detail:       claim.Detail,
			PolicyPdf:    h.Config.BaseUrl + "/upload/policy/pdfs/" + claim.PolicyId + ".pdf",
			RegistrantId: claim.RegistrantId,
		})
	}

	Response.Status = true
	Response.Message = "Products retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) UpdateClaim(w http.ResponseWriter, r *http.Request) {
    claimType := mux.Vars(r)["type"]
	
	var Request struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		CoverCost int    `json:"cover_cost"`
		PayProof  string `json:"pay_proof"`
		ClaimId   string `json:"claim_id"`
	}

	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request payload"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		Response.Status = false
		Response.Message = "Invalid email or password"
		helpers.ResponseJSON(w, http.StatusUnauthorized, Response)
		return
	}

	if user.Type != "admin" {
		Response.Status = false
		Response.Message = "Unauthorized"
		helpers.ResponseJSON(w, http.StatusUnauthorized, Response)
		return
	}

	var claim interface{}

	switch claimType {
	case "safari":
		claim = &models.ClaimSafari{}
	case "abror":
		claim = &models.ClaimAbror{}
	default:
		Response.Status = false
		Response.Message = "Invalid claim type"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	if err := h.DB.Where("claim_id = ?", Request.ClaimId).First(claim).Error; err != nil {
		Response.Status = false
		Response.Message = "Claim not found"
		helpers.ResponseJSON(w, http.StatusNotFound, Response)
		return
	}

	switch c := claim.(type) {
	case *models.ClaimSafari:
		c.Status = Request.Status
		c.Message = Request.Message
		c.CoverCost = Request.CoverCost
	case *models.ClaimAbror:
		c.Status = Request.Status
		c.Message = Request.Message
		c.CoverCost = Request.CoverCost
	}

	if Request.PayProof != "" {
		imageFormat := helpers.GetTypeBase64(Request.PayProof)
		imageName := "ClaimProof-" + Request.ClaimId + imageFormat

		var count int64
		if claimType == "safari" {
			for {
				if err := h.DB.Model(&models.ClaimSafari{}).Where("pay_proof = ?", imageName).Count(&count).Error; err != nil {
					Response.Status = false
					Response.Message = "Error checking image name"
					helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
					return
				}
				if count == 0 {
					break
				}
				count++
				imageName = fmt.Sprintf("ClaimProof-%s%d%s", Request.ClaimId, count, imageFormat)
			}
		} else if claimType == "abror" {
			for {
				if err := h.DB.Model(&models.ClaimAbror{}).Where("pay_proof = ?", imageName).Count(&count).Error; err != nil {
					Response.Status = false
					Response.Message = "Error checking image name"
					helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
					return
				}
				if count == 0 {
					break
				}
				count++
				imageName = fmt.Sprintf("ClaimProof-%s%d%s", Request.ClaimId, count, imageFormat)
			}
		}

		decodedImage, err := base64.StdEncoding.DecodeString(Request.PayProof)
		if err != nil {
			Response.Status = false
			Response.Message = "Failed to decode image"
			helpers.ResponseJSON(w, http.StatusBadRequest, Response)
			return
		}

		imagePath := filepath.Join("upload/claim", imageName)
		if err := os.WriteFile(imagePath, decodedImage, 0644); err != nil {
			Response.Status = false
			Response.Message = "Failed to save image"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}

		switch c := claim.(type) {
		case *models.ClaimSafari:
			c.PayProof = imageName
		case *models.ClaimAbror:
			c.PayProof = imageName
		}
	}

	if err := h.DB.Save(claim).Error; err != nil {
		Response.Status = false
		Response.Message = "Failed to update claim"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	Response.Status = true
	Response.Message = "Claim updated successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}
