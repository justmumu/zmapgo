package zmapgo

import (
	"fmt"
)

func multiPassChecker(args []string, checkArgument string) error {
	for _, arg := range args {
		if arg == checkArgument {
			return fmt.Errorf("found already added %s argument. Zmap does not allow multiple %s value", checkArgument, checkArgument)
		}
	}
	return nil
}
