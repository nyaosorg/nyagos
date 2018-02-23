package dos

func Copy(src, dst string, isFailIfExists bool) error {
	rc, err := copyFile(src, dst, isFailIfExists)
	if rc == 0 {
		return err
	}
	return nil
}

func Move(src, dst string) error {
	rc, err := moveFileEx(src, dst,
		MOVEFILE_REPLACE_EXISTING|MOVEFILE_COPY_ALLOWED|MOVEFILE_WRITE_THROUGH)
	if rc == 0 {
		return err
	}
	return nil
}
