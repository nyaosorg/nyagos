package timecheck

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
	var mday_max int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		mday_max = 31
	case 4, 6, 9, 11:
		mday_max = 30
	case 2:
		if DoesYearHave366Days(year) {
			mday_max = 29
		} else {
			mday_max = 28
		}
	default:
		return false
	}
	if mday > mday_max {
		return false
	}
	if year < 1900 {
		return false
	}
	return true
}
