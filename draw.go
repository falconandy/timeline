package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"strconv"
	"time"
)

func drawDropShadow(gc *draw2dimg.GraphicContext, x1, y1, x2, y2 float64) {
	var opacity uint8
	opacity = 0x44
	offset := 4.0

	gc.Save()
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, opacity})
	gc.MoveTo(x1+offset, y1+offset)
	gc.LineTo(x2+offset, y1+offset)
	gc.LineTo(x2+offset, y2+offset)
	gc.LineTo(x1+offset, y2+offset)
	gc.Close()
	gc.Fill()
	gc.Restore()
}

func drawBlock(d *Data, gc *draw2dimg.GraphicContext, x1, y1, x2, y2 float64, strokeColor, fillColor color.Color, label string, showLabel bool) {
	gc.Save()
	gc.BeginPath()
	gc.SetFillColor(fillColor)
	gc.SetStrokeColor(strokeColor)
	gc.SetLineWidth(1)
	gc.MoveTo(x1, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x1, y2)
	gc.Close()
	gc.FillStroke()
	gc.Fill()
	gc.Restore()

	if showLabel == false {
		return
	}

	//label
	gc.Save()
	gc.SetFillColor(d.FrameBorderColor)
	tx1, ty1, tx2, ty2 := gc.GetStringBounds(label)
	w := x2 - x1
	tw := tx2 - tx1
	th := ty2 - ty1

	adjustedTextY := y2 - th*0.5
	//shift or hide if label too long
	if tw > w || x1+tw > d.ChartW {
		label = ""
	}

	midX := x1 + (x2-x1)*0.5
	adjustedTextX := midX - tw*0.5

	//normalize label x-position
	if adjustedTextX < 0 {
		adjustedTextX = 0
	} else {
		for adjustedTextX+tw > d.ChartW {
			adjustedTextX--
		}
	}

	gc.FillStringAt(label, adjustedTextX, adjustedTextY)
	gc.Restore()
}
func drawArrowHeadN(d *Data, gc *draw2dimg.GraphicContext, x, y float64) {
	gc.Save()
	gc.BeginPath()
	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	var r float64
	r = 3.0 * d.Scale
	gc.MoveTo(x-r, y+r)
	gc.LineTo(x, y+r-r/2)
	gc.LineTo(x+r, y+r)
	gc.LineTo(x, y-r)
	gc.Close()
	gc.FillStroke()
	gc.Fill()
	gc.Restore()
}
func drawArrowHeadE(d *Data, gc *draw2dimg.GraphicContext, x, y float64) {
	gc.Save()
	gc.BeginPath()
	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	var r float64
	r = 3.0 * d.Scale
	gc.MoveTo(x, y)
	gc.LineTo(x-2*r, y-r)
	gc.LineTo(x-r, y)
	gc.LineTo(x-2*r, y+r)
	gc.Close()
	gc.FillStroke()
	gc.Fill()
	gc.Restore()
}
func drawArrowHeadS(d *Data, gc *draw2dimg.GraphicContext, x, y float64) {
	gc.Save()
	gc.BeginPath()
	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	var r float64
	r = 3.0 * d.Scale
	gc.MoveTo(x, y)
	gc.LineTo(x-r, y-2*r)
	gc.LineTo(x, y-r)
	gc.LineTo(x+r, y-2*r)
	gc.Close()
	gc.FillStroke()
	gc.Fill()
	gc.Restore()
}

func drawMilestone(d *Data, gc *draw2dimg.GraphicContext, x, y float64) {
	gc.Save()
	gc.BeginPath()
	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	var r float64
	r = 4 * d.Scale
	gc.MoveTo(x-r, y)
	gc.LineTo(x, y-2*r)
	gc.LineTo(x+r, y)
	gc.LineTo(x, y+2*r)
	gc.Close()
	gc.FillStroke()
	gc.Fill()
	gc.Restore()
}

func drawDateStamp(d *Data, gc *draw2dimg.GraphicContext, x, y float64, label string) {
	//vertical line
	gc.Save()
	gc.BeginPath()
	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	var r float64
	r = 6 * d.Scale
	gc.MoveTo(x, y-2*r)
	gc.LineTo(x, y+2*r)
	gc.FillStroke()
	gc.Restore()

	//label
	gc.Save()
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	label = label[5:]
	x1, y1, x2, y2 := gc.GetStringBounds(label)
	w := x2 - x1
	h := y2 - y1

	adjustedTextY := y + 2.5*h
	adjustedTextX := x - w*0.5

	//normalize label x-position
	if adjustedTextX < 0 {
		adjustedTextX = 0
	} else {
		for adjustedTextX+w > d.ChartW {
			adjustedTextX--
		}
	}

	//strip out year as it's shown at the top
	gc.FillStringAt(label, adjustedTextX, adjustedTextY)
	gc.Restore()
}

func drawCalendarGuides(d *Data, gc *draw2dimg.GraphicContext, y1, y2 float64, fn func(time.Time) string) {
	var period string

	gc.Save()
	gc.SetStrokeColor(d.CanvasGridColor)

	for i := 0; i <= d.Days; i++ {
		t := time.Date(d.First.Year(), d.First.Month(), d.First.Day()+i, 0, 0, 0, 0, time.UTC)
		currentPeriod := fn(t)

		if i == 0 {
			prevT := time.Date(d.First.Year(), d.First.Month(), d.First.Day()-1, 0, 0, 0, 0, time.UTC)
			period = fn(prevT)
		}
		if currentPeriod != period {
			x := float64(i) * d.DayW
			gc.Save()
			gc.BeginPath()
			a := []float64{2.0, 2.0}
			gc.SetLineDash(a, 0.0)
			gc.MoveTo(x, y1)
			gc.LineTo(x, float64(y2))
			gc.Stroke()
			gc.Restore()
			period = currentPeriod
		}
	}

	gc.Restore()
}

func drawCalendarRow(d *Data, gc *draw2dimg.GraphicContext, y float64, strokeColor, fillColor color.Color, fn func(time.Time) string) {
	var period string
	var from int
	for i := 0; i <= d.Days; i++ {
		t := time.Date(d.First.Year(), d.First.Month(), d.First.Day()+i, 0, 0, 0, 0, time.UTC)
		currentPeriod := fn(t)

		last := i == d.Days

		if i == 0 {
			prevT := time.Date(d.First.Year(), d.First.Month(), d.First.Day()-1, 0, 0, 0, 0, time.UTC)
			period = fn(prevT)
		}

		// new period or end of timeline
		if currentPeriod != period || last {
			x1 := float64(from) * d.DayW
			x2 := float64(i) * d.DayW
			if last {
				x2 += d.DayW
			}
			y2 := y + d.RowH

			//TODO: debug points here
			drawBlock(d, gc, x1, y, x2, y2, strokeColor, fillColor, period, true)

			// now update from for next section
			period = currentPeriod
			from = i
		}
	}
}

func drawStripe(d *Data, gc *draw2dimg.GraphicContext, index int, y1, y2 float64) {
	color := d.CanvasColor1
	if index%2 != 0 {
		color = d.CanvasColor2
	}

	y1 -= d.RowH / 2
	y2 += d.RowH / 2

	gc.Save()
	gc.SetStrokeColor(color)
	gc.SetFillColor(color)
	gc.BeginPath()
	gc.MoveTo(0, y1)
	gc.LineTo(d.ChartW, y1)
	gc.LineTo(d.ChartW, y2)
	gc.LineTo(0, y2)
	gc.Close()
	gc.FillStroke()
	gc.Restore()
}

func drawScene(d *Data, path string) {
	w, h, rowH := d.W, d.H, d.RowH

	dest := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetFontSize(d.FontSize) //preserve scale

	// calendar functions
	fnYear := func(t time.Time) string {
		return strconv.Itoa(t.Year())
	}
	fnMonth := func(t time.Time) string {
		return t.Month().String()
	}
	fnWeek := func(t time.Time) string {
		_, week := t.ISOWeek()
		return strconv.Itoa(week)
	}

	// label block
	var maxLabelWidth float64
	for _, task := range d.Tasks {
		label := task.Label
		x1, _, x2, _ := gc.GetStringBounds(label)
		labelWidth := x2 - x1
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}

	d.LabelW = maxLabelWidth + 5 //leave room at start for milestones
	d.ChartW = w - d.LabelW - 5  //leave room at end for milestones

	//var dayW float64
	d.DayW = d.ChartW / float64(d.Days)
	var x, y float64

	// guides
	var fnGuide func(t time.Time) string
	if d.Days < d.LayoutSteps[0] {
		fnGuide = fnWeek
	} else if d.Days < d.LayoutSteps[1] {
		fnGuide = fnMonth
	} else {
		fnGuide = fnYear
	}

	// increment y as needed
	y = 0

	gc.Save()
	gc.Translate(d.LabelW, 0)
	// year
	drawCalendarRow(d, gc, y, d.FrameBorderColor, d.FrameFillColor, fnYear)
	// month
	y += rowH
	drawCalendarRow(d, gc, y, d.FrameBorderColor, d.FrameFillColor, fnMonth)

	// weeks
	if d.Days < d.LayoutSteps[1] {
		y += rowH
		drawCalendarRow(d, gc, y, d.FrameBorderColor, d.FrameFillColor, fnWeek)
	}

	// days
	if d.Days < d.LayoutSteps[0] {
		y += rowH
		var weekend bool
		for i := 0; i <= d.Days; i++ {
			//determine if weekend
			t := time.Date(d.First.Year(), d.First.Month(), d.First.Day()+i, 0, 0, 0, 0, time.UTC)
			weekend = t.Weekday() == 0 || t.Weekday() == 6
			shade := d.CanvasColor2
			if weekend {
				shade = d.CanvasColor1
			}

			x = float64(i) * d.DayW

			drawBlock(d, gc, x, y, x+d.DayW, y+rowH, d.CanvasGridColor, shade, "", true)
		}
	}

	// stripes
	stripeY := y + 1.5*rowH + 1.0
	for index, _ := range d.Tasks {
		drawStripe(d, gc, index, stripeY, stripeY+d.RowH)
		stripeY += rowH * 2
	}

	// draw guides from calendar block onwards
	y += rowH
	drawCalendarGuides(d, gc, y, y+float64(len(d.Tasks))*2.0*rowH+2*rowH, fnGuide)

	gc.Restore()

	y += d.RowH / 2

	//iterate over tasks
	for index, task := range d.Tasks {
		//y = 4.5*rowH + float64(index)*rowH*2 + 1.0

		start := dayIndex(task.StartTime, d.First, d.Last)
		end := dayIndex(task.EndTime, d.First, d.Last)

		if start == -1 || end == -1 {
			//can't place task on timeline
			continue
		}

		//one-day tasks: draw full day
		if start == end {
			end++
		}

		x1 := float64(start) * d.DayW
		x2 := float64(end) * d.DayW
		y1 := y
		y2 := y + rowH

		gc.Save()
		gc.Translate(d.LabelW, 0)
		drawDropShadow(gc, x1, y1, x2, y2)
		drawBlock(d, gc, x1, y1, x2, y2, task.BorderColor, task.FillColor, task.Label, false)
		recur, _ := strconv.Atoi(task.Recur)
		if recur > 0 {
			increment := float64(recur) * d.DayW
			recurW := x2 - x1
			for recurX := x1 + increment; recurX < w; recurX += increment {
				drawDropShadow(gc, recurX, y1, recurX+recurW, y2)
				drawBlock(d, gc, recurX, y1, recurX+recurW, y2, task.BorderColor, task.FillColor, "", false)
			}
		}
		gc.Restore()

		//write out label
		label := task.Label
		gc.Save()
		gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
		_, ty1, _, ty2 := gc.GetStringBounds(label)
		th := ty2 - ty1
		adjustedTextY := y2 - th*0.5

		gc.FillStringAt(label, 0, adjustedTextY)
		gc.Restore()

		//draw milestones
		gc.Save()
		gc.Translate(d.LabelW, 0)
		for _, milestone := range task.Milestones {
			milestoneTime := parseDateStamp(milestone)

			for i := 0; i <= d.Days; i++ {
				itTime := time.Date(d.First.Year(), d.First.Month(), d.First.Day()+i, 0, 0, 0, 0, time.UTC)
				if milestoneTime == itTime {
					drawMilestone(d, gc, float64(i)*d.DayW, y1+rowH/2)
				}
			}
		}
		gc.Restore()

		//draw datestamps
		gc.Save()
		gc.Translate(d.LabelW, 0)
		for _, dateStamp := range task.DateStamps {
			dateStampTime := parseDateStamp(dateStamp)
			i := dayIndex(dateStampTime, d.First, d.Last)
			if i == -1 {
				continue
			}
			drawDateStamp(d, gc, float64(i)*d.DayW, y1+rowH/2, dateStamp)
		}
		gc.Restore()

		//draw arrows (end)
		gc.Save()
		gc.Translate(d.LabelW, 0)
		for _, endTo := range task.EndTo {
			arrowStartTime := task.EndTime

			//reject index if out of bounds
			if len(d.Tasks) < index+endTo {
				continue
			}

			//reject index if dest start < source start
			if dayIndex(task.StartTime, d.First, d.Last) > dayIndex(d.Tasks[index+int(endTo)].StartTime, d.First, d.Last) {
				continue
			}

			arrowEndTime := d.Tasks[index+int(endTo)].StartTime

			if dayIndex(arrowEndTime, d.First, d.Last) >= dayIndex(arrowStartTime, d.First, d.Last) {
				x1 := float64(dayIndex(arrowStartTime, d.First, d.Last)) * d.DayW
				x2 := float64(dayIndex(arrowEndTime, d.First, d.Last)) * d.DayW

				y1 := y + d.RowH/2
				y2 := y + float64(endTo)*2*d.RowH

				gc.BeginPath()
				gc.MoveTo(x1, y1)
				gc.LineTo(x2, y1)
				gc.LineTo(x2, y2)
				drawArrowHeadS(d, gc, x2, y2)
				gc.FillStroke()
			} else {
				//special case: arrow moves in direct line down
				x1 := float64(dayIndex(arrowEndTime, d.First, d.Last)) * d.DayW
				x2 := x1
				y1 := y + d.RowH
				y2 := y + float64(endTo)*2*d.RowH

				gc.BeginPath()
				gc.MoveTo(x1, y1)
				gc.LineTo(x2, y2)
				drawArrowHeadS(d, gc, x2, y2)
				gc.FillStroke()
			}
		}
		gc.Restore()

		//draw arrows (start)
		gc.Save()
		gc.Translate(d.LabelW, 0)
		for _, startTo := range task.StartTo {
			arrowStartTime := task.StartTime

			//reject index if out of bounds
			if len(d.Tasks) < index+startTo {
				continue
			}

			//reject index if dest start < source start
			if dayIndex(task.StartTime, d.First, d.Last) > dayIndex(d.Tasks[index+int(startTo)].StartTime, d.First, d.Last) {
				continue
			}

			arrowEndTime := d.Tasks[index+int(startTo)].StartTime

			x1 := float64(dayIndex(arrowStartTime, d.First, d.Last)) * d.DayW
			x2 := float64(dayIndex(arrowEndTime, d.First, d.Last)) * d.DayW

			y1 := y + d.RowH
			y2 := y + float64(startTo)*2*d.RowH

			if x1 == x2 {
				gc.BeginPath()
				gc.MoveTo(x1, y1)
				gc.LineTo(x2, y2)
				drawArrowHeadS(d, gc, x2, y2)
				gc.FillStroke()
			} else if x1 < x2 {
				y2 += d.RowH / 2

				gc.BeginPath()
				gc.MoveTo(x1, y1)
				gc.LineTo(x1, y2)
				gc.LineTo(x2, y2)
				drawArrowHeadE(d, gc, x2, y2)
				gc.FillStroke()
			}
		}
		gc.Restore()
		y += d.RowH * 2
	}

	//crop dest
	y += rowH
	rect := image.Rect(0.0, 0.0, int(w), int(y+rowH))
	cropped := dest.SubImage(rect)
	// Save to file
	draw2dimg.SaveToPngFile(path, cropped)
}