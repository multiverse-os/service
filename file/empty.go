package file


type unit int64

const (
	GiB unit = 1024
	GB  unit = 1000
)

func CreateEmptyFile(path string, size int64, u unit) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(path)
		}
	}()

	switch u {
	case GiB:
		size = size * 1024 * 1024 * 1024
	case GB:
		// The image size will be reduced to fit commercial drives that are
		// smaller than what they claim, 975 comes from 97.5% of the total size
		// but we want to be a multiple of 512 (and size is an int) we divide by
		// 512 and multiply it again
		size = size * 1000 * 1000 * 975 / 512 * 512
	default:
		panic("improper sizing unit used")
	}

	if err := file.Truncate(size); err != nil {
		return errors.New(fmt.Sprintf("Error creating %s of size %d to stage image onto", path, size))
	}
	return nil
}
