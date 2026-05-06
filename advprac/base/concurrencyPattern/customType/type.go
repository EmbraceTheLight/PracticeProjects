package customType

import "time"

type Order struct {
	ID        int
	UserID    string
	Amount    float64
	Status    string
	CreatedAt time.Time
}
