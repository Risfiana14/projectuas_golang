// route/route.go
package route

import (
    "projectuas/app/repository"
    "projectuas/app/service"
    "projectuas/middleware"

    "github.com/gofiber/fiber/v2"
    "database/sql"
    "go.mongodb.org/mongo-driver/mongo"
)

func Setup(appFiber *fiber.App, pg *sql.DB, mongoClient *mongo.Client) {
    repository.Init(pg, mongoClient)

    // Public
    appFiber.Post("/api/v1/auth/login", service.Login)

    // Protected
    api := appFiber.Group("/api/v1", middleware.JWT())
    api.Get("/profile", service.Profile)

    // === ADMIN ONLY: Manage Users (FR-009) ===
    admin := api.Group("", middleware.Role("admin"))
    admin.Get("/users", service.GetUsers)
    admin.Post("/users", service.CreateUserAdmin)
    admin.Put("/users/:id/role", service.UpdateUserRole)
    admin.Delete("/users/:id", service.DeleteUserAdmin)
    admin.Get("/achievements/all", service.GetAllAchievements)
    
    // Mahasiswa
    mhs := api.Group("", middleware.Role("mahasiswa"))
    mhs.Post("/achievements", service.CreateAchievement)
    mhs.Get("/achievements", service.GetMyAchievements)
    mhs.Get("/achievements/:id", service.GetAchievementDetail)
    mhs.Put("/achievements/:id", service.UpdateAchievement)
    mhs.Delete("/achievements/:id", service.DeleteAchievement)
    mhs.Post("/achievements/:id/submit", service.SubmitAchievement)

    // Dosen Wali
    dosen := api.Group("", middleware.Role("dosen_wali"))
    dosen.Get("/achievements/pending", service.GetPendingAchievements)
    dosen.Post("/achievements/:id/verify", service.VerifyAchievement)
    dosen.Post("/achievements/:id/reject", service.RejectAchievement)

}