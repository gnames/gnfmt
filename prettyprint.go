package gnfmt

import (
	"encoding/json"
	"fmt"
)

// Ppr is a pretty print of an object
func Ppr(obj any) string {
	res, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return string(res)
}
