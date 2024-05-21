package date

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate - функция для вычисления следующей даты
func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("нет значения повтора")
	}

	// Разделяем строку повтора на части
	repeatParts := strings.Split(repeat, " ")

	// В зависимости от первой буквы (d - дни, y - год, w - неделя, m - месяц)
	// используем соответствующий вычислитель
	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", errors.New("неправильный формат повтора")
		}

		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days > 400 {
			return "", errors.New("недопустимое количество дней или неверный формат числа")
		}

		for {
			startDate = startDate.AddDate(0, 0, days)

			if !startDate.Before(now) && !startDate.Equal(now) {
				break
			}
		}

	case "y":

		if len(repeatParts) != 1 {
			return "", errors.New("неправильный формат повтора")
		}
		for {
			startDate = startDate.AddDate(1, 0, 0)

			if !startDate.Before(now) && !startDate.Equal(now) {
				break
			}
		}

	case "w":

		var weekDays []int

		if len(repeatParts) != 2 {
			return "", errors.New("неправильный формат повтора")
		}

		// strings.Join(repeatParts, " ")

		repeatNumW := strings.Split(repeatParts[1], ",")
		if len(repeatNumW) == 0 {
			return "", errors.New("не указан интервал в днях недели")
		}

		for _, dayStr := range repeatNumW { // 4,5
			repDay, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", errors.New("неверный формат дня недели")
			}

			if repDay < 1 || repDay > 7 {
				return "", errors.New("недопустимое значение дня недели")
			}
			weekDays = append(weekDays, repDay)
		}

		for i, day := range weekDays {
			if day == 7 {
				weekDays[i] = 0
			}
		}

		sort.Ints(weekDays)

		// log.Printf("weekDays: %v", weekDays)

		var nextWeekDay int
		for _, wd := range weekDays { // 4,5
			// log.Printf("wd: %v", wd)
			if wd >= int(startDate.Weekday()) {
				nextWeekDay = wd
				break
			}

		}
		if nextWeekDay == 0 {
			nextWeekDay = weekDays[0]
		}
		for {
			startDate = startDate.AddDate(0, 0, 1)
			if startDate.After(now) && int(startDate.Weekday()) == nextWeekDay {
				// !startDate.Before(now) && !startDate.Equal(now)
				break
			}
		}

	case "m":
		return "", errors.New("неподдерживаемый формат")
	// case "m":
	// 	switch {
	// 	case len(repeatParts) == 2 || len(repeatParts) == 3:

	// 		dayArgs := strings.Split(repeatParts[1], ",")
	// 		var dayDone []int
	// 		for _, dayInMonthSTR := range dayArgs {
	// 			dayInMonthINT, err := strconv.Atoi(dayInMonthSTR)
	// 			if err != nil {
	// 				if dayInMonthSTR == "-1" {
	// 					dayInMonthINT = 31
	// 				} else if dayInMonthSTR == "-2" {
	// 					dayInMonthINT = 30
	// 				} else {
	// 					return "", errors.New("неверный формат дня")
	// 				}
	// 			} else if dayInMonthINT < 1 || dayInMonthINT > 31 {
	// 				return "", errors.New("неверный формат дня")
	// 			}
	// 			dayDone = append(dayDone, dayInMonthINT)
	// 		}

	// 		sort.Ints(dayDone)

	// 		var monthDone []int
	// 		if len(repeatParts) == 3 {
	// 			monthArgs := strings.Split(repeatParts[2], ",")
	// 			for _, monthInYearSTR := range monthArgs {
	// 				monthInYearINT, err := strconv.Atoi(monthInYearSTR)
	// 				if err != nil || monthInYearINT < 1 || monthInYearINT > 12 {
	// 					return "", errors.New("неверный формат месяца")
	// 				}
	// 				monthDone = append(monthDone, monthInYearINT)
	// 			}
	// 		}

	// 		sort.Ints(monthDone)

	// 		for _, dayNext := range dayDone {
	// 			for startDate.Day() != dayNext {
	// 				startDate = startDate.AddDate(0, 0, 1)
	// 				if startDate.Month() != now.Month() {
	// 					startDate = time.Date(startDate.Year(), startDate.Month(), dayNext, 0, 0, 0, 0, startDate.Location())
	// 					break
	// 				}
	// 			}

	// 			if startDate.After(now) {
	// 				break
	// 			}
	// 		}

	// 	default:
	// 		return "", errors.New("неподдерживаемый формат")
	// 	}

	// Если первым параметром передан неизвестный формат, то возвращаем ошибку
	default:
		return "", errors.New("неподдерживаемый формат")
	}

	// Возвращаем дату
	return startDate.Format("20060102"), nil

}
