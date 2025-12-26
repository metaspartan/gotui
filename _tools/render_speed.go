package main
import (
	"fmt"
	"log"
	"math/rand"
	"time"
	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)
func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()
	var grid []ui.Drawable
	rows := 20
	cols := 20
	w, h := ui.TerminalDimensions()
	blockW := w / cols
	blockH := h / rows
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			sl := widgets.NewSparkline()
			sl.Data = data
			sl.LineColor = ui.ColorGreen
			group := widgets.NewSparklineGroup(sl)
			group.SetRect(c*blockW, r*blockH, (c+1)*blockW, (r+1)*blockH)
			group.Border = false 
			group.Border = true
			grid = append(grid, group)
		}
	}
	frameCount := 100
	start := time.Now()
	for i := 0; i < frameCount; i++ {
		ui.Render(grid...)
		if i%10 == 0 {
			for _, item := range grid {
				if g, ok := item.(*widgets.SparklineGroup); ok {
					g.Sparklines[0].Data[rand.Intn(len(data))] = float64(rand.Intn(100))
				}
			}
		}
	}
	elapsed := time.Since(start)
	avgPerFrame := elapsed.Seconds() * 1000 / float64(frameCount)
	fps := float64(frameCount) / elapsed.Seconds()
	ui.Close()
	fmt.Printf("Render Speed Test Complete:\n")
	fmt.Printf("Widgets: %d\n", len(grid))
	fmt.Printf("Frames: %d\n", frameCount)
	fmt.Printf("Total Time: %v\n", elapsed)
	fmt.Printf("Avg Time per Frame: %.2f ms\n", avgPerFrame)
	fmt.Printf("FPS Potential: %.2f\n", fps)
}
