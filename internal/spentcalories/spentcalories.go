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

	if len(parsedStr) < 3 {
		return 0, "", 0, fmt.Errorf("Spentclories: A string does not have enough data")
	}

	steps, err := strconv.Atoi(parsedStr[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Spentclories: The step count is zero or less")
	}

	parsedDuration, err := time.ParseDuration(parsedStr[2])
	if err != nil {
		return 0, "", 0, err
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
	if steps <= 0 || duration <= 0 {
		return 0, fmt.Errorf("Some parameters if wrong. Parameters must be bigger than zero")
	}
	return (weight * meanSpeed(steps, height, duration) * time.Duration(duration).Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || duration <= 0 {
		return 0, fmt.Errorf("Some parameters if wrong. Parameters must be bigger than zero")
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

	resultSS := []string{"Тип тренировки: ",
		"Длительность: ",
		"Дистанция: ",
		"Скорость: ",
		"Сожгли калорий: "}

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

	resultSS[0] += workout
	resultSS[1] += strconv.FormatFloat(resultDuration, 'f', 2, 64) + " " + "ч."
	resultSS[2] += strconv.FormatFloat(resultDistance, 'f', 2, 64) + " " + "км."
	resultSS[3] += strconv.FormatFloat(resultSpeed, 'f', 2, 64) + " " + "км/ч"
	resultSS[4] += strconv.FormatFloat(resultCalories, 'f', 2, 64)

	return strings.Join(resultSS, "\n"), nil
}
