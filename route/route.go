package route

import (
	"github.com/gofiber/fiber/v2"
	"projectuas/app/service"
	"projectuas/middleware"
	"projectuas/database"
	"projectuas/app/repository"
)

func Setup(app *fiber.App) {

	// --- Database ---
	database.ConnectPostgres()
	mongoClient := database.ConnectMongo()
	repository.InitMongo(mongoClient)

	// --- PUBLIC AUTH ---
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", service.Login)
	auth.Post("/refresh", service.Refresh)
	auth.Post("/logout", service.Logout)

	// --- PROTECTED (JWT untuk semua user) ---
	api := app.Group("/api/v1", middleware.JWT())
	api.Get("/auth/profile", service.Profile)

	// --- USERS (ADMIN ONLY) ---
	users := api.Group("/users")
	users.Get("/", middleware.Role("admin"), service.AdminGetUsers)
	users.Get("/:id", middleware.Role("admin"), service.AdminGetUserDetail)
	users.Post("/", middleware.Role("admin"), service.AdminCreateUser)
	users.Put("/:id", middleware.Role("admin"), service.AdminUpdateUser)
	users.Delete("/:id", middleware.Role("admin"), service.AdminDeleteUser)
	users.Put("/:id/role", middleware.Role("admin"), service.AdminUpdateUserRole)

	// --- ACHIEVEMENTS ---
	achievements := api.Group("/achievements")

	// --- SPESIFIK DULU (urutan penting) ---
	achievements.Post("/:id/attachments", middleware.Role("admin"), service.UploadAchievementAttachments)
	achievements.Post("/:id/verify", middleware.Role("dosen_wali"), service.VerifyAchievement)
	achievements.Post("/:id/reject", middleware.Role("dosen_wali"), service.RejectAchievement)
	achievements.Post("/:id/submit", middleware.Role("mahasiswa"), service.SubmitAchievement)
	achievements.Put("/:id", middleware.Role("mahasiswa"), service.UpdateAchievement)
	achievements.Delete("/:id", middleware.Role("mahasiswa"), service.DeleteAchievement)
	achievements.Get("/:id/history", service.GetAchievementHistory)

	// --- GENERAL ROUTES (setelah spesifik) ---
	achievements.Get("/", service.GetAllAchievements)
	achievements.Get("/:id", service.GetAchievementDetail)
	achievements.Post("/", middleware.Role("mahasiswa"), service.CreateAchievement)

	// --- STUDENTS ---
	students := api.Group("/students")
	students.Get("/", service.GetStudents)
	students.Get("/:id", service.GetStudentDetail)
	students.Get("/:id/achievements", service.GetStudentAchievements)
	students.Put("/:id/advisor", middleware.Role("admin"), service.AssignAdvisor)

	// --- LECTURERS ---
	lecturers := api.Group("/lecturers")
	lecturers.Get("/", service.GetLecturers)
	lecturers.Get("/:id/advisees", middleware.Role("dosen_wali", "admin"), service.GetLecturerAdvisees)

	// --- REPORTS / STATISTICS (ADMIN ONLY) ---
	reports := api.Group("/reports")
	reports.Get("/statistics", middleware.Role("admin"), service.GetAchievementStatistics)
	reports.Get("/student/:id", middleware.Role("admin"), service.GetStudentReport)
}
