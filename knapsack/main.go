package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Item struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Weight int    `json:"weight"`
}
type Knapsack struct {
	Capacity int    `json:"capacity"`
	Profit   int    `json:"profit"`
	Items    []Item `json:"items"`
}
type Items struct {
	Items     []Item
	Kparallel [][]int
}

var items Items
var knapSack Knapsack
var tmpl *template.Template

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func knapSackParallel(W int, weight []int, value []int, n int) (int, []Item) {
	var wg sync.WaitGroup

	items.Kparallel = make([][]int, n+1)
	for i := 0; i <= n; i++ {
		items.Kparallel[i] = make([]int, W+1)
	}

	for i := 0; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			for w := 0; w <= W; w++ {
				if i == 0 || w == 0 {
					items.Kparallel[i][w] = 0
				} else if weight[i-1] <= w {
					items.Kparallel[i][w] = max(value[i-1]+items.Kparallel[i-1][w-weight[i-1]], items.Kparallel[i-1][w])
				} else {
					items.Kparallel[i][w] = items.Kparallel[i-1][w]
				}
			}
			wg.Done()
		}(i)
		wg.Wait()

	}
	result := items.Kparallel[n][W]
	addeditems := fetcAddedItems(result, n, W, value, weight)

	return result, addeditems

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

func setItemsSlice(values []int, weights []int) []Item {
	var items []Item
	for i := 0; i < len(values); i++ {

		item := Item{Name: string(rune(i + 65)), Value: values[i], Weight: weights[i]}
		items = append(items, item)
	}

	return items
}
func fetcAddedItems(result int, n int, w int, value []int, weight []int) []Item {
	var addedItemsWeights = make([]int, 0)
	var addedItems = make([]Item, 0)

	for i := n; i > 0 && result > 0; i-- {
		if result != items.Kparallel[i-1][w] {
			addedItemsWeights = append(addedItemsWeights, weight[i-1])
			result = result - value[i-1]
			w = w - weight[i-1]
		}
	}

	for _, w := range addedItemsWeights {
		for _, item := range items.Items {
			if item.Weight == w {
				addedItems = append(addedItems, item)
			}
		}
	}

	return addedItems
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", HandlerInputForm)
	mux.HandleFunc("/calculate", calculate)
	//mux.HandleFunc("/css", staticHandler)
}

func calculate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Print(err)
		}
		values := strings.Split(r.FormValue("values"), " ")
		weights := strings.Split(r.FormValue("weights"), " ")
		knapSack.Capacity, err = strconv.Atoi(r.FormValue("weight"))
		if err != nil {
			fmt.Print(err)
		}

		log.Print(values, weights, knapSack.Capacity)

		valuesSlice := make([]int, 0, len(values))
		weightsSlice := make([]int, 0, len(values))

		for i := range values {
			currentValue, err := strconv.Atoi(values[i])
			if err != nil {
				fmt.Print(err)
			}
			currentWeight, err := strconv.Atoi(weights[i])
			if err != nil {
				fmt.Print(err)
			}

			valuesSlice = append(valuesSlice, currentValue)
			weightsSlice = append(weightsSlice, currentWeight)
		}

		items.Items = setItemsSlice(valuesSlice, weightsSlice)
		log.Println("items:", items)

		knapSack.Profit, knapSack.Items = knapSackParallel(knapSack.Capacity, weightsSlice, valuesSlice, len(valuesSlice))

		log.Println("profit is:", knapSack.Profit)
		w.Write([]byte(strconv.Itoa(knapSack.Profit)))

		log.Println("added items are:", knapSack.Items)

	}
}

func HandlerInputForm(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Print(err)
	}
	err = tmpl.Execute(w, knapSack)
	if err != nil {
		fmt.Print(err)
	}
}
func HandlerCreateTables(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/table.html")
	if err != nil {
		fmt.Print(err)
	}
	err = tmpl.Execute(w, items)
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
