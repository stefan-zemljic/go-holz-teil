package main

import (
	"errors"
	"fmt"
	"github.com/deeean/go-vector/vector3"
	"math"
)

type Done struct {
	Work
	Config
}

type Config struct {
	Direction int
	Rotation  int
	Part      int
}

type Work struct {
	v     vector3.Vector3
	index int
}

func main() {
	var matrix [5][5][5]int
	var workList []Work
	workIndex := 0
	for z := 0.0; z < 5.0; z++ {
		for y := 0.0; y < 5.0; y++ {
			for x := 0.0; x < 5.0; x++ {
				workList = append(workList, Work{
					v:     vector3.Vector3{X: x, Y: y, Z: z},
					index: workIndex,
				})
				workIndex++
			}
		}
	}
	var doneList []Done
	found := 0
	//lastLog := time.Now()
	for workIndex = 0; workIndex < 125; workIndex++ {
		work := workList[workIndex]
		var err error
		var done Done
		if matrix[int(work.v.X)][int(work.v.Y)][int(work.v.Z)] != 0 {
			if workIndex == 124 {
				found++
				fmt.Println(found)
				printMatrix(&matrix)
				done = doneList[len(doneList)-1]
				doneList = doneList[:len(doneList)-1]
				workIndex = done.index
				err = errors.New("solved")
			} else {
				continue
			}
		}
		if err == nil {
			config, newErr := insert(&matrix, work.v, len(doneList)+1, Config{})
			done = Done{Work: work, Config: config}
			err = newErr
		}
		if err == nil && workIndex == 124 {
			found++
			fmt.Println(found)
			printMatrix(&matrix)
			err = errors.New("solved")
		}
		for err != nil {
			if err.Error() == "dead end" {
				done = doneList[len(doneList)-1]
				doneList = doneList[:len(doneList)-1]
			}
			Remove(&matrix, done)
			if incremented, ok := Increment(done.Config); !ok {
				continue
			} else {
				done.Config = incremented
			}
			done.Config, err = insert(&matrix, done.v, len(doneList)+1, done.Config)
			if err == nil && workIndex == 124 {
				found++
				fmt.Println(found)
				printMatrix(&matrix)
				err = errors.New("solved")
			}
		}
		workIndex = done.index
		doneList = append(doneList, done)
	}
}

func insert(matrix *[5][5][5]int, v vector3.Vector3, number int, start Config) (Config, error) {
	for dir := start.Direction; dir < 6; dir++ {
		rot := 0
		if dir == start.Direction {
			rot = start.Rotation
		}
		for ; rot < 4; rot++ {
			part := 0
			if dir == start.Direction && rot == start.Rotation {
				part = start.Part
			}
			for ; part < 5; part++ {
				config := Config{dir, rot, part}
				if tryInsert(matrix, v, config, number) {
					var err error
					if number != 0 {
						err = checkSolvable(matrix)
					}
					return config, err
				}
			}
		}
	}
	return Config{}, errors.New("dead end")
}

func checkSolvable(matrix *[5][5][5]int) error {
	for z := 0; z < 5; z++ {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				if matrix[x][y][z] != 0 {
					continue
				}
				_, err := insert(matrix, vector3.Vector3{X: float64(x), Y: float64(y), Z: float64(z)}, 0, Config{})
				if err != nil {
					return errors.New("unsolvable")
				}
			}
		}
	}
	return nil
}

func tryInsert(matrix *[5][5][5]int, coord vector3.Vector3, config Config, number int) bool {
	part := computePart(config)
	for _, v := range part {
		x := coord.X + v.X
		y := coord.Y + v.Y
		z := coord.Z + v.Z
		if x < 0 || x >= 5 || y < 0 || y >= 5 || z < 0 || z >= 5 {
			return false
		}
		if matrix[int(x)][int(y)][int(z)] != 0 {
			return false
		}
	}
	if number == 0 {
		return true
	}
	for _, v := range part {
		x := coord.X + v.X
		y := coord.Y + v.Y
		z := coord.Z + v.Z
		matrix[int(x)][int(y)][int(z)] = number
	}
	return true
}

var parts [][][][]vector3.Vector3

func computePart(config Config) []vector3.Vector3 {
	if parts == nil {
		parts = make([][][][]vector3.Vector3, 6)
		for i := 0; i < 6; i++ {
			parts[i] = make([][][]vector3.Vector3, 4)
			for j := 0; j < 4; j++ {
				parts[i][j] = make([][]vector3.Vector3, 5)
				for k := 0; k < 5; k++ {
					parts[i][j][k] = doComputePart(Config{i, j, k})
				}
			}
		}
	}
	return parts[config.Direction][config.Rotation][config.Part]
}

func doComputePart(config Config) []vector3.Vector3 {
	part := []vector3.Vector3{
		{},
		{Z: 1},
		{Z: 2},
		{Z: 3},
		{Y: 1, Z: 2},
	}
	for i, v := range part {
		switch config.Direction {
		case 0:
			v = RotateVector(v, vector3.Vector3{Y: 1}, 90)
		case 1:
			v = RotateVector(v, vector3.Vector3{X: 1}, 90)
		case 2:
		case 3:
			v = RotateVector(v, vector3.Vector3{Y: 1}, 270)
		case 4:
			v = RotateVector(v, vector3.Vector3{X: 1}, 270)
		case 5:
			v = RotateVector(v, vector3.Vector3{Y: 1}, 180)
		}
		part[i] = v
	}
	part[4] = RotateVector(part[4], part[1], 90*float64(config.Rotation))
	newZeroCoord := part[config.Part]
	for j, v := range part {
		part[j] = *v.Sub(&newZeroCoord)
	}
	return part
}

func RotateVector(v vector3.Vector3, axis vector3.Vector3, angleDeg float64) vector3.Vector3 {
	// Normalize the axis
	mag := math.Sqrt(axis.X*axis.X + axis.Y*axis.Y + axis.Z*axis.Z)
	if mag == 0 {
		return v
	}
	ux, uy, uz := axis.X/mag, axis.Y/mag, axis.Z/mag

	// Convert angle to radians
	theta := angleDeg * math.Pi / 180.0
	cosT := math.Cos(theta)
	sinT := math.Sin(theta)
	oneMinusCosT := 1 - cosT

	// Rotation matrix elements
	r11 := cosT + ux*ux*oneMinusCosT
	r12 := ux*uy*oneMinusCosT - uz*sinT
	r13 := ux*uz*oneMinusCosT + uy*sinT
	r21 := uy*ux*oneMinusCosT + uz*sinT
	r22 := cosT + uy*uy*oneMinusCosT
	r23 := uy*uz*oneMinusCosT - ux*sinT
	r31 := uz*ux*oneMinusCosT - uy*sinT
	r32 := uz*uy*oneMinusCosT + ux*sinT
	r33 := cosT + uz*uz*oneMinusCosT

	// Apply rotation: newV = R * v
	newX := r11*v.X + r12*v.Y + r13*v.Z
	newY := r21*v.X + r22*v.Y + r23*v.Z
	newZ := r31*v.X + r32*v.Y + r33*v.Z

	return vector3.Vector3{X: math.Round(newX), Y: math.Round(newY), Z: math.Round(newZ)}
}

func Remove(matrix *[5][5][5]int, done Done) {
	part := computePart(done.Config)
	for _, v := range part {
		x := done.v.X + v.X
		y := done.v.Y + v.Y
		z := done.v.Z + v.Z
		matrix[int(x)][int(y)][int(z)] = 0
	}
	return
}

func Increment(config Config) (Config, bool) {
	config.Part++
	if config.Part == 5 {
		config.Part = 0
		config.Rotation++
		if config.Rotation == 4 {
			config.Rotation = 4
			config.Direction++
			return config, config.Direction != 6
		}
	}
	return config, true
}

func DecreaseVector(v vector3.Vector3) (vector3.Vector3, bool) {
	v.X--
	if v.X < 0 {
		v.X = 4
		v.Y--
		if v.Y < 0 {
			v.Y = 4
			v.Z--
			return v, v.Z >= 0
		}
	}
	return v, true
}

func check(matrix *[5][5][5]int, coord vector3.Vector3) {
	fmt.Printf("Checking %v\n", coord)
	for {
		if matrix[int(coord.X)][int(coord.Y)][int(coord.Z)] == 0 {
			printMatrix(matrix)
			panic(fmt.Errorf("empty cell at %v", coord))
		}
		if coord.X == 0 && coord.Y == 0 && coord.Z == 0 {
			break
		}
		coord, _ = DecreaseVector(coord)
	}
}

func printMatrix(matrix *[5][5][5]int) {
	for y := 0; y < 5; y++ {
		for z := 0; z < 5; z++ {
			for x := 0; x < 5; x++ {
				fmt.Printf("%2d ", matrix[x][y][z])
			}
			if z != 4 {
				fmt.Print(" | ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
