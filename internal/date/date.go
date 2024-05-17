package date

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Функция для вычисления следующей даты
func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("нет значения повтора")
	}

	repeatParts := strings.Split(repeat, " ")

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
	// 	if len(repeatParts) < 2 {
	// 		return "", errors.New("недостаточно аргументов для формата месяца")
	// 	}

	// 	dayArgs := strings.Split(repeatParts[0], ",")
	// 	monthArgs := strings.Split(repeatParts[1], ",")

	// 	if len(dayArgs) == 0 || len(monthArgs) == 0 {
	// 		return "", errors.New("недопустимый формат дня или месяца")
	// 	}

	// 	for _, dayStr := range dayArgs {
	// 		day, err := strconv.Atoi(dayStr)
	// 		if err != nil {
	// 			return "", errors.New("неверный формат дня месяца")
	// 		}

	// 		if day == -1 {
	// 			// Handle last day of the month
	// 		} else if day == -2 {
	// 			// Handle second to last day of the month
	// 		} else if day < 1 || day > 31 {
	// 			return "", errors.New("недопустимое значение дня месяца")
	// 		}
	// 	}

	// 	for _, monthStr := range monthArgs {
	// 		month, err := strconv.Atoi(monthStr)
	// 		if err != nil || month < 1 || month > 12 {
	// 			return "", errors.New("недопустимое значение месяца")
	// 		}
	// }

	default:
		return "", errors.New("неподдерживаемый формат")
	}

	return startDate.Format("20060102"), nil

}
