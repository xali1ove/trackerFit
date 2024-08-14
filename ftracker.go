package ftracker

import (
	"fmt"
)

const (
	lenStep   = 0.65  // средняя длина шага
	mInKm     = 1000  // количество метров в километре
	minInH    = 60    // количество минут в часе
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с
	cmInM     = 100   // количество сантиметров в метре
)

func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	distance := distance(action)
	return distance / duration
}

func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	var dist float64
	var speed float64
	var calories float64

	switch trainingType {
	case "Бег":
		dist = distance(action)
		speed = meanSpeed(action, duration)
		calories = RunningSpentCalories(action, weight, duration)
	case "Ходьба":
		dist = distance(action)
		speed = meanSpeed(action, duration)
		calories = WalkingSpentCalories(action, duration, weight, height)
	case "Плавание":
		dist = float64(lengthPool*countPool) / mInKm // корректировка дистанции для плавания
		speed = dist / duration
		calories = SwimmingSpentCalories(lengthPool, countPool, duration, weight)
	default:
		return "неизвестный тип тренировки"
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, dist, speed, calories)
}

const (
	runningCaloriesMeanSpeedMultiplier = 18   // множитель средней скорости
	runningCaloriesMeanSpeedShift      = 1.79 // среднее количество сжигаемых калорий при беге
)

func RunningSpentCalories(action int, weight, duration float64) float64 {
	speed := meanSpeed(action, duration)
	return ((runningCaloriesMeanSpeedMultiplier * speed * runningCaloriesMeanSpeedShift) * weight / mInKm * duration * minInH)
}

const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста
)

func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
	speed := meanSpeed(action, duration) * kmhInMsec
	return ((walkingCaloriesWeightMultiplier*weight + (speed*speed/height)*walkingSpeedHeightMultiplier*weight) * duration * minInH)
}

const (
	swimmingCaloriesMeanSpeedShift   = 1.1 // среднее количество сжигаемых калорий при плавании относительно скорости
	swimmingCaloriesWeightMultiplier = 2   // множитель веса при плавании
)

func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	return float64(lengthPool) * float64(countPool) / mInKm / duration
}

func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	speed := swimmingMeanSpeed(lengthPool, countPool, duration)
	return (speed + swimmingCaloriesMeanSpeedShift) * swimmingCaloriesWeightMultiplier * weight * duration
}
