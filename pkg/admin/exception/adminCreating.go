package exception

import "fmt"

type AdminCreating struct {
	AdminID string
}

func (e *AdminCreating) Error() string {
	return fmt.Sprintf("creating AdminID: %s failed", e.AdminID)
}
