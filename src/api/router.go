package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/api/handlers"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/middle"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Config struct {
	MasterDB data.MasterDB
}

// GetRouter создает и возвращает маршрутизатор без запуска сервера
func GetRouter(config Config) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), handlers.MasterDBContextKey, config.MasterDB)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/pet-and-health", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Post("/auth", handlers.Auth)
			r.Put("/forgot-password", handlers.ForgotPassword)
			r.Post("/logout", handlers.Logout)
			r.Post("/registration", handlers.Registration)
		})

		r.Route("/", func(r chi.Router) {
			r.Use(middle.JWTAuth)

			r.Route("/admin", func(r chi.Router) {
				r.Use(middle.CheckRole("admin"))

				r.Route("/pets", func(r chi.Router) {
					r.Get("/", handlers.GetPets)                  // Список всіх тварин
					r.Post("/", handlers.AddPet)                  // Додавання тварини
					r.Get("/{id}", handlers.PetInfo)              // Інформація про конкретну тварину
					r.Put("/{id}", handlers.UpdatePet)            // Оновлення інформації про тварину
					r.Delete("/{id}", handlers.DeletePet)         // Видалення тварини
					r.Post("/{id}/report", handlers.GetPetReport) // Генерація звіту про тварину
				})

				r.Route("/devices", func(r chi.Router) {
					r.Get("/", handlers.GetDevices)                   // Перегляд всіх пристроїв
					r.Post("/", handlers.AddDevice)                   // Додавання пристрою
					r.Put("/{id}", handlers.UpdateDevice)             // Оновлення пристрою
					r.Post("/{id}/assign", handlers.AssignDevice)     // Прив'язка пристрою до тварини
					r.Post("/{id}/unassign", handlers.UnassignDevice) // Відв'язка пристрою від тварини
				})

				r.Route("/users", func(r chi.Router) {
					r.Get("/", handlers.GetUsers)          // Перегляд всіх користувачів
					r.Post("/", handlers.CreateUser)       // Створення користувача
					r.Get("/{id}", handlers.UserInfo)      // Інформація про користувача
					r.Put("/{id}", handlers.UpdateUser)    // Оновлення користувача
					r.Delete("/{id}", handlers.DeleteUser) // Видалення користувача
				})
				r.Get("/profile", handlers.UserInfo)          // Профіль адміністратора
				r.Put("/profile", handlers.UpdateUserProfile) // Оновлення профілю адміна
			})

			// Маршрути для ветеринара
			r.Route("/vet", func(r chi.Router) {
				r.Use(middle.CheckRole("vet"))
				r.Get("/pets", handlers.GetPets)                   // Перегляд всіх тварин
				r.Get("/pets/{id}", handlers.PetInfo)              // Інформація про конкретну тварину
				r.Post("/pets/{id}/report", handlers.GetPetReport) // Генерація звіту про тварину
				r.Get("/users", handlers.GetUsers)
				r.Get("/profile", handlers.UserInfo)          // Профіль ветеринара
				r.Put("/profile", handlers.UpdateUserProfile) // Оновлення профілю ветеринара
			})

			r.Route("/owner", func(r chi.Router) {
				r.Use(middle.CheckRole("user"))
				r.Get("/pets", handlers.GetOwnerPets)         // Список тварин власника
				r.Get("/pets/{id}", handlers.OwnerPetInfo)    // Інформація про конкретну тварину власника
				r.Get("/profile", handlers.UserInfo)          // Профіль власника
				r.Put("/profile", handlers.UpdateUserProfile) // Оновлення профілю власника

				// Отримання сповіщень
				//r.Get("/notifications", handlers.GetNotifications)   // Отримання сповіщень про стан здоров'я
			})
		})

		// Точка прийому даних від пристроїв (можливо потрібна окрема авторизація)
		r.Route("/health-data", func(r chi.Router) {
			r.Post("/device", handlers.AddHealthData)
		})
	})

	return r
}

// Run создает маршрутизатор и запускает HTTP-сервер
func Run(config Config) {
	r := GetRouter(config)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
