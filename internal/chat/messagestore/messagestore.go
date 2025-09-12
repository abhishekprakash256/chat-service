/*
the package to save the message from chat 
from user using pgsql 

*/

package messagestore

import (

	"time"
	"context"

	"chat-service/internal/config"
	
)



funn SaveMessage(Sender string) {

	// get the pgsql conn 
	pool := config.GlobalDbConn.PgsqlConn

	// start the context 
	ctx := context.Background()
}