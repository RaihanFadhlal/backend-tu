package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backendtku/app/helpers"
	"backendtku/app/middleware"
	"backendtku/app/models"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

// safari
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		TrxId string `json:"trx_id"`
	}

	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	var fullname string

	if err := h.DB.Model(&models.User{}).Where("email = ?", email).Select("name").First(&fullname).Error; err != nil {
		Response.Status = false
		Response.Message = "User not found"
		helpers.ResponseJSON(w, http.StatusNotFound, Response)
		return
	}

	var trx models.Transaction
	var trxa models.TransactionAbror
	var err error

	createSnapRequest := func(orderID string, totalPrice int64, productCode string, productPrice int64, capacity int32, productName string) {
		req := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: totalPrice,
			},
			CreditCard: &snap.CreditCardDetails{
				Secure: true,
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: fullname,
				Email: email,
			},
			Items: &[]midtrans.ItemDetails{
				{
					ID:    productCode,
					Price: productPrice,
					Qty:   capacity,
					Name:  productName,
				},
			},
		}

		snapResp, err := middleware.SnapClient.CreateTransaction(req)
		if err != nil {
			log.Println("Error creating transaction:", err)
			http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(snapResp)
	}

	if Request.TrxId[3] == 'S' {
		err = h.DB.Where("registrant_id = ? AND transaction_id = ? AND status = 'Menunggu Pembayaran'", email, Request.TrxId).First(&trx).Error
		if err != nil {
			Response.Status = false
			Response.Message = "Transaction not found"
			helpers.ResponseJSON(w, http.StatusNotFound, Response)
			return
		}
		createSnapRequest(trx.TransactionId, int64(trx.TotalPrice), trx.ProductCode, int64(trx.ProductPrice), int32(trx.Capacity), trx.ProductName)
	} else {
		err = h.DB.Where("registrant_id = ? AND transaction_id = ? AND status = 'Menunggu Pembayaran'", email, Request.TrxId).First(&trxa).Error
		if err != nil {
			Response.Status = false
			Response.Message = "Transaction not found"
			helpers.ResponseJSON(w, http.StatusNotFound, Response)
			return
		}
		createSnapRequest(trxa.TransactionId, int64(trxa.TotalPrice), trxa.ProductCode, int64(trxa.ProductPrice), int32(trxa.Capacity), trxa.ProductName)
	}
}

func (h *Handler) DownloadPdf(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		PolicyId string `json:"policy_id"`
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var enroll models.EnrollmentSafari
	if err := h.DB.Where("registrant_id = ? AND policy_id = ?", email, Request.PolicyId).First(&enroll).Error; err != nil {
		var enrollAbror models.EnrollmentAbror
		if err := h.DB.Where("registrant_id = ? AND policy_id = ?", email, Request.PolicyId).First(&enrollAbror).Error; err != nil {
			http.Error(w, "Policy not found or unauthorized", http.StatusNotFound)
			return
		}
		enroll = models.EnrollmentSafari{
			PolicyId: enrollAbror.PolicyId,
		}
	}

	filePath := "./upload/policy/pdfs/" + enroll.PolicyId + ".pdf"

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+enroll.PolicyId+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, filePath)
}

func (h *Handler) RequestProduct(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		ProductCode  string `json:"product_code"`
		ProductName  string `json:"product_name"`
		Capacity     int    `json:"capacity"`
		ProductPrice int64  `json:"product_price"`
		Phone        string `json:"phone"`
		From         string `json:"from"`
		Destination  string `json:"destination"`
		DateStart    string `json:"date_start"`
		DateEnd      string `json:"date_end"`
		Contribution string `json:"contribution"`
		FullName     string `json:"fullname"`
		Birthdate    string `json:"birthdate"`
		Birthplace   string `json:"birthplace"`
		Gender       string `json:"gender"`
		Passport     string `json:"passport"`
		Others       []struct {
			Fullname  string `json:"fullname"`
			Birthdate string `json:"birthdate"`
		} `json:"others"`
	}

	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		TrxId   string `json:"trx_id"`
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", Request.DateStart)
	if err != nil {
		Response.Status = false
		Response.Message = "Invalid start date format"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	endDate, err := time.Parse("2006-01-02", Request.DateEnd)
	if err != nil {
		Response.Status = false
		Response.Message = "Invalid end date format"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	if totalDays <= 0 {
		Response.Status = false
		Response.Message = "End date must be after start date"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	var count int64
	h.DB.Model(&models.ProductSafari{}).Where("code = ? AND contribution = ? AND price = ? AND day_min <= ? AND day_max >= ?",
		Request.ProductCode, Request.Contribution, Request.ProductPrice, totalDays, totalDays).Count(&count)

	if count == 0 {
		Response.Status = false
		Response.Message = "No matching product found"
		helpers.ResponseJSON(w, http.StatusNotFound, Response)
		return
	}

	grossAmt := int64(Request.Capacity) * Request.ProductPrice

	var transactionId string
	for {
		transactionId = "T-" + Request.ProductCode + "-" + helpers.RandomString(5)
		var count int64
		if err := h.DB.Model(&models.Transaction{}).Where("transaction_id = ?", transactionId).Count(&count).Error; err == nil && count == 0 {
			break
		}
	}

	transaction := models.Transaction{
		ID:            uuid.New(),
		TransactionId: transactionId,
		RegistrantId:  email,
		ProductCode:   Request.ProductCode,
		ProductName:   Request.ProductName,
		ProductPrice:  int(Request.ProductPrice),
		Capacity:      int(Request.Capacity),
		TotalPrice:    int(grossAmt),
		Status:        "Menunggu Pembayaran",
		CreatedAt:     time.Now(),
		ExpiredAt:     time.Now().Add(24 * time.Hour),
	}

	if err := h.DB.Create(&transaction).Error; err != nil {
		Response.Status = false
		Response.Message = "Error saving transaction: " + err.Error()
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	var enrollmentId string
	for {
		enrollmentId = "E-" + Request.ProductCode + "-" + helpers.RandomString(5)
		var count int64
		if err := h.DB.Model(&models.EnrollmentSafari{}).Where("enrollment_id = ?", enrollmentId).Count(&count).Error; err == nil && count == 0 {
			break
		}
	}

	enrollment := models.EnrollmentSafari{
		ID:            uuid.New(),
		EnrollmentId:  enrollmentId,
		RegistrantId:  email,
		TransactionId: transactionId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Phone:         Request.Phone,
		ProductCode:   Request.ProductCode,
		ProductName:   Request.ProductName,
		From:          Request.From,
		Destination:   Request.Destination,
		DateStart:     Request.DateStart,
		DateEnd:       Request.DateEnd,
		Contribution:  Request.Contribution,
		Capacity:      Request.Capacity,
		Name:          Request.FullName,
		Birthdate:     Request.Birthdate,
		Birthplace:    Request.Birthplace,
		Gender:        Request.Gender,
		Passport:      Request.Passport,
	}

	if err := h.DB.Create(&enrollment).Error; err != nil {
		Response.Status = false
		Response.Message = "Error saving enrollment: " + err.Error()
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	for _, other := range Request.Others {
		var otherEnrollmentId string
		for {
			otherEnrollmentId = "E-" + Request.ProductCode + "-" + helpers.RandomString(5)
			var count int64
			if err := h.DB.Model(&models.EnrollmentSafari{}).Where("enrollment_id = ?", otherEnrollmentId).Count(&count).Error; err == nil && count == 0 {
				break
			}
		}

		otherEnrollment := models.EnrollmentSafari{
			ID:            uuid.New(),
			EnrollmentId:  otherEnrollmentId,
			RegistrantId:  email,
			TransactionId: transactionId,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			ProductCode:   Request.ProductCode,
			ProductName:   Request.ProductName,
			From:          Request.From,
			Destination:   Request.Destination,
			DateStart:     Request.DateStart,
			DateEnd:       Request.DateEnd,
			Contribution:  Request.Contribution,
			Capacity:      Request.Capacity,
			Name:          other.Fullname,
			Birthdate:     other.Birthdate,
		}

		if err := h.DB.Create(&otherEnrollment).Error; err != nil {
			Response.Status = false
			Response.Message = "Error saving additional enrollment: " + err.Error()
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
	}

	Response.Status = true
	Response.Message = "Enrollment successfull"
	Response.TrxId = transactionId
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) PaymentStatus(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		TrxId string `json:"trx_id"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request body"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	var trx models.Transaction
	if err := h.DB.Where("registrant_id = ? AND transaction_id = ?", email, Request.TrxId).First(&trx).Error; err != nil {
		Response.Status = false
		Response.Message = "Transaction not found"
		helpers.ResponseJSON(w, http.StatusNotFound, Response)
		return
	}

	midtransResponse, err := middleware.VerifyMidtransTrx(Request.TrxId)
	if err != nil {
		Response.Status = false
		Response.Message = "Failed to verify transaction status"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	status := midtransResponse.TransactionStatus
	log.Println("Midtrans Transaction Status:", status)

	if status == "pending" {
		Response.Status = true
		Response.Message = "Menunggu Pembayaran"
		helpers.ResponseJSON(w, http.StatusOK, Response)
		return
	}

	if status == "expire" || status == "deny" || status == "cancel" {
		trx.Status = "Gagal"
		if err := h.DB.Save(&trx).Error; err != nil {
			Response.Status = false
			Response.Message = "Failed to update transaction status"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
		Response.Status = true
		Response.Message = "Pembayaran Gagal"
		helpers.ResponseJSON(w, http.StatusOK, Response)
		return
	}

	if status == "settlement" {
		trx.Status = "Berhasil"
		if err := h.DB.Save(&trx).Error; err != nil {
			Response.Status = false
			Response.Message = "Failed to update transaction status"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}

		var enroll models.EnrollmentSafari
		if err := h.DB.Where("registrant_id = ? AND transaction_id = ? AND LENGTH(phone) > 0", email, Request.TrxId).First(&enroll).Error; err != nil {
			Response.Status = false
			Response.Message = "Transaction not found"
			helpers.ResponseJSON(w, http.StatusNotFound, Response)
			return
		}

		uniqueCode := trx.TransactionId
		policyId := "policy-" + uniqueCode[len(uniqueCode)-5:]

		group := enroll.ProductCode[:len(enroll.ProductCode)-2]

		var product models.ProductSafari
		if err := h.DB.Where("group_code = ?", group).First(&product).Error; err != nil {
			Response.Status = false
			Response.Message = "Product not found"
			helpers.ResponseJSON(w, http.StatusNotFound, Response)
			return
		}

		var benefits []models.ProductBenefitSafari
		if err := h.DB.Where("group_code = ?", product.GroupCode).Find(&benefits).Error; err != nil {
			Response.Status = false
			Response.Message = "Benefits not found"
			helpers.ResponseJSON(w, http.StatusNotFound, Response)
			return
		}

		var result []string
		for _, benefit := range benefits {
			var contributionField string
			switch product.Contribution {
			case "Basic":
				contributionField = benefit.Basic
			case "Gold":
				contributionField = benefit.Gold
			case "Platinum":
				contributionField = benefit.Platinum
			case "Titanium":
				contributionField = benefit.Titanium
			default:
				contributionField = ""
			}

			if contributionField != "" {
				result = append(result, fmt.Sprintf("%s : %s", benefit.Detail, contributionField))
			}
		}

		var others []models.EnrollmentSafari
		if err := h.DB.Where("registrant_id = ? AND transaction_id = ?", email, Request.TrxId).Find(&others).Error; err != nil {
			Response.Status = false
			Response.Message = "Transaction not found"
			helpers.ResponseJSON(w, http.StatusNotFound, Response)
			return
		}
		log.Println(others)

		m := helpers.GetMaroto(policyId, trx.ProductName, enroll.Name, enroll.Contribution, enroll.DateStart, enroll.DateEnd, trx.TotalPrice, result, others)
		document, err := m.Generate()
		if err != nil {
			log.Fatal(err.Error())
		}
		err = document.Save("upload/policy/pdfs/" + policyId + ".pdf")
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := h.DB.Model(&models.EnrollmentSafari{}).
			Where("registrant_id = ? AND transaction_id = ?", email, Request.TrxId).
			Update("policy_id", policyId).Error; err != nil {
			Response.Status = false
			Response.Message = "Failed to save policies"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}

		Response.Status = true
		Response.Message = "Pembayaran Berhasil"
		helpers.ResponseJSON(w, http.StatusOK, Response)
	}
}

func (h *Handler) GetPolicies(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		ProductName string `json:"product_name"`
		Destination string `json:"destination"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			Contribution string `json:"contribution"`
			Destination  string `json:"destination"`
			DateStart    string `json:"sdate"`
			DateEnd      string `json:"edate"`
			Image        string `json:"image"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request body"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)
	query := h.DB.Where("registrant_id = ? AND LENGTH(phone) > 0 AND LENGTH(policy_id) > 0", email)

	if Request.Destination != "" {
		query = query.Where("destination = ?", Request.Destination)
	}
	if Request.ProductName != "" {
		query = query.Where("product_name = ?", Request.ProductName)
	}

	var enrolls []models.EnrollmentSafari
	if err := query.Order("created_at DESC").Find(&enrolls).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving products"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	for _, enroll := range enrolls {
		var product models.ProductSafari
		query2 := h.DB.Where("name = ? AND LENGTH(image) > 0", enroll.ProductName)
		if err := query2.Find(&product).Error; err != nil {
			Response.Status = false
			Response.Message = "Error retrieving products"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
		Response.Data = append(Response.Data, struct {
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			Contribution string `json:"contribution"`
			Destination  string `json:"destination"`
			DateStart    string `json:"sdate"`
			DateEnd      string `json:"edate"`
			Image        string `json:"image"`
		}{
			PolicyId:     enroll.PolicyId,
			ProductName:  enroll.ProductName,
			Contribution: enroll.Contribution,
			Destination:  enroll.Destination,
			DateStart:    enroll.DateStart,
			DateEnd:      enroll.DateEnd,
			Image:        h.Config.BaseUrl + "/upload/product/" + product.Image,
		})
	}

	Response.Status = true
	Response.Message = "Products retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

// abror
func (h *Handler) RequestAbror(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		Contribution string `json:"contribution"`
		ProductCode  string `json:"product_code"`
		ProductName  string `json:"product_name"`
		CarBrand     string `json:"car_brand"`
		CarType      string `json:"car_type"`
		Year         string `json:"year"`
		DateStart    string `json:"date_start"`
		DateEnd      string `json:"date_end"`
		Price        int    `json:"price"`
		Plat         string `json:"plat"`
		Chassis      string `json:"chassis"`
		Engine       string `json:"engine"`
		Image1       string `json:"image1"`
		Image2       string `json:"image2"`
		Image3       string `json:"image3"`
		Image4       string `json:"image4"`
		Fullname     string `json:"fullname"`
		Birthdate    string `json:"birtdate"`
		Gender       string `json:"gender"`
		Phone        string `json:"phone"`
		IdUser       string `json:"id_user"`
	}

	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		TrxId   string `json:"trx_id"`
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)
	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, Request.DateStart)
	if err != nil {
		http.Error(w, "Invalid DateStart format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse(layout, Request.DateEnd)
	if err != nil {
		http.Error(w, "Invalid DateEnd format", http.StatusBadRequest)
		return
	}

	duration := endDate.Sub(startDate)
	if duration.Hours() > 366*24 {
		Response.Status = false
		Response.Message = "Date range exceeds 366 days"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	var vehicle models.Car
	if err := h.DB.Where("name = ?",
		Request.CarType).First(&vehicle).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving car"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	var vehicleType models.VehicleType
	if err := h.DB.Where("min <= ? AND max >= ?",
		vehicle.Price, vehicle.Price).First(&vehicleType).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving car type"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	plat := helpers.ExtractPlateCode(Request.Plat)

	var region models.Region
	query := `SELECT * FROM regions WHERE ? = ANY(string_to_array(plat, ','))`
	if err := h.DB.Raw(query, plat).Scan(&region).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving region"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	var product models.ProductAbror
	if err := h.DB.Where("type_name = ? AND region_code = ? AND vehicle_code = ? ",
		Request.Contribution, region.Code, vehicleType.Code).First(&product).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving product"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	price := product.Percentage / 100 * float32(vehicle.Price)

	var transactionId string
	for {
		transactionId = "T-" + Request.ProductCode + "-" + helpers.RandomString(5)
		var count int64
		if err := h.DB.Model(&models.Transaction{}).Where("transaction_id = ?", transactionId).Count(&count).Error; err == nil && count == 0 {
			break
		}
	}

	transaction := models.TransactionAbror{
		ID:            uuid.New(),
		TransactionId: transactionId,
		RegistrantId:  email,
		ProductCode:   Request.ProductCode,
		ProductName:   Request.ProductName,
		ProductPrice:  int(price),
		Capacity:      1,
		TotalPrice:    int(price),
		Status:        "Menunggu Pembayaran",
		CreatedAt:     time.Now(),
		ExpiredAt:     time.Now().Add(24 * time.Hour),
	}

	if err := h.DB.Create(&transaction).Error; err != nil {
		Response.Status = false
		Response.Message = "Error saving transaction: " + err.Error()
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	var enrollmentId string
	for {
		enrollmentId = "E-" + Request.ProductCode + "-" + helpers.RandomString(5)
		var count int64
		if err := h.DB.Model(&models.EnrollmentSafari{}).Where("enrollment_id = ?", enrollmentId).Count(&count).Error; err == nil && count == 0 {
			break
		}
	}

	basePath := "upload/enroll"
	imageNames := []string{}
	for _, image := range []struct {
		Base64   string
		BaseName string
	}{
		{Request.Image1, enrollmentId},
		{Request.Image2, enrollmentId},
		{Request.Image3, enrollmentId},
		{Request.Image4, enrollmentId},
		{Request.IdUser, enrollmentId},
	} {
		imageName, err := saveImage(image.Base64, image.BaseName, basePath, h.DB)
		if err != nil {
			Response.Status = false
			Response.Message = err.Error()
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
		imageNames = append(imageNames, imageName)
	}

	enrollment := models.EnrollmentAbror{
		ID:            uuid.New(),
		EnrollmentId:  enrollmentId,
		RegistrantId:  email,
		TransactionId: transactionId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Phone:         Request.Phone,
		ProductCode:   Request.ProductCode,
		ProductName:   Request.ProductName,
		DateStart:     Request.DateStart,
		DateEnd:       Request.DateEnd,
		Contribution:  Request.Contribution,
		CarBrand:      Request.CarBrand,
		CarType:       Request.CarType,
		Year:          Request.Year,
		Plat:          Request.Plat,
		Chassis:       Request.Chassis,
		Engine:        Request.Engine,
		Image1:        imageNames[0],
		Image2:        imageNames[1],
		Image3:        imageNames[2],
		Image4:        imageNames[3],
		Name:          Request.Fullname,
		Birthdate:     Request.Birthdate,
		Gender:        Request.Gender,
		IdUser:        imageNames[4],
	}

	if err := h.DB.Create(&enrollment).Error; err != nil {
		Response.Status = false
		Response.Message = "Error saving enrollment: " + err.Error()
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	Response.Status = true
	Response.Message = "Enrollment successfull"
	Response.TrxId = transactionId
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) PaymentStatusAbror(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		TrxId string `json:"trx_id"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request body"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	var trx models.TransactionAbror
	if err := h.DB.Where("registrant_id = ? AND transaction_id = ?", email, Request.TrxId).First(&trx).Error; err != nil {
		Response.Status = false
		Response.Message = "Transaction not found"
		helpers.ResponseJSON(w, http.StatusNotFound, Response)
		return
	}

	midtransResponse, err := middleware.VerifyMidtransTrx(Request.TrxId)
	if err != nil {
		Response.Status = false
		Response.Message = "Failed to verify transaction status"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	status := midtransResponse.TransactionStatus
	log.Println("Midtrans Transaction Status:", status)

	if status == "pending" {
		Response.Status = true
		Response.Message = "Menunggu Pembayaran"
		helpers.ResponseJSON(w, http.StatusOK, Response)
		return
	}

	if status == "expire" || status == "deny" || status == "cancel" {
		trx.Status = "Gagal"
		if err := h.DB.Save(&trx).Error; err != nil {
			Response.Status = false
			Response.Message = "Failed to update transaction status"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
		Response.Status = true
		Response.Message = "Pembayaran Gagal"
		helpers.ResponseJSON(w, http.StatusOK, Response)
		return
	}

	if status == "settlement" {
		trx.Status = "Berhasil"
		if err := h.DB.Save(&trx).Error; err != nil {
			Response.Status = false
			Response.Message = "Failed to update transaction status"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}

		var enroll models.EnrollmentAbror
		if err := h.DB.Where("registrant_id = ? AND transaction_id = ? AND LENGTH(phone) > 0", email, Request.TrxId).First(&enroll).Error; err != nil {
			Response.Status = false
			Response.Message = "Transaction not found"
			helpers.ResponseJSON(w, http.StatusNotFound, Response)
			return
		}

		uniqueCode := trx.TransactionId
		policyId := "policy-" + uniqueCode[len(uniqueCode)-5:]

		char6 := Request.TrxId[5]
		var benefits []models.ProductBenefitAbror
		if err := h.DB.Order("id ASC").Find(&benefits).Error; err != nil {
			http.Error(w, "Error retrieving benefits", http.StatusInternalServerError)
			return
		}

		var result []string
		for _, benefit := range benefits {
			var value string
			switch char6 {
			case 'S':
				value = benefit.Standard
			case 'P':
				value = benefit.Premium
			default:
				continue
			}

			if value != "" {
				result = append(result, fmt.Sprintf("%s : %s", benefit.Description, value))
			}
		}

		m := helpers.GetMarotoAbror(policyId, enroll.CarBrand, enroll.CarType, enroll.Plat, enroll.Name, enroll.Contribution, enroll.DateStart, enroll.DateEnd, trx.TotalPrice, result, enroll.Chassis, enroll.Engine, enroll.Image1, enroll.Image2, enroll.Image3, enroll.Image4)
		document, err := m.Generate()
		if err != nil {
			log.Fatal(err.Error())
		}
		err = document.Save("upload/policy/pdfs/" + policyId + ".pdf")
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := h.DB.Model(&models.EnrollmentAbror{}).
			Where("registrant_id = ? AND transaction_id = ?", email, Request.TrxId).
			Update("policy_id", policyId).Error; err != nil {
			Response.Status = false
			Response.Message = "Failed to save policies"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}

		Response.Status = true
		Response.Message = "Pembayaran Berhasil"
		helpers.ResponseJSON(w, http.StatusOK, Response)
	}
}

func (h *Handler) GetPoliciesAbror(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		CarType   string `json:"car_type"`
		DateStart string `json:"date_start"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			Contribution string `json:"contribution"`
			CarType      string `json:"car_type"`
			DateStart    string `json:"date_start"`
			DateEnd      string `json:"date_end"`
			Image        string `json:"image"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request body"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)
	query := h.DB.Where("registrant_id = ? AND LENGTH(phone) > 0 AND LENGTH(policy_id) > 0", email)

	if Request.CarType != "" {
		query = query.Where("car_type = ?", Request.CarType)
	}
	if Request.DateStart != "" {
		query = query.Where("date_start = ?", Request.DateStart)
	}

	var enrolls []models.EnrollmentAbror
	if err := query.Order("created_at DESC").Find(&enrolls).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving products"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	for _, enroll := range enrolls {
		var product models.ProductAbror
		query2 := h.DB.Where("name = ? AND LENGTH(image) > 0", enroll.ProductName)
		if err := query2.Find(&product).Error; err != nil {
			Response.Status = false
			Response.Message = "Error retrieving products"
			helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
			return
		}
		Response.Data = append(Response.Data, struct {
			PolicyId     string `json:"policy_id"`
			ProductName  string `json:"product_name"`
			Contribution string `json:"contribution"`
			CarType      string `json:"car_type"`
			DateStart    string `json:"date_start"`
			DateEnd      string `json:"date_end"`
			Image        string `json:"image"`
		}{
			PolicyId:     enroll.PolicyId,
			ProductName:  enroll.ProductName,
			Contribution: enroll.Contribution,
			CarType:      enroll.CarType,
			DateStart:    enroll.DateStart,
			DateEnd:      enroll.DateEnd,
			Image:        h.Config.BaseUrl + "/upload/product/" + product.Image,
		})
	}

	Response.Status = true
	Response.Message = "Products retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) GetTrx(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		ProductName string `json:"product_name"`
		Status      string `json:"status"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			TransactionId string `json:"transaction_id"`
			ProductName   string `json:"product_name"`
			TotalPrice    int    `json:"total_price"`
			CreatedAt     string `json:"created_at"`
			ExpiredAt     string `json:"expired_at"`
			Status        string `json:"status"`
			Image         string `json:"image"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request body"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}

	email := r.Context().Value(middleware.UserEmailKey).(string)

	handleTransactions := func(query *gorm.DB, productTable string, isSafari bool) error {
		var trxs []models.Transaction
		if err := query.Order("created_at DESC").Find(&trxs).Error; err != nil {
			return err
		}

		for _, trx := range trxs {
			var product interface{}
			if isSafari {
				product = &models.ProductSafari{}
			} else {
				product = &models.ProductAbror{}
			}

			productQuery := h.DB.Table(productTable).Where("name = ? AND LENGTH(image) > 0", trx.ProductName)
			if err := productQuery.Find(product).Error; err != nil {
				return err
			}

			var image string
			if isSafari {
				image = product.(*models.ProductSafari).Image
			} else {
				image = product.(*models.ProductAbror).Image
			}

			Response.Data = append(Response.Data, struct {
				TransactionId string `json:"transaction_id"`
				ProductName   string `json:"product_name"`
				TotalPrice    int    `json:"total_price"`
				CreatedAt     string `json:"created_at"`
				ExpiredAt     string `json:"expired_at"`
				Status        string `json:"status"`
				Image         string `json:"image"`
			}{
				TransactionId: trx.TransactionId,
				ProductName:   trx.ProductName,
				TotalPrice:    trx.TotalPrice,
				CreatedAt:     trx.CreatedAt.Format("2006-01-02 15:04:05"),
				ExpiredAt:     trx.ExpiredAt.Format("2006-01-02 15:04:05"),
				Status:        trx.Status,
				Image:         h.Config.BaseUrl + "/upload/product/" + image,
			})
		}
		return nil
	}

	query1 := h.DB.Where("registrant_id = ?", email)
	if Request.Status != "" {
		query1 = query1.Where("status = ?", Request.Status)
	}
	if Request.ProductName != "" {
		query1 = query1.Where("product_name = ?", Request.ProductName)
	}
	if err := handleTransactions(query1.Model(&models.Transaction{}), "product_safaris", true); err != nil {
		Response.Status = false
		Response.Message = "Error retrieving transactions"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	query2 := h.DB.Where("registrant_id = ?", email)
	if Request.Status != "" {
		query2 = query2.Where("status = ?", Request.Status)
	}
	if Request.ProductName != "" {
		query2 = query2.Where("product_name = ?", Request.ProductName)
	}
	if err := handleTransactions(query2.Model(&models.TransactionAbror{}), "product_abrors", false); err != nil {
		Response.Status = false
		Response.Message = "Error retrieving transactions2"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	Response.Status = true
	Response.Message = "Products retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func saveImage(imageBase64, baseName, path string, db *gorm.DB) (string, error) {
	if imageBase64 == "" {
		return "", nil
	}

	imageFormat := helpers.GetTypeBase64(imageBase64)
	timestamp := time.Now().UnixNano()
	imageName := fmt.Sprintf("%s_%d%s", baseName, timestamp, imageFormat)

	decodedImage, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	imagePath := filepath.Join(path, imageName)
	if err := os.WriteFile(imagePath, decodedImage, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return imageName, nil
}
