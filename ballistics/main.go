package ballistics

import (
	"fmt"
	"math"

	"github.com/cnkei/gospline"
)

/*
	Based on http://www.x-ballistics.eu/cms/ballistics/how-to-calculate-the-trajectory/
*/

// Solver stores bullet, velocity and environmental information
type Solver struct {
	Bullet
	AirDensity      float64 // kg/m^3
	MuzzleVelocityX float64 // m/s
	MuzzleVelocityY float64 // m/s
}

// Bullet stores bullet specifications
type Bullet struct {
	Name     string
	Diameter float64 //Meausred in meters
	Mass     float64 //Measured in kilograms
	G1       float64 //G1 Ballistic Coefficient
}

// Frame stores trajectory data
type Frame struct {
	FlightTime float64
	VelocityX  float64
	VelocityY  float64
	DistanceX  float64
	DistanceY  float64
}

// Frames stores frames
type Frames []Frame

// Solve calculates the trajectory of the bullet
func (s Solver) Solve() Frames {
	diameter := s.Bullet.Diameter * 39.3701 // Inches
	mass := s.Bullet.Mass * 2.20462         // Pounds
	formFactor := mass / (math.Pow(diameter, 2) * s.Bullet.G1)
	timeStep := 0.01
	// Release direction
	vx, vy := s.MuzzleVelocityX, s.MuzzleVelocityY // m/s
	// Starting trajectory values
	sx, sy := 0.0, 0.0         // m
	airDensity := s.AirDensity // kg/m^3
	gravity := 9.81            // m/s^2d
	flightTime := 0.0
	crossSectionalArea := math.Pow(s.Bullet.Diameter, 2.0) * math.Pi / 4.0 //m^2
	mach1 := 340.3                                                         //m/s
	/*
		G1 LUT (Look Up Table)
	*/
	machValues := []float64{
		0.00, 0.20, 0.30, 0.40, 0.50, 0.60, 0.70, 0.80, 0.90, 1.00, 1.10, 1.20, 1.30, 1.40, 1.50, 1.60, 1.70, 1.80, 1.90, 2.00, 2.50, 3.00, 4.00, 5.00,
	}
	cwValues := []float64{
		0.26, 0.23, 0.22, 0.21, 0.20, 0.20, 0.21, 0.25, 0.34, 0.48, 0.58, 0.63, 0.65, 0.66, 0.65, 0.64, 0.63, 0.62, 0.60, 0.59, 0.53, 0.51, 0.50, 0.49,
	}

	frames := Frames{}

	for sx < 1000.0 {
		fmt.Println("------------------------------------------------------------")
		fmt.Printf("DistanceX: %f, DistanceY: %f, VelocityX: %f, VelocityY: %f\n", sx, sy, vx, vy)
		totalVelocity := math.Sqrt(math.Pow(vx, 2.0) + math.Pow(vy, 2.0))
		mach := totalVelocity / mach1
		// Fetch an interpolated air resistance value using mach number
		spline := gospline.NewCubicSpline(machValues, cwValues)
		airResistance := spline.At(mach)
		// Air drag force (N)
		force := 0.5 * airDensity * crossSectionalArea * math.Pow(totalVelocity, 2) * airResistance * formFactor

		// Horizontal and Vertical force (N)
		forceX := force * vx / totalVelocity
		forceY := force * vy / totalVelocity

		// Horizontal and Vertical negative acceleration (m/s^2)
		accelerationX := -forceX / s.Bullet.Mass
		accelerationY := -forceY/s.Bullet.Mass - gravity

		// Difference in velocity
		deltaVelocityX := accelerationX * timeStep
		deltaVelocityY := accelerationY * timeStep

		// Difference in distance
		deltaDistanceX := vx*timeStep + deltaVelocityX*timeStep/2
		deltaDistanceY := vy*timeStep + deltaVelocityY*timeStep/2

		// Update distance
		sx += deltaDistanceX
		sy += deltaDistanceY

		// Update velocity
		vx += deltaVelocityX
		vy += deltaVelocityY

		// Update flightTime
		flightTime += timeStep

		// Create a new frame and append to frame slice
		newFrame := Frame{
			flightTime, vx, vy, sx, sy,
		}
		frames = append(frames, newFrame)
	}
	return frames
}

func (frames Frames) toCSV() string {
	csvBuilder := "flightTime, velocityX, velocityY, distanceX, distanceY\n"
	for _, frame := range frames {
		csvBuilder += fmt.Sprintf("%f, %f, %f, %f, %f\n", frame.FlightTime, frame.VelocityX, frame.VelocityY, frame.DistanceX, frame.VelocityY)
	}
	return csvBuilder
}
