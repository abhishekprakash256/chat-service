/*
the main to test the function. 
*/


package main 

import (

	"chat-service/chat_hash_service"
	"fmt"

)


func main() {

	num := chat_hash_generation.GenerateRandomHash(5,10)

	fmt.Println(num)
	
}