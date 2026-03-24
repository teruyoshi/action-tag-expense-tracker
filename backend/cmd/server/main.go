package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"action-tag-expense-tracker/backend/handlers"
	"action-tag-expense-tracker/backend/repositories"
	"action-tag-expense-tracker/backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := getEnv("DATABASE_URL", "root:root@tcp(db:3306)/expense_tracker?charset=utf8mb4&parseTime=True&loc=Local")

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("waiting for db... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	tagRepo := &repositories.ActionTagRepository{DB: db}
	eventRepo := &repositories.EventRepository{DB: db}
	expenseRepo := &repositories.ExpenseRepository{DB: db}
	summaryRepo := &repositories.SummaryRepository{DB: db}
	balanceRepo := &repositories.BalanceRepository{DB: db}
	incomeCategoryRepo := &repositories.IncomeCategoryRepository{DB: db}
	incomeRepo := &repositories.IncomeRepository{DB: db}

	tagHandler := &handlers.ActionTagHandler{Repo: tagRepo}
	eventHandler := &handlers.EventHandler{Repo: eventRepo}
	expenseHandler := &handlers.ExpenseHandler{Repo: expenseRepo, BalanceRepo: balanceRepo}
	summaryHandler := &handlers.SummaryHandler{Repo: summaryRepo}
	balanceService := &services.BalanceService{
		BalanceRepo:   balanceRepo,
		ActionTagRepo: tagRepo,
		EventRepo:     eventRepo,
		ExpenseRepo:   expenseRepo,
	}
	balanceHandler := &handlers.BalanceHandler{
		Repo:    balanceRepo,
		Service: balanceService,
	}
	incomeCategoryHandler := &handlers.IncomeCategoryHandler{Repo: incomeCategoryRepo}
	incomeService := &services.IncomeService{
		IncomeRepo:  incomeRepo,
		BalanceRepo: balanceRepo,
	}
	incomeHandler := &handlers.IncomeHandler{
		Repo:    incomeRepo,
		Service: incomeService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/tags", tagHandler.List)
	r.Post("/tags", tagHandler.Create)
	r.Put("/tags/{id}", tagHandler.Update)
	r.Delete("/tags/{id}", tagHandler.Delete)

	r.Post("/events", eventHandler.Create)

	r.Post("/expenses", expenseHandler.Create)
	r.Put("/expenses/{id}", expenseHandler.Update)

	r.Get("/summary/month", summaryHandler.MonthTotal)
	r.Get("/summary/tag", summaryHandler.TagMonthTotals)
	r.Get("/summary/tag/diff", summaryHandler.TagMonthTotalsWithDiff)
	r.Get("/summary/tag/details", summaryHandler.TagExpenseDetails)

	r.Get("/balance", balanceHandler.Get)
	r.Put("/balance", balanceHandler.Update)

	r.Get("/income-categories", incomeCategoryHandler.List)
	r.Post("/income-categories", incomeCategoryHandler.Create)
	r.Put("/income-categories/{id}", incomeCategoryHandler.Update)
	r.Delete("/income-categories/{id}", incomeCategoryHandler.Delete)

	r.Get("/incomes", incomeHandler.List)
	r.Post("/incomes", incomeHandler.Create)
	r.Put("/incomes/{id}", incomeHandler.Update)
	r.Delete("/incomes/{id}", incomeHandler.Delete)

	port := getEnv("PORT", "8080")
	fmt.Printf("Server running on :%s\n", port)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
