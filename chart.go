package main

import (
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func renderWeeklyChart(w http.ResponseWriter, currentWeek int) {
	weeklyCalories, basalMetabolism := getWeeklyCalories(currentWeek)

	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	caloriesData := []opts.LineData{}
	balanceData := []opts.LineData{}

	for _, day := range days {
		calories := weeklyCalories[day]
		caloriesData = append(caloriesData, opts.LineData{Value: calories})

		balance := basalMetabolism - calories
		balanceData = append(balanceData, opts.LineData{Value: balance})
	}

	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico de calorias"}))
	line.SetXAxis(days).
		AddSeries("Calorias Consumidas", caloriesData).
		AddSeries("Saldo Calórico", balanceData)

	f, _ := os.Create("./weekly_calories.html")
	defer f.Close()
	line.Render(f)
}
