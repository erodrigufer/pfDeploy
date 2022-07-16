package shCmd

import (
	"fmt"
	"io"
	"os"
)

// reboot, reboots the system.
func reboot() error {
	_, err := app.shCmd("reboot")
	if err != nil {
		err = fmt.Errorf("reboot attempt failed: %w", err)
		return err
	}
	// The program should never come this far, since it would reboot before.
	return nil
}

// TODO: change file attributes after copy or create a special function for that
// copyFile, copy the src file to dst. Any existing file will be overwritten
// and will not copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error while trying to open file ('%s'): %w", src, err)
	}
	defer in.Close()

	// The 'dst' file will be created, or truncated if it already exists
	// (overwritten). 'dst' file has file mode 0666.
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error while trying to create file ('%s'): %w", dst, err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("error while trying to copy data from file '%s' to file '%s': %w", src, dst, err)
	}

	return nil
}
