package main

import (
	"fmt"
	"net/http"
	"time"
)

func renderFormPage(w http.ResponseWriter, r *http.Request) {
	_, currentWeek := time.Now().ISOWeek()

	if r.Method == http.MethodPost {
		metabolism := r.FormValue("metabolism")
		calories := r.FormValue("calories")

		dayOfWeek := time.Now().Weekday().String()

		if metabolism != "" && !isBasalMetabolismSet(currentWeek) {
			var metab float64
			fmt.Sscanf(metabolism, "%f", &metab)
			insertBasalMetabolism(metab, currentWeek)
		}

		if calories != "" {
			var cal float64
			fmt.Sscanf(calories, "%f", &cal)
			addCalories(dayOfWeek, cal, currentWeek)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	metabolismDisabled := ""
	if isBasalMetabolismSet(currentWeek) {
		metabolismDisabled = "disabled"
	}

	fmt.Fprintf(w, `
		<html>
		<head>
			<meta charset="utf-8">
			<title>Consumo Calórico</title>
		</head>
		<body>
			<h1>Registro de Consumo Calórico</h1>
			<form method="POST" action="/">
				<label for="metabolism">Metabolismo Basal (em calorias):</label><br>
				<input type="text" id="metabolism" name="metabolism" %s><br><br>

				<label for="calories">Calorias Consumidas:</label><br>
				<input type="text" id="calories" name="calories"><br><br>
				<input type="submit" value="Salvar">
			</form>
			<br>
			<a href="./weekly_calories.html" onclick="event.preventDefault(); renderChart();">Ver gráfico semanal</a>
			<script>
				function renderChart() {
					fetch("/generate_chart").then(response => {
						if (response.ok) {
							alert("Gráfico gerado!");
							window.open('./weekly_calories.html');
						} else {
							alert("Erro ao gerar gráfico.");
						}
					});
				}
			</script>
		</body>
		</html>
	`, metabolismDisabled)
}

func generateChartHandler(w http.ResponseWriter, r *http.Request) {
	_, currentWeek := time.Now().ISOWeek()
	renderWeeklyChart(w, currentWeek)
	fmt.Fprint(w, "Gráfico gerado.")
}
