package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomAlphanumericCode(n int) (string, error) {
  b := make([]byte, n)
  _, err := rand.Read(b)
  if err != nil {
    return "", fmt.Errorf("failed to read random bytes: %w", err)
  }

  for i := range b {
    // Loop until a valid alphanumeric character is generated
    for {
      randomByte := b[i]
      asciiValue := int(randomByte)
      if (asciiValue >= 'a' && asciiValue <= 'z') ||
         (asciiValue >= 'A' && asciiValue <= 'Z') ||
         (asciiValue >= '0' && asciiValue <= '9') {
        break // Valid alphanumeric character found, exit loop
      }
      // If not alphanumeric, regenerate the byte
      _, err := rand.Read(b[i:i+1])
      if err != nil {
        return "", fmt.Errorf("failed to read random byte: %w", err)
      }
    }
  }

  return base64.URLEncoding.EncodeToString(b)[:n], nil
}

// func main() {
//   // Example usage
//   code, err := GenerateRandomAlphanumericCode(10)
//   if err != nil {
//     panic(err) // Replace with proper error handling in production
//   }
//   println(code)
// }
