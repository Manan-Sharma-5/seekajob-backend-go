package repository

// import (
// 	"log"
// 	"time"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// 	"github.com/gin-gonic/gin"
// )

// func PreviousYearUpload(c *gin.Context) {
//     region := "eu-north-1"
//     bucket := "software-engineering-project-s3"

//     // Retrieve query parameters for file path and name
//     filename := c.Query("filename")
//     if filename == "" {
//         c.JSON(400, gin.H{"error": "Filename is required"})
//         return
//     }

//     // Get the user ID from the cookie (or however it is stored in your application)
//     userID, err := c.Cookie("user_id")
//     if err != nil {
//         c.JSON(401, gin.H{"error": "User not authenticated"})
//         return
//     }

//     // Create an AWS session
//     sess, err := session.NewSession(&aws.Config{
//         Region: aws.String(region),
//     })
//     if err != nil {
//         log.Println("Failed to create AWS session:", err)
//         c.JSON(500, gin.H{"error": "Failed to create AWS session"})
//         return
//     }

//     // Create an S3 service client
//     svc := s3.New(sess)

//     // Generate a pre-signed PUT URL for the client to upload the file
//     key := "previous-year-question-paper/" + year + "/" + subjectCode + "/" + filename
//     putReq, _ := svc.PutObjectRequest(&s3.PutObjectInput{
//         Bucket: aws.String(bucket),
//         Key:    aws.String(key),
//     })
//     putURL, err := putReq.Presign(15 * time.Minute)
//     if err != nil {
//         log.Println("Failed to sign PUT request:", err)
//         c.JSON(500, gin.H{"error": "Failed to sign PUT request"})
//         return
//     }

//     // Generate a pre-signed GET URL for accessing the uploaded file
//     getReq, _ := svc.GetObjectRequest(&s3.GetObjectInput{
//         Bucket: aws.String(bucket),
//         Key:    aws.String(key),
//     })
// 		getURL, err := getReq.Presign(7 * 24 * time.Hour)
//     if err != nil {
//         log.Println("Failed to sign GET request:", err)
//         c.JSON(500, gin.H{"error": "Failed to sign GET request"})
//         return
//     }

// }
