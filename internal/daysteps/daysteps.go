package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parsedStr := strings.Split(data, ",")
	if len(parsedStr) != 2 {
		return 0, 0, fmt.Errorf("Daysteps: A string does not have enough data or too big")
	}

	steps, err := strconv.Atoi(parsedStr[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("Daysteps: The step count is zero or less")
	}

	parsedDuration, err := time.ParseDuration(parsedStr[1])
	if parsedDuration <= 0 || err != nil {
		return 0, 0, fmt.Errorf("Wrong time duration")
	}

	return steps, parsedDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	ss := []string{
		"Количество шагов: ",
		"Дистанция составила ",
		"Вы сожгли ",
	}

	steps, training, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceKM := (float64(steps) * stepLength) / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, training)
	if err != nil {
		log.Println(err)
		return ""
	}

	ss[0] += strconv.Itoa(steps) + "."
	ss[1] += strconv.FormatFloat(distanceKM, 'f', 2, 64) + " " + "км."
	ss[2] += strconv.FormatFloat(calories, 'f', 2, 64) + " " + "ккал." + "\n"

	return strings.Join(ss, "\n")
}
