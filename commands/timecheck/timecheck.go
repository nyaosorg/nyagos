package timecheck

// DoesYearHave366Days returns true if the `year` has 366 days.
func DoesYearHave366Days(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	if year%4 == 0 {
		return true
	}
	return false
}

// IsOk returns true if year,month,mday,hour,min,sec isnot invalid.
func IsOk(year, month, mday, hour, min, sec int) bool {
	if sec < 0 || sec > 60 {
		return false
	}
	if min < 0 || min >= 60 {
		return false
	}
	if hour < 0 || hour >= 24 {
		return false
	}
	if mday <= 0 {
		return false
	}
	var mdayMax int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		mdayMax = 31
	case 4, 6, 9, 11:
		mdayMax = 30
	case 2:
		if DoesYearHave366Days(year) {
			mdayMax = 29
		} else {
			mdayMax = 28
		}
	default:
		return false
	}
	if mday > mdayMax {
		return false
	}
	if year < 1900 {
		return false
	}
	return true
}
