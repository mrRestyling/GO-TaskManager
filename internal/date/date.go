package date

import (
	"errors"
	"fmt"
	"sort"
	"start/internal/models"
	"strconv"
	"strings"
	"time"
)

// NextDate - функция для вычисления следующей даты
func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse(models.FormatDate, date)
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", errors.New("нет значения повтора")
	}

	// Разделяем строку повтора на части
	repeatParts := strings.Split(repeat, " ") // w 1,2
	if repeatParts[0] != "y" && len(repeatParts) < 2 {
		return "", errors.New("неправильный формат повтора")
	}

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

	// case "m":
	// 	return "", errors.New("неподдерживаемый формат")
	case "m":

		// switch len(repeatParts) {

		// case 3:
		var dates []time.Time
		var months []int

		if len(repeatParts) == 2 {

			daysSTR := strings.Split(repeatParts[1], ",")

			for _, specificDay := range daysSTR {
				daysINT, err := strconv.Atoi(specificDay)
				if err != nil {
					return "", errors.New("неверный формат дня или недопустимое значение дня")
				} else if daysINT == 31 {
					months = []int{1, 3, 5, 7, 8, 10, 12}
				} else if daysINT == 30 || daysINT == 29 {
					months = []int{1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
				} else {
					months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
				}
			}
			// months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		} else if len(repeatParts) == 3 {

			monthsSTR := strings.Split(repeatParts[2], ",")

			fmt.Printf("monthsSTR: %v\n", monthsSTR)

			months = make([]int, len(monthsSTR))
			for i, month := range monthsSTR {
				months[i], err = strconv.Atoi(month)
				if err != nil || months[i] < 1 || months[i] > 12 {
					return "", errors.New("неверный формат месяца или недопустимое значение месяца")
				}
			}
			fmt.Printf("months: %v\n", months)

		}

		daysSTR := strings.Split(repeatParts[1], ",")
		fmt.Printf("daysSTR: %v\n", daysSTR)

		days := make([]int, len(daysSTR))
		for i, day := range daysSTR {
			days[i], err = strconv.Atoi(day)
			if err != nil {
				return "", errors.New("неверный формат дня или недопустимое значение дня")
			}
			if days[i] < -2 || days[i] > 31 {
				return "", errors.New("неверный формат дня или недопустимое значение дня")
			}
		}

		fmt.Printf("days: %v\n", days)

		for _, month := range months {
			for _, day := range days {

				var date time.Time

				if day < 0 {
					endOfMonth := time.Date(startDate.Year(), time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
					// fmt.Printf("endOfMonth: %v\n", endOfMonth)
					date = endOfMonth.AddDate(0, 0, day+1)
				} else if day >= 1 {
					date = time.Date(startDate.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
				}
				if date.Before(now) {
					date = date.AddDate(1, 0, 0)
				}
				dates = append(dates, date)
			}

		}

		// fmt.Printf("dates: %v\n", dates)

		// Сортируем даты по возрастанию
		sort.Slice(dates, func(i, j int) bool {
			return dates[i].Before(dates[j])
		})

		// Находим следующую дату после текущей
		var firstDateAfterParsedDate time.Time
		for _, date := range dates {
			if date.After(startDate) && date.After(now) { //
				firstDateAfterParsedDate = date
				break
			}
		}

		if !firstDateAfterParsedDate.IsZero() {
			if firstDateAfterParsedDate.After(now) {
				return firstDateAfterParsedDate.Format(models.FormatDate), nil
			} else {
				// fmt.Println("не удалось найти следующую дату после текущей")
				return "", errors.New("не удалось найти следующую дату после текущей")
			}
		} else {
			// fmt.Println("не удалось найти следующую дату после текущей2")
			return "", errors.New("не удалось найти следующую дату после текущей2")
		}

	default:
		return "", errors.New("неподдерживаемый формат")
	}

	// Возвращаем дату
	return startDate.Format(models.FormatDate), nil

}
