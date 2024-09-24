package routes

import (
	"log"
	"net/http"
	"backendtku/app/controllers"
	"backendtku/app/middleware"
	"backendtku/config"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitializeRoutes(router *mux.Router, db *gorm.DB) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	middleware.InitRedis()
	middleware.InitMidtrans()
	handler := controllers.NewHandler(db, cfg)
	router.HandleFunc("/", handler.Home).Methods("GET")

	//auth
	router.HandleFunc("/register", handler.Register).Methods("POST")
    router.HandleFunc("/verify", handler.VerifyEmail).Methods("GET")
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/refresh-token", handler.RefreshToken).Methods("POST")
	router.Handle("/change-pass", middleware.Authenticate(http.HandlerFunc(handler.ChangePassword))).Methods("POST")
	router.Handle("/logout", middleware.Authenticate(http.HandlerFunc(handler.Logout))).Methods("POST")

	//product safari
    router.HandleFunc("/products", handler.GetProducts).Methods("GET")
	router.HandleFunc("/detail", handler.GetProductDetail).Methods("GET")
	router.HandleFunc("/countries", handler.GetCountries).Methods("GET")
	router.HandleFunc("/get-price", handler.GetSafariPrice).Methods("POST")
	router.HandleFunc("/prod-daymax", handler.GetDayMax).Methods("POST")

	//product abror
	router.HandleFunc("/get-abror", handler.GetAbrorDetail).Methods("GET")
	router.HandleFunc("/get-price-abror", handler.GetAbrorPrice).Methods("POST")
	router.HandleFunc("/get-cars", handler.GetCars).Methods("GET")

	//profile
	router.Handle("/profile", middleware.Authenticate(http.HandlerFunc(handler.GetUserProfile))).Methods("GET")
	router.Handle("/profile", middleware.Authenticate(http.HandlerFunc(handler.UpdateUserProfile))).Methods("PUT")

	//request
	router.Handle("/req-prod", middleware.Authenticate(http.HandlerFunc(handler.RequestProduct))).Methods("POST")
	router.Handle("/generate-trx", middleware.Authenticate(http.HandlerFunc(handler.CreateTransaction))).Methods("POST")
	router.Handle("/enroll-status", middleware.Authenticate(http.HandlerFunc(handler.PaymentStatus))).Methods("POST")
	router.Handle("/enroll-status-abror", middleware.Authenticate(http.HandlerFunc(handler.PaymentStatusAbror))).Methods("POST")
	router.Handle("/get-trx", middleware.Authenticate(http.HandlerFunc(handler.GetTrx))).Methods("POST")
	router.Handle("/req-abror", middleware.Authenticate(http.HandlerFunc(handler.RequestAbror))).Methods("POST")

	//policy
	router.Handle("/get-policies", middleware.Authenticate(http.HandlerFunc(handler.GetPolicies))).Methods("POST")
	router.Handle("/get-policies-abror", middleware.Authenticate(http.HandlerFunc(handler.GetPoliciesAbror))).Methods("POST")
	router.Handle("/download-pdf", middleware.Authenticate(http.HandlerFunc(handler.DownloadPdf))).Methods("POST")

	//claim 
	router.Handle("/req-claim", middleware.Authenticate(http.HandlerFunc(handler.RequestClaim))).Methods("POST")
	router.Handle("/get-claim", middleware.Authenticate(http.HandlerFunc(handler.GetClaim))).Methods("POST")
	router.Handle("/req-claim-abror", middleware.Authenticate(http.HandlerFunc(handler.RequestClaimAbror))).Methods("POST")
	router.Handle("/get-claim-abror", middleware.Authenticate(http.HandlerFunc(handler.GetClaimAbror))).Methods("POST")
	router.Handle("/get-claim-detail", middleware.Authenticate(http.HandlerFunc(handler.GetClaimDetail))).Methods("POST")
	

	//media
	router.PathPrefix("/upload/product/").Handler(http.StripPrefix("/upload/product/", http.FileServer(http.Dir("./upload/product"))))
	router.PathPrefix("/upload/users/").Handler(http.StripPrefix("/upload/users/", http.FileServer(http.Dir("./upload/users"))))
	router.PathPrefix("/upload/claim/").Handler(http.StripPrefix("/upload/claim/", http.FileServer(http.Dir("./upload/claim"))))
	router.PathPrefix("/upload/policy/pdfs/").Handler(http.StripPrefix("/upload/policy/pdfs/", http.FileServer(http.Dir("./upload/policy/pdfs"))))

	//admin
	router.Handle("/get-claim-sfr", middleware.Authenticate(http.HandlerFunc(handler.GetClaimSafariAll))).Methods("POST")
	router.Handle("/get-claim-abr", middleware.Authenticate(http.HandlerFunc(handler.GetClaimAbrorAll))).Methods("POST")
	router.Handle("/update-claim/{type}", middleware.Authenticate(http.HandlerFunc(handler.UpdateClaim))).Methods("PUT")
}	