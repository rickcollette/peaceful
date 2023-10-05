package main

import (
    "fmt"
    "github.com/rickcollette/peaceful/tools" 
)

func main() {
    // Generate a UUID using the GenerateUUID function
    uuidString, err := tools.GenerateUUID()
    if err != nil {
        fmt.Println("Error generating UUID:", err)
        return
    }

    fmt.Println("Generated UUID:", uuidString)

    // You can now use 'uuidString' in your application as needed.
}
