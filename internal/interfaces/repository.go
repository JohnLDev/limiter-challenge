package interfaces

type Repository interface {
	Count(key string) (int, error)
	Save(key string, id string) error
	CheckLock(key string) (bool, error)
	LockKey(key string) error
}
