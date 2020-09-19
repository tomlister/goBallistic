package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/cnkei/gospline"
	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/tomlister/goBallistic/ballistics"
)

type Game struct {
	Solver     ballistics.Solver
	Frames     ballistics.Frames
	Spline     gospline.Spline
	BakedCurve *ebiten.Image
	vxmax      float64
	ftmax      float64
	finaly     float64
}

func (g *Game) Update(screen *ebiten.Image) error {
	// Write your game's logical update.
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.BakedCurve, op)

	xline, _ := ebiten.NewImage(640, 1, ebiten.FilterDefault)
	xline.Fill(color.White)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, g.finaly)
	screen.DrawImage(xline, op)
	yline, _ := ebiten.NewImage(1, int(g.finaly), ebiten.FilterDefault)
	yline.Fill(color.White)
	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(yline, op)
	ytext, _ := ebiten.NewImage(100, 20, ebiten.FilterDefault)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(-1.5708)
	op.GeoM.Translate(10, (g.finaly/2)+50)
	text.Draw(ytext, "X Velocity (m/s)", bitmapfont.Gothic12r, 5, 8, color.White)
	screen.DrawImage(ytext, op)
	text.Draw(screen, "Flight Time (s)", bitmapfont.Gothic12r, 320-38, int(g.finaly)+20, color.White)
	text.Draw(screen, fmt.Sprintf("Bullet Specifications:\nName - %s\nDiameter - %fm\nMass - %fkg\nG1 Coefficient: %f\nExternal Factors:\nMuzzle Velocity X - %fm/s\nAir Density - %fkg/m^3", g.Solver.Name, g.Solver.Bullet.Diameter, g.Solver.Bullet.Mass, g.Solver.Bullet.G1, g.Solver.MuzzleVelocityX, g.Solver.AirDensity), bitmapfont.Gothic12r, 20, 480-140, color.White)
	cx, cy := ebiten.CursorPosition()
	if float64(cy) <= g.finaly {
		point, _ := ebiten.NewImage(4, 4, ebiten.FilterDefault)
		point.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cx)-2, (480-(g.Spline.At((float64(cx)/640)*g.ftmax)/g.vxmax)*480)-2)
		screen.DrawImage(point, op)
		xpos := cx
		if 640-cx <= 180 {
			xpos = cx - 180
		}
		text.Draw(screen, fmt.Sprintf("vx:%fm/s, ft:%fs", g.Spline.At((float64(cx)/640)*g.ftmax), (float64(cx)/640)*g.ftmax), bitmapfont.Gothic12r, xpos, int(480-(g.Spline.At((float64(cx)/640)*g.ftmax)/g.vxmax)*480)-10, color.RGBA{0xff, 0x00, 0x00, 0xff})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	// Specify bullet specs
	/*
		Based on ADI Australian Outback 155.5gr Berger Target
	*/
	bullet := ballistics.Bullet{
		Name:     "Australian Outback .308 155.5gr Berger Target",
		Diameter: 0.0078232,   // m
		Mass:     0.010076231, // kg
		G1:       0.473,       // G1 Ballistic Coefficient
	}
	solver := ballistics.Solver{
		Bullet:          bullet,
		AirDensity:      1.2029,  // kg/m^3
		MuzzleVelocityX: 871.728, // m/s
		MuzzleVelocityY: 10,      // m/s
	}
	frames := solver.Solve()
	vxmax := 0.0
	ftmax := 0.0
	vx := []float64{}
	ft := []float64{}
	for _, frame := range frames {
		vx = append(vx, frame.VelocityX)
		ft = append(ft, frame.FlightTime)
		if vxmax < frame.VelocityX {
			vxmax = frame.VelocityX
		}
		if ftmax < frame.FlightTime {
			ftmax = frame.FlightTime
		}
	}
	curve := gospline.NewCubicSpline(ft, vx)
	img, _ := ebiten.NewImage(640, 480, ebiten.FilterDefault)
	finaly := 0.0
	for x := 0; x < 640; x++ {
		point, _ := ebiten.NewImage(1, 1, ebiten.FilterDefault)
		point.Fill(color.White)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), 480-(curve.At((float64(x)/640)*ftmax)/vxmax)*480)
		img.DrawImage(point, op)
		finaly = 480 - (curve.At((float64(x)/640)*ftmax)/vxmax)*480
	}
	game := &Game{
		Solver:     solver,
		Frames:     frames,
		Spline:     curve,
		vxmax:      vxmax,
		ftmax:      ftmax,
		BakedCurve: img,
		finaly:     finaly,
	}
	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("goBallistic")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
