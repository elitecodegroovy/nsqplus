package tiny

import (
	"io"
	"os"
	"fmt"
	"strings"
)





var out io.Writer = os.Stdout // modified during testing

func Echo(newline bool, sep string, args []string) error {
	fmt.Fprint(out, strings.Join(args, sep))
	if newline {
		fmt.Fprintln(out)
	}
	return nil
}
