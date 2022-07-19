package rollback

// RollbackableAction is an interface to represents an install action that may be rolled back.
type RollbackableAction interface {
	Rollback() error
}
