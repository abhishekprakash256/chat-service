package register

import (
	"context"
	"encoding/json"
	//"log"
	"net/http"

	"chat-service/internal/config"
	"chat-service/internal/hash"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"
)

type SuccessResponse struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}



// writeError writes a standardized JSON error response.
//
// Params:
//   - w: http.ResponseWriter to write the response.
//   - code: HTTP status code (e.g., 400, 500).
//   - msg: human-readable error message.
//
// Example Response:
//   {
//     "status": "error",
//     "code": 400,
//     "message": "Invalid JSON"
//   }
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: msg,
	})
}




// UserRegistration handles chat registration requests.
//
// It expects a POST request with JSON body in the form:
//   {
//     "userOne": "Sam",
//     "userTwo": "Bob"
//   }
//
// Steps performed:
//   1. Validates that the request method is POST.
//   2. Decodes the JSON body into RegistrationtData struct.
//   3. Generates a unique hash using Redis.
//   4. Saves the userOne, userTwo, and hash into the PostgreSQL "login" table.
//   5. Returns a JSON response with the inserted data and generated hash.
//
// Success Response Example:
//   {
//     "status": "success",
//     "code": 200,
//     "data": {
//       "userOne": "Sam",
//       "userTwo": "Bob",
//       "hash": "abc123"
//     }
//   }
//
// Error Response Example:
//   {
//     "status": "error",
//     "code": 400,
//     "message": "Invalid JSON"
//   }
//
// Dependencies:
//   - config.GlobalDbConn (PostgreSQL + Redis connections)
//   - hash (for unique hash generation)
//   - pgsqlcrud (for DB insertion)
//
// Route:
//   POST /chat-server/registration
func UserRegistration(w http.ResponseWriter, r *http.Request) {

	pool := config.GlobalDbConn.PgsqlConn
	client := config.GlobalDbConn.RedisConn
	ctx := context.Background()

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	// Decode request
	var data config.RegistrationtData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Generate hash
	hash.GenerateUniqueHash(config.UniqueHashSet, config.UsedHashSet, 5, 10, 20, client)
	uniqueHash := hash.PopUniqueHash(config.UniqueHashSet, config.UsedHashSet, client)

	// Save login data
	logindata := config.LoginData{
		ChatID:  uniqueHash,
		UserOne: data.UserOne,
		UserTwo: data.UserTwo,
	}
	if !pgsqlcrud.InsertLoginData(ctx, config.LoginTable , pool, logindata) {
		writeError(w, http.StatusInternalServerError, "Database insert failed")
		return
	}

	// Success Response
	resp := SuccessResponse{
		Status: "success",
		Code:   http.StatusOK,
		Data: map[string]string{
			"userOne": data.UserOne,
			"userTwo": data.UserTwo,
			"hash":    uniqueHash,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
