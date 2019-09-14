package file

import (
	"bufio"
	"io"
	"os"
)

func CopyFile(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, srcStat.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	reader := bufio.NewReader(srcFile)
	writer := bufio.NewWriter(dstFile)
	defer func() {
		if err != nil {
			writer.Flush()
		}
	}()
	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	writer.Flush()
	return nil
}
