package main

import (
	"time"
	"fmt"
	"math"
)

// Алгоритм управляет работой канального вентилятора
// Алгоритм запускается на периодической основе через равные интервалы времени.
// В момент запуска анализируются доступные на этот момент данные датчиков.
//
// Основной алгоритм использует данные датчиков во влажных помещениях и в сухих помещениях.
// Наибольшее значение влажности во влажных помещениях сравнивается с наибольшим значением влажности
// в сухих помещениях. При разнице влажности более 10% алгоритм включает канальный вентилятор.
//
// Алгоритм продолжает контролировать разницу влажностей. При разнице влажности менее 5% алгоритм
// выключает вентилятор. Дополнительно контролируется темп падения влажности во влажных помещениях.
// Если влажность во влажном помещении не изменяется (с допуском +-3%) в течение 30 минут, то
// алгоритм выключает вентилятор и прекращает контроль на другие 30 минут. По истечении этого
// времени контроль возобновляется в обычном режиме.
//
// В любом случае (даже в последнем) алгоритм продолжает публиковать управляющие сигналы
// получатель использует это для измеения пульса алгоритма, чтобы понять, что сервис, управляющий
// вентилятором не завис и продолжает работать. Получатель должен уметь работать с повторяющимися
// управляющими сигналами.
// 
// Также, могут возникнуть ситуации, когда не поступают сигналы от датчиков. Если "умерли" все датчики
// сухих помещений, то в качестве условной влажности в сухих помещениях принимается 60%.
//
// Если "умерли" все датчики влажных помещений, то алгоритм выключает канальный вентилятор. При появлении
// сигнала хотя бы от одного датчика во влажном помещении возобновляется обработка данных с основного
// алгоритма. И опять: пересылка сигналов отключения вентилятора продолжается на периодической основе,
// чтобы потребитель мог измерить пульс работы сервиса.


const calcTimeout = time.Second * 30
var calcTick *time.Ticker
var savedHumidity float64
var savedTime time.Time
var startSteadyHumidityControl = false
var stopControl = false

func calcLoop() {
	for range calcTick.C {
		fmt.Println("Calc Tick")
		doCalc()
	}
}

func doCalc() {
	// Если мокрые датчики умерли - то вообще нечего контролировать
	if !wetSensorsAlive {
		ventServiceState.off()
		return
	}
	// Контроль может быть остановлен на 30 минут, если влажность не менялась, не смотря на работу
	// вентилятора в течение 30 минут
	if stopControl && time.Since(savedTime) < time.Minute * 30 {
		ventServiceState.off()
	} else if stopControl {
		// Время прекращения контроля закончилось
		stopControl = false
	}
	
	wet := maxHumidity(wetSensors)
	dry := 60.0
	if drySensorsAlive {
		dry = maxHumidity(drySensors)
	}

	if wet - dry > 10 {
		if !startSteadyHumidityControl {
			// Начинаем контроль показания влажности - оно должно меняться
			startSteadyHumidityControl = true
			savedHumidity = wet
			savedTime = time.Now()
		} else if math.Abs(savedHumidity - wet) > 3.0 {
			// Если влажность обновилась больше, чем на 3% - то обновим опорное значение
			savedHumidity = wet
			savedTime = time.Now()
		} else if time.Since(savedTime) > time.Minute * 30 {
			// Если с момента начала контроля в течение 30 минут значение не изменилось
			// более, чем на 3%, то выключаем вентилятор и прекращаем контроль на другие 30 минут
			startSteadyHumidityControl = false
			stopControl = true
			savedTime = time.Now()
			ventServiceState.off()
			return
		}
		ventServiceState.on()
	}
	if wet - dry < 5 {
		startSteadyHumidityControl = false
		ventServiceState.off()
	}
}

func maxHumidity(arr []DHTSensorData) float64 {
	max := -1.
	for _, d := range arr {
		if d.Humidity > max {
			max = d.Humidity
		}
	}
	return max
}

func minHumidity(arr []DHTSensorData) float64 {
	min := 101.
	for _, d := range arr {
		if d.Humidity < min {
			min = d.Humidity
		}
	}
	return min
}
