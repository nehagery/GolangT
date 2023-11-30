package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// user represents data about a new user.
type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	PhoneNo    string    `json:"phone_number"`
	Otp        string    `json:"otp"`
	OtpExpTime time.Time `json:"otp_expiration_time"`
}

// albums slice to seed record album data.
/*
var users = []user{
	{Id: 1, Name: "Train", PhoneNo: "12345678", Otp: "123", OtpExpTime: "2012-04-21:25:43-05:00"},
	{Id: 2, Name: "Jeru", PhoneNo: "102345678", Otp: "456", OtpExpTime: "2012-04-21T18:25:43-05:00"},
	{Id: 3, Name: "Sarah ", PhoneNo: "12045678", Otp: "789", OtpExpTime: "2012-04-21T18:25:43-05:00"},
}
*/

var users = []User{
	{Id: 1, Name: "Train", PhoneNo: "12345678"},
	{Id: 2, Name: "Jeru", PhoneNo: "102345678"},
	{Id: 3, Name: "Sarah ", PhoneNo: "12045678"},
}

// getUsers responds with the list of all albums as JSON.
func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

// postUsers adds a user from JSON received in the request body.
func postUsers(c *gin.Context) {
	var newUser User

	// Call BindJSON to bind the received JSON to
	// newUser.
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	// if phone no exist throw 400 error
	if getUserByPhoneNo(newUser.PhoneNo) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "400 Error ,user not unique"})
	} else {
		users = append(users, newUser)
		c.IndentedJSON(http.StatusCreated, newUser)
	}

}

// parameter sent by the client, then returns that user as a response.
func getUserByPhoneNo(phoneNo string) bool {
	fmt.Println("Phone no is", phoneNo)
	// Loop over the list of users, looking for
	// a user whose phoneNumber value matches the parameter.
	for _, a := range users {
		if a.PhoneNo == phoneNo {
			return true

		}
	}
	return false
}

// generateOtp generates Otp for the valid phone no, adds a user from JSON received in the request body.
func generateOtp(c *gin.Context) {
	var existingUser User

	// Call BindJSON to bind the received JSON to
	// existingUser.
	if err := c.BindJSON(&existingUser); err != nil {
		return
	}
	// if phone no exists, generate OTP and save user
	if getUserByPhoneNo(existingUser.PhoneNo) {
		fmt.Println("Phone no is available in DB", existingUser.PhoneNo)
		existingUser.Otp = generateRandom4DigitOTP()
		existingUser.OtpExpTime = time.Now().Add(1 * time.Minute)
		//TODO update the user in slice with OTP and estimated time
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": " 404 Error , user not found"})
	}

}
func generateRandom4DigitOTP() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

// generateOtp generates Otp for the valid phone no, adds a user from JSON received in the request body.
func verifyOtp(c *gin.Context) {
	var existingUser User

	// Call BindJSON to bind the received JSON to
	// existingUser.
	if err := c.BindJSON(&existingUser); err != nil {
		return
	}
	// if phone no exists, generate OTP and save user
	if validateOTP(existingUser.PhoneNo, existingUser.Otp) {
		if validateExpirationTimeOTP(existingUser.PhoneNo) {
			c.IndentedJSON(http.StatusOK, gin.H{"message": " OTP Validation Success"})
		} else {
			c.IndentedJSON(http.StatusRequestTimeout, gin.H{"message": "Error ,OTP has expired"})
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": " Error , OTP is incorrect"})
	}

}

func validateExpirationTimeOTP(phoneNo string) bool {
	fmt.Println("OTP exp time given is", otp)
	// Loop over the list of users, looking for
	// a user whose OTP value matches the parameter.
	for _, a := range users {
		if a.PhoneNo == phoneNo {
			if time.Now().Before(a.OtpExpTime) {
				return true
			}
		}
	}
	return false
}

func validateOTP(phoneNo string, otp string) bool {
	fmt.Println("OTP given is", otp)
	// Loop over the list of users, looking for
	// a user whose OTP value matches the parameter.
	for _, a := range users {
		if a.PhoneNo == phoneNo {
			if a.Otp == otp {
				return true
			}
		}
	}
	return false
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/api/users", postUsers)
	router.POST("/api/users/generateotp", generateOtp)
	router.POST("/api/users/verifyotp", verifyOtp)

	router.Run("localhost:8080")
}
