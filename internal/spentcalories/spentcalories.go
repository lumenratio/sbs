package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parsedStr := strings.Split(data, ",")

	if len(parsedStr) != 3 {
		return 0, "", 0, fmt.Errorf("a string does not have enough data or too big")
	}

	steps, err := strconv.Atoi(parsedStr[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("the step count is zero or less")
	}

	parsedDuration, err := time.ParseDuration(parsedStr[2])
	if err != nil {
		return 0, "", 0, err
	} else if parsedDuration <= 0 {
		return 0, "", 0, fmt.Errorf("wrong time duration")
	}

	return steps, parsedStr[1], parsedDuration, nil
}

func distance(steps int, height float64) float64 {
	return (float64(steps) * (height * stepLengthCoefficient)) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / time.Duration(duration).Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case weight <= 0:
		return 0, fmt.Errorf("weight is zero or less")
	case height <= 0:
		return 0, fmt.Errorf("weight is zero or less")
	case steps <= 0:
		return 0, fmt.Errorf("steps is zero or less") // useless check
	case duration <= 0:
		return 0, fmt.Errorf("duration is zero or less") // useless check again
	}

	return (weight * meanSpeed(steps, height, duration) * time.Duration(duration).Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case weight <= 0:
		return 0, fmt.Errorf("weight is zero or less")
	case height <= 0:
		return 0, fmt.Errorf("weight is zero or less")
	case steps <= 0:
		return 0, fmt.Errorf("steps is zero or less") // useless check
	case duration <= 0:
		return 0, fmt.Errorf("duration is zero or less") // useless check again
	}
	return ((weight * meanSpeed(steps, height, duration) * time.Duration(duration).Minutes()) / minInH) * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, workout, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var (
		resultCalories float64
		resultDuration float64
		resultDistance float64
		resultSpeed    float64
	)

	switch workout {
	case "Ходьба":
		resultCalories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		resultCalories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	resultDuration = time.Duration(duration).Hours()
	resultDistance = distance(steps, height)
	resultSpeed = meanSpeed(steps, height, duration)

	// Rewrote to use the Sprintf. It's ugly looks. Just ugly.
	result := fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`, workout, resultDuration, resultDistance, resultSpeed, resultCalories)

	return result, nil
}
