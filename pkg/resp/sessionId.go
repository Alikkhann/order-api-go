package resp 

import "github.com/google/uuid"

func CreateSessionId() string {
    return uuid.New().String()
}