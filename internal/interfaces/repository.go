package interfaces

type Repository interface {
	GetAccessByToken(token string) (int, error)
	GetAccessByIp(ip string) (int, error)
}
