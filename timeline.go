package main

import (
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Data struct {
	Tasks                            []*Task   `json:"tasks"`
	MySettings                       *Settings `json:"settings"`
	MyTheme                          *Theme    `json:"theme"`
	First, Last                      time.Time
	Days                             int
	FontSize, Scale                  float64
	W, H, LabelW, ChartW, DayW, RowH float64
	FrameBorderColor                 color.Color
	FrameFillColor                   color.Color
	StripeColorDark                  color.Color
	StripeColorLight                 color.Color
	GridColor                        color.Color
}

type Task struct {
	Start              string   `json:"start"`
	End                string   `json:"end"`
	Label              string   `json:"label"`
	Recur              int      `json:"recur"`
	Milestones         []string `json:"milestones"`
	DateStamps         []string `json:"dateStamps"`
	StartTo            []int    `json:"startTo"`
	EndTo              []int    `json:"endTo"`
	StartTime, EndTime time.Time
	BorderColor        color.Color
	FillColor          color.Color
}

type Theme struct {
	ColorScheme      string `json:"colorScheme"`
	BorderColor1     string `json:"borderColor1"`
	FillColor1       string `json:"fillColor1"`
	BorderColor2     string `json:"borderColor2"`
	FillColor2       string `json:"fillColor2"`
	FrameBorderColor string `json:"frameBorderColor"`
	FrameFillColor   string `json:"frameFillColor"`
	StripeColorDark  string `json:"stripeColorDark"`
	StripeColorLight string `json:"stripeColorLight"`
	GridColor        string `json:"gridColor"`
}

type Settings struct {
	Lang          string `json:"lang"`
	End           string `json:"end"`
	Zoom          int    `json:"zoom"`
	HideDaysFrom  int    `json:"hideDaysFrom"`
	HideWeeksFrom int    `json:"hideWeeksFrom"`
}

type Result struct {
	Message string
	Code    int
	Context *gg.Context
}

type ShortResult struct {
	Message string
	Code    int
}

type Locale struct {
	Lang   string     `json:"lang"`
	Layout string     `json:"layout"`
	Months [12]string `json:"months"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: ./%s [<JSON file> [<JSON file>]]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}
	port := flag.Int("p", 8000, "listen on port")
	hostname := flag.String("n", "localhost", "Hostname")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		serve(*hostname, *port)
		return
	}

	ch := make(chan ShortResult)

	for _, input := range args {
		go processFile(input, ch)
	}

	var mu sync.Mutex
	var code int
	for range args {
		result := <-ch
		mu.Lock()
		code += result.Code
		mu.Unlock()
		fmt.Println(result.Message)
	}
	os.Exit(code)
}
