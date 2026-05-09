package main

import (
	"fmt"
	"log"
	"ocs-room-booking/config"
	"ocs-room-booking/db"
	"ocs-room-booking/handlers"
	"ocs-room-booking/middleware"
	"ocs-room-booking/models"
	"ocs-room-booking/repository"
	"ocs-room-booking/services"
	"ocs-room-booking/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	config.Load()
	db.Connect()

	runMigrations()
	seedData()

	router := setupRouter()
	addr := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations() {
	gdb := db.GetDB()
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	_ = sqlDB

	// Run migrations via raw SQL since golang-migrate requires a file URL
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role VARCHAR(50) NOT NULL CHECK (role IN ('admin','core','viewer')),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS blocks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS rooms (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			block_id UUID REFERENCES blocks(id) ON DELETE CASCADE,
			room_name VARCHAR(100) NOT NULL,
			capacity INTEGER NOT NULL,
			is_available BOOLEAN DEFAULT true,
			allowed_purposes TEXT[] DEFAULT '{OA,Interview,PPT}',
			notes TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			date DATE NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			purpose VARCHAR(50) CHECK (purpose IN ('OA','Interview','PPT')),
			participant_count INTEGER NOT NULL,
			status VARCHAR(50) DEFAULT 'confirmed' CHECK (status IN ('confirmed','cancelled')),
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS permissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			scope VARCHAR(100) NOT NULL
		)`,
	}

	for _, m := range migrations {
		if err := gdb.Exec(m).Error; err != nil {
			log.Printf("Migration warning: %v", err)
		}
	}
	log.Println("Migrations applied")
}

func seedData() {
	gdb := db.GetDB()

	// --- Seed Admin ---
	cfg := config.AppConfig
	var adminCount int64
	gdb.Model(&models.User{}).Where("role = 'admin'").Count(&adminCount)
	if adminCount == 0 {
		hash, err := utils.HashPassword(cfg.AdminPassword)
		if err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}
		admin := models.User{
			Name:         "OCS Admin",
			Email:        cfg.AdminEmail,
			PasswordHash: hash,
			Role:         "admin",
			IsActive:     true,
		}
		if err := gdb.Create(&admin).Error; err != nil {
			log.Printf("Admin seed warning: %v", err)
		} else {
			log.Println("Admin user seeded")
		}
	}

	// --- Seed Blocks ---
	blockNames := []string{
		"A Block", "B Block", "C Block", "CSE Block",
		"LHC", "BT/BM", "CY", "EE", "MA", "MSME", "PH",
	}
	blockMap := make(map[string]uuid.UUID)
	for _, name := range blockNames {
		var block models.Block
		if err := gdb.Where("name = ?", name).First(&block).Error; err != nil {
			block = models.Block{Name: name}
			gdb.Create(&block)
		}
		blockMap[name] = block.ID
	}
	log.Println("Blocks seeded")

	// --- Seed Rooms ---
	type roomSeed struct {
		block    string
		name     string
		capacity int
	}
	rooms := []roomSeed{
		{"A Block", "A-Class Room 320", 80},
		{"A Block", "A-AUDITORIUM", 289},
		{"A Block", "A-Class Room 111", 70},
		{"A Block", "A-Class Room 112", 80},
		{"A Block", "A-Class Room 114", 36},
		{"A Block", "A-Class Room 117", 84},
		{"A Block", "A-Class Room 118", 84},
		{"A Block", "A-Class Room 119", 108},
		{"A Block", "A-Class Room 220", 40},
		{"A Block", "A-Class Room 221", 120},
		{"A Block", "A-LH-1", 184},
		{"A Block", "A-LH-2", 184},
		{"BT/BM", "BT/BM-009", 24},
		{"BT/BM", "BT/BM-010", 24},
		{"BT/BM", "BT/BM-118", 60},
		{"C Block", "C-LH-10", 68},
		{"C Block", "C-LH-2", 138},
		{"C Block", "C-LH-3", 100},
		{"C Block", "C-LH-4", 60},
		{"C Block", "C-LH-5", 60},
		{"C Block", "C-LH-6", 60},
		{"C Block", "C-LH-7", 70},
		{"C Block", "C-LH-9", 66},
		{"CSE Block", "CSE-LH-01", 70},
		{"CSE Block", "CSE-LH-02", 70},
		{"CSE Block", "CSE-LH-03", 70},
		{"CY", "CY-LH-1", 30},
		{"CY", "CY-LH-2", 40},
		{"CY", "CY-LH-3", 90},
		{"EE", "EE-004(GF)", 80},
		{"EE", "EE-20 (SF)", 60},
		{"LHC", "LHC-01", 72},
		{"LHC", "LHC-02", 72},
		{"LHC", "LHC-03", 120},
		{"LHC", "LHC-04", 200},
		{"LHC", "LHC-05", 800},
		{"LHC", "LHC-06", 320},
		{"LHC", "LHC-07", 200},
		{"LHC", "LHC-08", 120},
		{"LHC", "LHC-09", 72},
		{"LHC", "LHC-10", 72},
		{"LHC", "LHC-11", 120},
		{"LHC", "LHC-12", 200},
		{"LHC", "LHC-13", 320},
		{"LHC", "LHC-14", 200},
		{"LHC", "LHC-15", 120},
		{"MA", "MA-01", 56},
		{"MA", "MA-02", 56},
		{"MA", "MA-114", 30},
		{"MSME", "MSME-LH-1", 36},
		{"MSME", "MSME-LH-2", 60},
		{"MSME", "MSME-LH-3", 106},
		{"PH", "PH-1", 80},
		{"PH", "PH-2", 60},
		{"PH", "PH-3", 50},
	}

	var roomCount int64
	gdb.Model(&models.Room{}).Count(&roomCount)
	if roomCount == 0 {
		for _, r := range rooms {
			blockID, ok := blockMap[r.block]
			if !ok {
				log.Printf("Block not found for room %s: %s", r.name, r.block)
				continue
			}
			room := models.Room{
				BlockID:         blockID,
				RoomName:        r.name,
				Capacity:        r.capacity,
				IsAvailable:     true,
				AllowedPurposes: models.StringArray{"OA", "Interview", "PPT"},
			}
			gdb.Create(&room)
		}
		log.Println("Rooms seeded")
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	gdb := db.GetDB()

	// Repositories
	userRepo := repository.NewUserRepository(gdb)
	roomRepo := repository.NewRoomRepository(gdb)
	bookingRepo := repository.NewBookingRepository(gdb)

	// Services
	userSvc := services.NewUserService(userRepo)
	roomSvc := services.NewRoomService(roomRepo)
	bookingSvc := services.NewBookingService(bookingRepo, roomRepo)

	// Handlers
	authH := handlers.NewAuthHandler(userRepo)
	userH := handlers.NewUserHandler(userSvc)
	roomH := handlers.NewRoomHandler(roomSvc)
	bookingH := handlers.NewBookingHandler(bookingSvc)

	api := router.Group("/api")

	// Public routes
	api.POST("/auth/login", authH.Login)

	// Protected routes (any authenticated user)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/rooms/search", roomH.SearchRooms)
		protected.GET("/rooms/:id", roomH.GetRoom)
		protected.GET("/blocks", roomH.GetAllBlocks)

		// core + admin can book
		protected.POST("/bookings", middleware.RequireRole("core", "admin"), bookingH.CreateBooking)
		protected.GET("/bookings/my", bookingH.GetMyBookings)
		protected.GET("/bookings/all", middleware.RequireRole("viewer", "admin"), bookingH.GetAllBookings)
		protected.DELETE("/bookings/:id", bookingH.CancelMyBooking)
	}

	// Admin only routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		admin.POST("/users", userH.CreateUser)
		admin.PATCH("/users/:id", userH.UpdateUser)
		admin.GET("/users", userH.ListUsers)

		admin.GET("/rooms", roomH.ListAllRooms)
		admin.POST("/rooms", roomH.CreateRoom)
		admin.PATCH("/rooms/:id", roomH.UpdateRoom)
		admin.DELETE("/rooms/:id", roomH.DeleteRoom)

		admin.GET("/bookings", bookingH.GetAllBookings)
		admin.DELETE("/bookings/:id", bookingH.AdminCancelBooking)
	}

	return router
}
