package humble

import (
	"fmt"
)

type ViewElementNotFoundError struct {
	viewId string
}

func (e ViewElementNotFoundError) Error() string {
	return fmt.Sprintf("Could not find element in index or DOM for view with id: %s.", e.viewId)
}

func NewViewElementNotFoundError(viewId string) ViewElementNotFoundError {
	return ViewElementNotFoundError{viewId: viewId}
}
