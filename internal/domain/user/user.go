package user

import "time"

type Role string

const (
	RoleSupplier Role = "supplier"
	RoleReseller Role = "reseller"
	RoleConsumer Role = "consumer"
	RoleAdmin    Role = "admin"
)


type User struct {
    ID        string    `bson:"_id"`
    Name      string    `bson:"name"`
    Username  string    `bson:"username"`
    Email     string    `bson:"email"`
    Password  string    `bson:"password"`
    Role      string    `bson:"role"`
    CreatedAt time.Time `bson:"created_at"`
}