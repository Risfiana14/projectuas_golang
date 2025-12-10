package route

import (
	"github.com/gofiber/fiber/v2"
	"projectuas/app/service"
	"projectuas/middleware"
	"projectuas/database"
	"projectuas/app/repository"
)

func Setup(app *fiber.App) {

	// Database
	database.ConnectPostgres()
	mongoClient := database.ConnectMongo()
	repository.InitMongo(mongoClient)

	// PUBLIC AUTH
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", service.Login)
	auth.Post("/refresh", service.Refresh)
	auth.Post("/logout", service.Logout)

	// PROTECTED (semua user)
	api := app.Group("/api/v1", middleware.JWT())
	api.Get("/auth/profile", service.Profile)

	// USERS (ADMIN)
	users := api.Group("/users", middleware.Role("admin"))
	users.Get("/", service.AdminGetUsers)
	users.Get("/:id", service.AdminGetUserDetail)
	users.Post("/", service.AdminCreateUser)
	users.Put("/:id", service.AdminUpdateUser)
	users.Delete("/:id", service.AdminDeleteUser)
	users.Put("/:id/role", service.AdminUpdateUserRole)

	// ACHIEVEMENTS (SRS)
	achievements := api.Group("/achievements")
	achievements.Get("/", service.GetAllAchievements)
	achievements.Get("/:id", service.GetAchievementDetail)
	achievements.Get("/:id/history", service.GetAchievementHistory)
	achievements.Post("/:id/attachments", service.UploadAchievementAttachments)

	// mahasiswa only on same path
	mhsAch := achievements.Group("/", middleware.Role("mahasiswa"))
	mhsAch.Post("/", service.CreateAchievement)
	mhsAch.Put("/:id", service.UpdateAchievement)
	mhsAch.Delete("/:id", service.DeleteAchievement)
	mhsAch.Post("/:id/submit", service.SubmitAchievement)

	// dosen wali only
	dosenAch := achievements.Group("/", middleware.Role("dosen_wali"))
	dosenAch.Post("/:id/verify", service.VerifyAchievement)
	dosenAch.Post("/:id/reject", service.RejectAchievement)
	dosenAch.Get("/pending", service.GetPendingAchievements)

	// STUDENTS
	students := api.Group("/students")
	students.Get("/", service.GetStudents)
	students.Get("/:id", service.GetStudentDetail)
	students.Get("/:id/achievements", service.GetStudentAchievements)
	// assign advisor only admin
	students.Put("/:id/advisor", middleware.Role("admin"), service.AssignAdvisor)

	// LECTURERS
	lecturers := api.Group("/lecturers")
	lecturers.Get("/", service.GetLecturers)
	lecturers.Get("/:id/advisees", service.GetLecturerAdvisees)

	// REPORTS (admin)
	reports := api.Group("/reports", middleware.Role("admin"))
	reports.Get("/statistics", service.GetAchievementStatistics)
	reports.Get("/student/:id", service.GetStudentReport)
}
