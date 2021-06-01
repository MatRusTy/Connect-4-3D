package main

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

const boardSize float32 = 2
const boardHeight float32 = 0.2

const rodSpace = boardSize / 4
const rodOffset = rodSpace / 2
const rodRadius float64 = float64(boardSize / 60)
const rodHeight float64 = float64(boardSize / 5)

const pigRadius = rodRadius * 1.4
const pigHeight = rodHeight / 4

var a = app.App()
var scene = core.NewNode()

func main() {

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 1.5, 4)
	cam.LookAt(math32.NewVector3(boardSize/2, 0, boardSize/2), math32.NewVector3(0, 1, 0))
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Add board
	geom := geometry.NewBox(boardSize, boardHeight, boardSize)
	mat := material.NewStandard(math32.NewColor("Sienna"))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(boardSize/2, -boardHeight/2, boardSize/2)
	scene.Add(mesh)

	// Add rods
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			rod := geometry.NewCylinder(rodRadius, rodHeight, 20, 5, true, false)
			mat := material.NewStandard(math32.NewColor("Peru"))
			mesh := graphic.NewMesh(rod, mat)
			x, z := getPosition(row, col)
			mesh.SetPosition(x, float32(rodHeight/2), z)
			scene.Add(mesh)
		}
	}

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(0, 3, 0)
	scene.Add(pointLight)

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	addPig(0, 0, 0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}

func addPig(player int, row int, col int) {
	pig := geometry.NewCylinder(pigRadius, pigHeight, 20, 5, true, false)
	mat := material.NewStandard(math32.NewColor("Blue"))
	mesh := graphic.NewMesh(pig, mat)
	x, z := getPosition(row, col)
	mesh.SetPosition(x, float32(pigHeight/2), z)
	scene.Add(mesh)
}

func getPosition(row int, col int) (float32, float32) {
	return float32(row)*rodSpace + rodOffset, float32(col)*rodSpace + rodOffset
}
