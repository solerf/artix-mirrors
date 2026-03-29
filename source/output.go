package source

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

func write(source string, content io.Reader, writeTo io.Writer) error {
	buff := bytes.NewBuffer([]byte(
		fmt.Sprintf(`###############################################################################
## %s
## created at: %s
###############################################################################

`,
			source, time.Now().Format(time.DateTime))),
	)

	_, err := buff.ReadFrom(content)
	if err != nil {
		return fmt.Errorf("build content: %w", err)
	}

	_, err = writeTo.Write(buff.Bytes())
	if err != nil {
		return fmt.Errorf("write content: %w", err)
	}
	return nil
}
