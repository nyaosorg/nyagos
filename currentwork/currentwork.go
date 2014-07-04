package currentwork

import "fmt"
import "os"
import "strings"
import "unicode"

var lastfolder = map[rune]string{}

func getFirst(s string) (rune, error) {
	reader := strings.NewReader(s)
	drive, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return unicode.ToUpper(drive), nil
}

func Getwd() (string, error) {
	wd, err := os.Getwd()
	if err == nil {
		drive, driveErr := getFirst(wd)
		if driveErr == nil {
			lastfolder[drive] = wd
		}
	}
	return wd, err
}

func Chdrive(drive string) error {
	driveLetter, driveErr := getFirst(drive)
	if driveErr != nil {
		return driveErr
	}
	folder, ok := lastfolder[driveLetter]
	if ok {
		os.Chdir(folder)
	} else {
		os.Chdir(fmt.Sprintf("%c:.", driveLetter))
	}
	return nil
}

func Chdir(folder string) error {
	err := os.Chdir(folder)
	Getwd()
	return err
}
