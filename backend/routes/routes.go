package routes

import (
	"net/http"
	"wechat/auth"
	"wechat/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, database *gorm.DB) {
	protected := router.Group("/api")
	protected.Use(auth.JWTMiddleware())

	//Register
	router.POST("/register", func(ctx *gin.Context) {
		var registerData struct {
			UserName string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&registerData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashedPassword, err := HashPassword(registerData.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure Password"})
			return
		}
		newUser := models.User{
			UserName: registerData.UserName,
			Email:    registerData.Email,
			Password: hashedPassword,
		}
		if err := database.Create(&newUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Create User"})
			return
		}
		safeUser := struct {
			UserName string
			Email    string
		}{
			UserName: newUser.UserName,
			Email:    newUser.Email,
		}
		ctx.JSON(http.StatusOK, safeUser)
	})
	//Login
	router.POST("/login", func(ctx *gin.Context) {
		var loginData struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&loginData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := database.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ok := CheckPasswordHash(loginData.Password, user.Password)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
			return
		}
		token, err := auth.GenerateJWT(user.ID.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token"})
			return
		}
		type SafeUser struct {
			ID       uuid.UUID         `json:"id"`
			UserName string            `json:"username"`
			Email    string            `json:"email"`
			Rooms    []*models.Room    `json:"rooms"`
			Messages []*models.Message `json:"messages"`
		}

		safeUser := SafeUser{
			ID:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
			Rooms:    user.Rooms,
			Messages: user.Messages,
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user":  safeUser,
			"token": token,
		})
	})
	//Request a room
	protected.POST("/room", func(ctx *gin.Context) {
		type RoomRequest struct {
			User1ID string `json:"user1_id" binding:"required"`
			User2ID string `json:"user2_id" binding:"required"`
		}

		var body RoomRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		user1UUID, err := uuid.Parse(body.User1ID)
		user2UUID, err2 := uuid.Parse(body.User2ID)
		if err != nil || err2 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID(s)"})
			return
		}
		if user1UUID == user2UUID {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create a room with yourself"})
			return
		}
		userIDValue, ok := ctx.Get("userID") // Your auth middleware
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not in context"})
			return
		}
		currentUserID, ok := userIDValue.(uuid.UUID)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
			return
		}

		// Validate users exist
		var users []models.User
		if err := database.Where("id IN ?", []uuid.UUID{user1UUID, user2UUID}).Find(&users).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "DB error fetching users"})
			return
		}
		if len(users) != 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "One or both users do not exist"})
			return
		}

		// Check existing room
		var existingRoom models.Room
		err = database.Joins("JOIN user_rooms ur1 ON ur1.room_id = rooms.id").
			Joins("JOIN user_rooms ur2 ON ur2.room_id = rooms.id").
			Where("ur1.user_id = ? AND ur2.user_id = ?", user1UUID, user2UUID).
			First(&existingRoom).Error

		if err == nil {
			// Room exists, return with dynamic name
			roomName, err := getRoomNameForUser(database, existingRoom.ID, currentUserID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room name"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"message":   "Room already exists",
				"room_id":   existingRoom.ID,
				"room_name": roomName,
				"room":      existingRoom,
			})
			return
		}

		// Create room + assign users
		var room models.Room
		err = database.Transaction(func(tx *gorm.DB) error {
			room = models.Room{}
			if err := tx.Create(&room).Error; err != nil {
				return err
			}

			userRooms := []models.UserRoom{
				{RoomID: room.ID, UserID: user1UUID},
				{RoomID: room.ID, UserID: user2UUID},
			}
			if err := tx.Create(&userRooms).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room: " + err.Error()})
			return
		}

		// Return the new room name for the current user
		roomName, err := getRoomNameForUser(database, room.ID, currentUserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room name"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":   "Room created",
			"room_id":   room.ID,
			"room_name": roomName,
			"room":      room,
		})
	})
	//Get all user
	protected.GET("/users", func(ctx *gin.Context) {
		var users []models.User
		if err := database.Select("id", "user_name", "email").Find(&users).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"users": users})
	})
	protected.GET("/room/:roomID", func(ctx *gin.Context) {
		roomIDstr := ctx.Param("roomID")
		roomID, err := uuid.Parse(roomIDstr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
			return
		}
		var room models.Room
		if err := database.Preload("Members").Where("id = ? ", roomID).First(&room).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Room does not exist"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "room fetched",
			"room":    room,
			"members": room.Members,
		})
	})

}
func getRoomNameForUser(db *gorm.DB, roomID uuid.UUID, currentUserID uuid.UUID) (string, error) {
	var userRooms []models.UserRoom
	if err := db.Where("room_id = ?", roomID).Find(&userRooms).Error; err != nil {
		return "", err
	}

	var otherUserID uuid.UUID
	for _, ur := range userRooms {
		if ur.UserID != currentUserID {
			otherUserID = ur.UserID
			break
		}
	}

	var otherUser models.User
	if err := db.First(&otherUser, "id = ?", otherUserID).Error; err != nil {
		return "", err
	}
	return otherUser.UserName, nil
}
