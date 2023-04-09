package main

import (
	"log"
	myDB "manager/db"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func main() {
	db := myDB.InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		// Get parameters from the request body
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		// Validate parameters
		log.Println(telephone, len(telephone))
		if len(telephone) != 11 {
			ctx.JSON(200, gin.H{
				"code": 422,
				"msg":  "Invalid phone number",
			})
			return
		}
		if len(password) < 6 {
			ctx.JSON(200, gin.H{
				"code": 422,
				"msg":  "Password must be at least 6 characters",
			})
			return
		}
		// Check if the phone number already exists
		if isTelephoneExist(db, telephone) {
			ctx.JSON(200, gin.H{
				"code": 422,
				"msg":  "Phone number already exists",
			})
			return
		}

		// validate the name
		if len(name) == 0 {
			name = RandomString(10)
		}

		// register the user
		newUser := myDB.User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "Success",
		})
	})
	panic(r.Run(":8080"))
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	rand.Seed(time.Now().Unix())
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user myDB.User
	db.Where("telephone = ?", telephone).First(&user)
	return user.ID != 0
}
