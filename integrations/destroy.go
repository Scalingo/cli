package integrations

import "fmt"

func Destroy(scm_type string) error {
	// Call to scalingo api for removal

	fmt.Printf("Integration '%s' has been deleted.\n", scm_type)
	return nil
}
