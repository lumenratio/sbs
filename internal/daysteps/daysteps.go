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
		return 0, 0, fmt.Errorf("a string does not have enough data or too big")
	}

	steps, err := strconv.Atoi(parsedStr[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("the step count is zero or less")
	}

	parsedDuration, err := time.ParseDuration(parsedStr[1])
	if err != nil {
		return 0, 0, err
	}
	if parsedDuration <= 0 {
		return 0, 0, fmt.Errorf("wrong time duration")
	}

	return steps, parsedDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, trainDuration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceKM := (float64(steps) * stepLength) / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, trainDuration)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Rewrote to use the Sprintf. But it's looks like a shit: broken formatting, visually ugly.
	result := fmt.Sprintf("Количество шагов: %d.\n"+
		"Дистанция составила %.2f км.\n"+
		"Вы сожгли %.2f ккал.\n", steps, distanceKM, calories)
	return result
}
