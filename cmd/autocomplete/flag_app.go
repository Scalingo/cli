package autocomplete

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-scalingo/v9/debug"
)

func FlagAppAutoComplete(ctx context.Context) bool {
	apps, err := appsList(ctx)
	if err != nil {
		debug.Println("fail to get apps list:", err)
		return false
	}

	for _, app := range apps {
		fmt.Println(app.Name)
	}

	return true
}
