package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Value  int `json:"value"`
	Weight int `json:"weight"`
}
type Knapsack struct {
	Weight int    `json:"weight"`
	Items  []Item `json:"items"`
}

var knapSack Knapsack
var tmpl *template.Template

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}

}

func knapSackParallel(W int, weight []int, value []int, n int) int {

	var Kparallel [4][51]int

	//Kparallel := make([][]int, n+1, W+1)

	var wg sync.WaitGroup
	//var m sync.Mutex
	//var Kparallel [10][51]int

	//wg.Add(n + 1)

	for i := 0; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			for w := 0; w <= W; w++ {

				if i == 0 || w == 0 {

					Kparallel[i][w] = 0

				} else if weight[i-1] <= w {

					Kparallel[i][w] = max(value[i-1]+Kparallel[i-1][w-weight[i-1]], Kparallel[i-1][w])

				} else {

					Kparallel[i][w] = Kparallel[i-1][w]

				}

			}
			wg.Done()
		}(i)
		wg.Wait()

	}
	return Kparallel[n][W]

}
func knapSackClassic(W int, weight []int, value []int, n int) int {

	var K [4][51]int

	for i := 0; i <= n; i++ {

		for w := 0; w <= W; w++ {

			if i == 0 || w == 0 {
				K[i][w] = 0
			} else if weight[i-1] <= w {
				K[i][w] = max(value[i-1]+K[i-1][w-weight[i-1]], K[i-1][w])
			} else {
				K[i][w] = K[i-1][w]
			}
		}
	}

	return K[n][W]

}
func Demo() {
	coutnOfElement := 100
	W := 20

	value := make([]int, coutnOfElement)
	weight := make([]int, coutnOfElement)

	rand.Seed(time.Now().UnixMilli())
	for i := range value {
		value[i] = rand.Intn(100)
		weight[i] = rand.Intn(100)
	}

	timeBefore := time.Now()
	log.Println("paraller result", knapSackParallel(W, weight, value, coutnOfElement))
	timeAfter := time.Now()
	result := timeAfter.UnixMilli() - timeBefore.UnixMilli()
	log.Println("Time after parallel ", result)

	timeBeforeC := time.Now()
	log.Println("non parallel result", knapSackClassic(W, weight, value, coutnOfElement))
	timeAfterC := time.Now()
	result1 := timeAfterC.UnixMilli() - timeBeforeC.UnixMilli()
	log.Println("Time after classic", result1)

}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/index.html"))
}
func Handle() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		coutnOfElement := 100
		//W := 20

		value := make([]int, coutnOfElement)
		weight := make([]int, coutnOfElement)
		rand.Seed(time.Now().UnixMilli())
		for i := range value {
			value[i] = rand.Intn(100)
			weight[i] = rand.Intn(100)
		}
		//result := knapSackClassic(W, weight, value, coutnOfElement)

		tmpl.ExecuteTemplate(w, "index.html", nil)

	}
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", index)
	mux.HandleFunc("/calculate", calculate)
	//mux.HandleFunc("/css", staticHandler)
}

func calculate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		err1 := r.ParseForm()
		if err1 != nil {
			fmt.Print(err1)
		}
		fmt.Print(r.Form)

		values := strings.Split(r.FormValue("values"), " ")
		weights := strings.Split(r.FormValue("weights"), " ")
		maxCapacity, err := strconv.Atoi(r.FormValue("weight"))
		if err != nil {
			fmt.Print(err)
		}

		log.Print(values, weights, maxCapacity)

		valuesSlice := make([]int, 0, len(values))
		weightsSlice := make([]int, 0, len(values))

		for i := range values {
			currentValue, err := strconv.Atoi(values[i])
			if err != nil {
				fmt.Print(err)
			}
			valuesSlice = append(valuesSlice, currentValue)
			currentWeight, err := strconv.Atoi(weights[i])
			if err != nil {
				fmt.Print(err)
			}
			weightsSlice = append(weightsSlice, currentWeight)
		}

		result := knapSackClassic(maxCapacity, weightsSlice, valuesSlice, len(valuesSlice))

		// knapSack.Weight, err = strconv.Atoi(maxCapacity)
		// if err != nil {
		// 	fmt.Print(err)
		// }

		// items := make([]Item, 0)
		// currentItem := Item{}

		// for i := range valuesSlice {
		// 	currentItem.Value, err = strconv.Atoi(valuesSlice[i])
		// 	currentItem.Weight, err = strconv.Atoi(weightsSlice[i])

		// 	items = append(items, currentItem)
		// }

		// knapSack.Items = items

		log.Println("result is:", result)
		w.Write([]byte(strconv.Itoa(result)))
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Print(err)
	}
	err = tmpl.Execute(w, knapSack)
	if err != nil {
		fmt.Print(err)
	}
}

func main() {

	var mux = http.NewServeMux()
	fs := http.FileServer(http.Dir("resources"))
	mux.Handle("/resources/", http.StripPrefix("/resources", fs))
	registerRoutes(mux)
	httpServer := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}

}
