package helper

// import (
// 	"fmt"
// 	"math/rand"
// 	"strings"
// 	"time"
// )

// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// func init() {
// 	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
// }

// func generateCode(length int) string {
// 	var sb strings.Builder
// 	for i := 0; i < length; i++ {
// 		sb.WriteByte(charset[rand.Intn(len(charset))])
// 	}
// 	return sb.String()
// }

// func GenerateRandomAlphanumericCode(n int) string {
//   // Seed the random number generator with current time
//   rand.Seed(time.Now().UnixNano())

//   b := make([]byte, n)
//   for i := range b {
//     // Generate random index within the alphanumeric string length
//     index := rand.Intn(len(charset))
//     b[i] = charset[index]
//   }
//   return string(b)
// }

// func main() {
//   // Example usage
//   fmt.Println(generateCode(8))

//   code := GenerateRandomAlphanumericCode(10)
//   println(code)
// }
