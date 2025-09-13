/*
make the logout 
*/

package logout 


import (

	"fmt"
	"net/http"
)


func LogOutUser(w http.ResponseWriter, r *http.Request ) {

	fmt.Println("....... logout clicked ...... ")

}