package date

import (
	"errors"
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
		if err != nil || days >= 401 {
			return "", errors.New("недопустимое количество дней или неверный формат числа")
		}

		for {
			startDate = startDate.AddDate(0, 0, days)

			if !startDate.Before(now) && !startDate.Equal(now) {
				break
			}
		}

	case "y":
		for {
			startDate = startDate.AddDate(1, 0, 0)

			if !startDate.Before(now) && !startDate.Equal(now) {
				break
			}
		}

	case "w":
		return "", errors.New("неподдерживаемый формат")

	case "m":
		return "", errors.New("неподдерживаемый формат")

	default:
		return "", errors.New("неподдерживаемый формат")
	}

	return startDate.Format("20060102"), nil

}
