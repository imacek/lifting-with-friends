package internal

import (
	"math"
	"sort"
	"time"
)

var AnalysisCutoffDate = time.Date(2021, 1, 1, 0, 0, 0, 0, LosAngelesTimeLocation)

type ExerciseAggData struct {
	Timestamp    time.Time `json:"timestamp"`
	MaxWeight    float64   `json:"maxWeight"`
	MaxOneRepMax float64   `json:"maxOneRepMax"`
	TotalVolume  float64   `json:"totalVolume"`
}

type UserExerciseTimeSeries = map[string][]ExerciseAggData

func calculateExerciseTimeSeries(liftingSets []LiftingSet, timeKeyFunc func(time.Time) time.Time) UserExerciseTimeSeries {
	m := make(map[string]map[time.Time]ExerciseAggData)

	for _, ls := range liftingSets {
		timeKey := timeKeyFunc(ls.timestamp)
		if timeKey.Before(AnalysisCutoffDate) {
			continue
		}

		if _, contains := m[ls.exerciseName]; !contains {
			m[ls.exerciseName] = make(map[time.Time]ExerciseAggData)
		}
		if _, contains := m[ls.exerciseName][timeKey]; !contains {
			m[ls.exerciseName][timeKey] = ExerciseAggData{
				Timestamp: timeKey,
			}
		}

		data := m[ls.exerciseName][timeKey]
		m[ls.exerciseName][timeKey] = ExerciseAggData{
			Timestamp:    timeKey,
			MaxWeight:    math.Max(data.MaxWeight, ls.weight),
			MaxOneRepMax: math.Max(data.MaxOneRepMax, ls.oneRepMax),
			TotalVolume:  data.TotalVolume + ls.weight*float64(ls.reps),
		}
	}

	// Drop the map
	m2 := make(map[string][]ExerciseAggData, len(m))

	for user, dataMap := range m {
		m2[user] = make([]ExerciseAggData, len(dataMap))

		index := 0
		for _, data := range dataMap {
			m2[user][index] = data
			index++
		}

		sort.Slice(m2[user], func(i, j int) bool {
			return m2[user][i].Timestamp.Before(m2[user][j].Timestamp)
		})
	}

	return m2
}

// Identity function
func timeToTime(t time.Time) time.Time {
	return t
}
func timeToDateStart(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
func timeToWeekStart(t time.Time) time.Time {
	year, month, day := t.AddDate(0, 0, -int(t.Weekday())).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
func timeToMonthStart(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func CalculateUserExerciseTimeSeries(userLiftingSets map[string][]LiftingSet) map[string][4]UserExerciseTimeSeries {
	result := make(map[string][4]UserExerciseTimeSeries, len(userLiftingSets))

	for user, ls := range userLiftingSets {
		result[user] = [4]UserExerciseTimeSeries{
			calculateExerciseTimeSeries(ls, timeToTime),
			calculateExerciseTimeSeries(ls, timeToDateStart),
			calculateExerciseTimeSeries(ls, timeToWeekStart),
			calculateExerciseTimeSeries(ls, timeToMonthStart),
		}
	}

	return result
}
