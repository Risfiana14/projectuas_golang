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

    // PUBLIC
    auth := app.Group("/api/v1/auth")

    auth.Post("/login", service.Login)
    auth.Post("/logout", service.Logout)
    auth.Post("/refresh", service.Refresh)

    // PROTECTED
    api := app.Group("/api/v1", middleware.JWT())

    api.Get("/auth/profile", service.Profile)

    // ADMIN ONLY
    users := api.Group("/users", middleware.Role("admin"))
    users.Get("/", service.AdminGetUsers)
    users.Get("/:id", service.AdminGetUserDetail)
    users.Post("/", service.AdminCreateUser)
    users.Put("/:id", service.AdminUpdateUser)
    users.Delete("/:id", service.AdminDeleteUser)
    users.Put("/:id/role", service.AdminUpdateUserRole)

    // MAHASISWA
    mhs := api.Group("/achievements", middleware.Role("mahasiswa"))
    mhs.Post("/", service.CreateAchievement)
    mhs.Get("/", service.GetMyAchievements)
    mhs.Get("/:id", service.GetAchievementDetail)
    mhs.Put("/:id", service.UpdateAchievement)
    mhs.Delete("/:id", service.DeleteAchievement)
    mhs.Post("/:id/submit", service.SubmitAchievement)

    // DOSEN WALI
    dosen := api.Group("/verify", middleware.Role("dosen_wali"))
    dosen.Get("/pending", service.GetPendingAchievements)
    dosen.Post("/:id/verify", service.VerifyAchievement)
    dosen.Post("/:id/reject", service.RejectAchievement)
}
