package request

import (
	"context"
	"fmt"
)

// ExtractUserID retrieves the user_id from the context
func ExtractUserID(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value("user_id").(uint)
	if !ok {
		return 0, fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}
