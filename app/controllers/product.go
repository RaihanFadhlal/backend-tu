package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"backendtku/app/helpers"
	"backendtku/app/models"
)

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Price int    `json:"price"`
			Image string `json:"image"`
		} `json:"data"`
	}

	country := r.URL.Query().Get("country")

	var products []models.ProductSafari
	query := h.DB.Where("code LIKE ?", "%B1")

	if country != "" {
		countryFilter := "%" + country + "%"
		query = query.Where("countries LIKE ?", countryFilter)
	}

	if err := query.Find(&products).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving products"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	for _, product := range products {
		Response.Data = append(Response.Data, struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Price int    `json:"price"`
			Image string `json:"image"`
		}{
			Code:  product.Code,
			Name:  product.Name,
			Price: product.Price,
			Image: h.Config.BaseUrl + "/upload/product/" + product.Image,
		})
	}

	Response.Status = true
	Response.Message = "Products retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	type PricePeriod struct {
		DayMin int `json:"day_min"`
		DayMax int `json:"day_max"`
		Price  int `json:"price"`
	}

	type PriceDetails struct {
		Basic    []PricePeriod `json:"basic"`
		Gold     []PricePeriod `json:"gold"`
		Platinum []PricePeriod `json:"platinum"`
		Titanium []PricePeriod `json:"titanium"`
	}

	type BenefitDetail struct {
		Detail   string `json:"detail"`
		Basic    string `json:"basic"`
		Gold     string `json:"gold"`
		Platinum string `json:"platinum"`
		Titanium string `json:"titanium"`
	}

	type BenefitCategory struct {
		Desc   string          `json:"desc"`
		Detail []BenefitDetail `json:"detail"`
	}

	type Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Name        string            `json:"name"`
			Description string            `json:"desc"`
			Image       string            `json:"image"`
			Terms       string            `json:"tnc"`
			Types       []string          `json:"type"`
			Countries   string            `json:"countries"`
			Price       PriceDetails      `json:"price"`
			Benefits    []BenefitCategory `json:"benefits"`
		} `json:"data"`
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		helpers.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": "Product ID is required"})
		return
	}

	var product models.ProductSafari
	if err := h.DB.Where("code = ?", id).First(&product).Error; err != nil {
		helpers.ResponseJSON(w, http.StatusNotFound, map[string]string{"message": "Product not found"})
		return
	}

	groupCode := id[:len(id)-2]

	var contributions []string

	if err := h.DB.Model(&models.ProductSafari{}).
		Distinct("contribution").
		Where("group_code = ?", groupCode).
		Pluck("contribution", &contributions).Error; err != nil {
		http.Error(w, "Failed to retrieve contributions", http.StatusInternalServerError)
		return
	}

	var priceDetails PriceDetails
	priceCategories := []string{"Basic", "Gold", "Platinum", "Titanium"}
	for _, category := range priceCategories {
		var prices []PricePeriod
		h.DB.Table("product_safaris").Select("day_min, day_max, price").
			Where("group_code = ? AND contribution = ?", groupCode, category).
			Order("day_min ASC").
			Scan(&prices)

		switch category {
		case "Basic":
			priceDetails.Basic = prices
		case "Gold":
			priceDetails.Gold = prices
		case "Platinum":
			priceDetails.Platinum = prices
		case "Titanium":
			priceDetails.Titanium = prices
		}
	}

	var benefits []BenefitCategory
	var distinctDesc []string
	h.DB.Table("product_benefit_safaris").Select("distinct COALESCE(description, '-') as desc").Where("group_code = ?", groupCode).Scan(&distinctDesc)
	for _, desc := range distinctDesc {
		var benefitDetails []BenefitDetail
		h.DB.Table("product_benefit_safaris").Select("COALESCE(detail, '-') as detail, COALESCE(basic, '-') as basic, COALESCE(gold, '-') as gold, COALESCE(platinum, '-') as platinum, COALESCE(titanium, '-') as titanium").
			Where("group_code = ? AND description = ?", groupCode, desc).
			Scan(&benefitDetails)

		benefits = append(benefits, BenefitCategory{
			Desc:   desc,
			Detail: benefitDetails,
		})
	}

	response := Response{
		Status:  true,
		Message: "Product details retrieved successfully",
	}
	response.Data.Name = product.Name
	response.Data.Description = product.Description
	response.Data.Image = h.Config.BaseUrl + "/upload/product/" + product.Image
	response.Data.Terms = product.Terms
	response.Data.Countries = product.Countries
	response.Data.Price = priceDetails
	response.Data.Benefits = benefits
	response.Data.Types = contributions

	helpers.ResponseJSON(w, http.StatusOK, response)
}

func (h *Handler) GetAbrorDetail(w http.ResponseWriter, r *http.Request) {
	type PricePeriod struct {
		C          string  `json:"c"`
		RangePrice string  `json:"range_price"`
		A1         float32 `json:"a1"`
		A2         float32 `json:"a2"`
		A3         float32 `json:"a3"`
	}

	type PriceDetails struct {
		Standard []PricePeriod `json:"standard"`
		Premium  []PricePeriod `json:"premium"`
	}

	type BenefitDetail struct {
		Standard string `json:"standard"`
		Premium  string `json:"premium"`
	}

	type BenefitCategory struct {
		Desc   string          `json:"desc"`
		Type   string          `json:"type"`
		Detail []BenefitDetail `json:"detail"`
	}

	type Cars struct {
		Brand string   `json:"brand"`
		Type  []string `json:"type"`
	}

	type Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Name        string            `json:"name"`
			Description string            `json:"desc"`
			Image       string            `json:"image"`
			Terms       string            `json:"terms"`
			Cars        []Cars            `json:"cars"`
			Price       PriceDetails      `json:"price"`
			Benefits    []BenefitCategory `json:"benefits"`
		} `json:"data"`
	}
	var response Response
	response.Status = true
	response.Message = "Success"

	var product models.ProductAbror
	if err := h.DB.Where("id = '1'").Find(&product).Error; err != nil {
		response.Status = false
		response.Message = "Error retrieving products"
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response.Data.Name = product.Name
	response.Data.Description = product.Description
	response.Data.Image = h.Config.BaseUrl + "/upload/product/" + product.Image
	response.Data.Terms = product.AllowedVehicle

	var productAbror models.ProductAbror
	if err := h.DB.Find(&productAbror).Error; err != nil {
		response.Status = false
		response.Message = "Error fetching vehicle types"
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var vehicleTypes []models.VehicleType
	if err := h.DB.Find(&vehicleTypes).Error; err != nil {
		response.Status = false
		response.Message = "Error fetching vehicle types"
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	fetchAndMapProductAbrors := func(typeName string) (map[string]map[string]float32, error) {
		var productAbrors []models.ProductAbror
		if err := h.DB.Where("type_name = ?", typeName).Find(&productAbrors).Error; err != nil {
			return nil, err
		}
		productAbrorMap := make(map[string]map[string]float32)
		for _, pa := range productAbrors {
			if _, exists := productAbrorMap[pa.VehicleCode]; !exists {
				productAbrorMap[pa.VehicleCode] = make(map[string]float32)
			}
			productAbrorMap[pa.VehicleCode][pa.RegionCode] = pa.Percentage
		}
		return productAbrorMap, nil
	}

	populatePriceDetails := func(typeName string, productAbrorMap map[string]map[string]float32) []PricePeriod {
		var pricePeriods []PricePeriod
		for _, vt := range vehicleTypes {
			pp := PricePeriod{
				C:          vt.Code,
				RangePrice: fmt.Sprintf("(%d-%d)", vt.Min, vt.Max),
				A1:         productAbrorMap[vt.Code]["A1"],
				A2:         productAbrorMap[vt.Code]["A2"],
				A3:         productAbrorMap[vt.Code]["A3"],
			}
			pricePeriods = append(pricePeriods, pp)
		}
		return pricePeriods
	}

	standardProductAbrorMap, err := fetchAndMapProductAbrors("Standard")
	if err != nil {
		response.Status = false
		response.Message = "Error fetching product abrors for Standard"
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	response.Data.Price.Standard = populatePriceDetails("Standard", standardProductAbrorMap)

	premiumProductAbrorMap, err := fetchAndMapProductAbrors("Premium")
	if err != nil {
		response.Status = false
		response.Message = "Error fetching product abrors for Premium"
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	response.Data.Price.Premium = populatePriceDetails("Premium", premiumProductAbrorMap)

	var brands []string
	if err := h.DB.Model(&models.Car{}).Distinct().Pluck("brand", &brands).Error; err != nil {
		response.Status = false
		response.Message = "Error fetching car brands"
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var cars []Cars
	for _, brand := range brands {
		var types []string
		if err := h.DB.Model(&models.Car{}).Where("brand = ?", brand).Pluck("name", &types).Error; err != nil {
			response.Status = false
			response.Message = "Error fetching car types for brand " + brand
			helpers.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
		cars = append(cars, Cars{
			Brand: brand,
			Type:  types,
		})
	}
	response.Data.Cars = cars

	var benefitCategories []BenefitCategory

	var benefits []models.ProductBenefitAbror
	if err := h.DB.Find(&benefits).Error; err != nil {
		response.Status = false
		response.Message = "Error fetching benefits"
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	benefitMap := make(map[string]BenefitCategory)

	for _, benefit := range benefits {
		key := benefit.Description + "_" + benefit.Type 
		if category, exists := benefitMap[key]; exists {
			benefitDetail := BenefitDetail{
				Standard: benefit.Standard,
				Premium:  benefit.Premium,
			}
			category.Detail = append(category.Detail, benefitDetail)
			benefitMap[key] = category
		} else {
			benefitDetail := BenefitDetail{
				Standard: benefit.Standard,
				Premium:  benefit.Premium,
			}
			benefitCategory := BenefitCategory{
				Desc:   benefit.Description,
				Type:   benefit.Type,
				Detail: []BenefitDetail{benefitDetail},
			}
			benefitMap[key] = benefitCategory
		}
	}

	for _, category := range benefitMap {
		benefitCategories = append(benefitCategories, category)
	}

	response.Data.Benefits = benefitCategories
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func (h *Handler) GetAbrorPrice(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		Contribution string `json:"contribution"`
		Type         string `json:"type"`
		PlatCode     string  `json:"plat_code"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ProductCode  string `json:"product_code"`
			Price int    `json:"price"`
			Percentage float32 `json:"percentage"`
			VehiclePrice int `json:"vehicle_price"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request payload"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}
	defer r.Body.Close()

	var vehicle models.Car
	if err := h.DB.Where("name = ?",
		Request.Type,).First(&vehicle).Error; err != nil {
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

	plat := helpers.ExtractPlateCode(Request.PlatCode)
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

	price := product.Percentage/100 * float32(vehicle.Price)

	Response.Data.ProductCode = product.Code
	Response.Data.Price = int(price)
	Response.Data.Percentage = product.Percentage
	Response.Data.VehiclePrice = vehicle.Price
	Response.Status = true
	Response.Message = "Product retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) GetDayMax(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		GroupCode string `json:"group_code"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			DayMax int `json:"day_max"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request payload"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}
	defer r.Body.Close()

	var price models.ProductSafari
	if err := h.DB.Where("group_code = ?",
		Request.GroupCode).Order("day_max desc").First(&price).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving product"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	Response.Data.DayMax = price.DayMax
	Response.Status = true
	Response.Message = "Product retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) GetCountries(w http.ResponseWriter, r *http.Request) {
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Countries string `json:"countries"`
		} `json:"data"`
	}

	var result struct {
		Countries string
	}

	err := h.DB.Model(&models.ProductSafari{}).
		Select("countries").
		Order("LENGTH(countries) DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		Response.Status = false
		Response.Message = "Error retrieving country"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	Response.Data.Countries = result.Countries
	Response.Status = true
	Response.Message = "Country retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) GetCars(w http.ResponseWriter, r *http.Request) {
	var Response struct {
		Status  bool     `json:"status"`
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}

	var cars []models.Car
	if err := h.DB.Select("name").Find(&cars).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving car names"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	var carNames []string
	for _, car := range cars {
		carNames = append(carNames, car.Name)
	}

	Response.Status = true
	Response.Message = "Car names retrieved successfully"
	Response.Data = carNames
	helpers.ResponseJSON(w, http.StatusOK, Response)
}

func (h *Handler) GetSafariPrice(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		GroupCode string `json:"group_code"`
		Type      string `json:"type"`
		Period    int    `json:"period"`
	}
	var Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Code  string `json:"code"`
			Price int    `json:"price"`
			Name  string `json:"name"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Request); err != nil {
		Response.Status = false
		Response.Message = "Invalid request payload"
		helpers.ResponseJSON(w, http.StatusBadRequest, Response)
		return
	}
	defer r.Body.Close()
	var product models.ProductSafari
	if err := h.DB.Where("group_code = ? AND contribution = ? AND day_min <= ? AND day_max >= ?",
		Request.GroupCode, Request.Type, Request.Period, Request.Period).First(&product).Error; err != nil {
		Response.Status = false
		Response.Message = "Error retrieving product"
		helpers.ResponseJSON(w, http.StatusInternalServerError, Response)
		return
	}

	Response.Data.Code = product.Code
	Response.Data.Price = product.Price
	Response.Data.Name = product.Name
	Response.Status = true
	Response.Message = "Product retrieved successfully"
	helpers.ResponseJSON(w, http.StatusOK, Response)
}