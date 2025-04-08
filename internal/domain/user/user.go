package user

type Role string

const (
	RoleSupplier Role = "supplier"
	RoleReseller Role = "reseller"
	RoleConsumer Role = "consumer"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt string
}
