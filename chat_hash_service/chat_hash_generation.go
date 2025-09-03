/*

The service for the registartion of the user 
register the user in the redis

*/


/*

import random

def generate_random_hash(low,high):
    """
    The function to generates a random hash of the given length.
    """

    random_length = random.randint(low, high)
    hash = ""
    for i in range(random_length):
        hash += random.choice("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    return has


*/


package main 

import (

	"fmt"
	"math/rand"
	"time"

)



func generateRandomHash(min int, max int ) string {
	/*
	The function to generate the random hash 
	*/

	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between min and max (inclusive)
	randomNumber := rand.Intn(max-min+1) + min

	var hash string

	choices := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 
	for i := 0; i <= randomNumber; i++ {
        
		choiceInt := rand.Intn(len(choices))

		hash += string(choices[choiceInt])

	
		}

	return hash

}




func main() {

	res := generateRandomHash(5,10)

	fmt.Println(res)
}


