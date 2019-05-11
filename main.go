package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	"github.com/n1try/wakapi/middlewares"
	"github.com/n1try/wakapi/models"
	"github.com/n1try/wakapi/routes"
	"github.com/n1try/wakapi/services"
	"github.com/n1try/wakapi/utils"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func readConfig() models.Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(os.Getenv("WAKAPI_PORT"))
	dbUser, valid := os.LookupEnv("WAKAPI_DB_USER")
	dbPassword, valid := os.LookupEnv("WAKAPI_DB_PASSWORD")
	dbHost, valid := os.LookupEnv("WAKAPI_DB_HOST")
	dbName, valid := os.LookupEnv("WAKAPI_DB_NAME")

	if err != nil {
		log.Fatal(err)
	}
	if !valid {
		log.Fatal("Config parameters missing.")
	}

	return models.Config{
		Port:       port,
		DbHost:     dbHost,
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbName:     dbName,
		DbDialect:  "mysql",
	}
}

func main() {
	// Read Config
	config := readConfig()

	// Connect to database
	db, err := gorm.Open(config.DbDialect, utils.MakeConnectionString(&config))
	if err != nil {
		// log.Fatal("Could not connect to database.")
		log.Fatal(err)
	}
	defer db.Close()

	// Migrate database schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Heartbeat{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	// Services
	heartbeatSrvc := &services.HeartbeatService{db}
	userSrvc := &services.UserService{db}
	aggregationSrvc := &services.AggregationService{db, heartbeatSrvc}

	// Handlers
	heartbeatHandler := &routes.HeartbeatHandler{HeartbeatSrvc: heartbeatSrvc}
	aggregationHandler := &routes.AggregationHandler{AggregationSrvc: aggregationSrvc}

	// Middlewares
	authenticate := &middlewares.AuthenticateMiddleware{UserSrvc: userSrvc}

	// Setup Routing
	router := mux.NewRouter()
	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter()

	// API Routes
	heartbeats := apiRouter.Path("/heartbeat").Subrouter()
	heartbeats.Methods("POST").HandlerFunc(heartbeatHandler.Post)

	aggreagations := apiRouter.Path("/aggregation").Subrouter()
	aggreagations.Methods("GET").HandlerFunc(aggregationHandler.Get)

	// Sub-Routes Setup
	router.PathPrefix("/api").Handler(negroni.Classic().With(
		negroni.HandlerFunc(authenticate.Handle),
		negroni.Wrap(apiRouter),
	))

	// Listen HTTP
	portString := ":" + strconv.Itoa(config.Port)
	s := &http.Server{
		Handler:      router,
		Addr:         portString,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Listening on %+s\n", portString)
	s.ListenAndServe()
}
