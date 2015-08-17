package main

import (
	"log"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const floatSize = 4

var curFade float32
var vboTriangle, vboTriangleColors uint32
var program uint32
var attributeCoord2d, attributeVColor int32
var uniformFade int32

var triangleAttributes = []float32{
	0.0, 0.8, 1.0, 1.0, 0.0,
	-0.8, -0.8, 0.0, 0.0, 1.0,
	0.8, -0.8, 1.0, 0.0, 0.0,
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func initResources() uint32 {
	vs := CreateShader("triangle.v.glsl", gl.VERTEX_SHADER)
	fs := CreateShader("triangle.f.glsl", gl.FRAGMENT_SHADER)

	var linkOk int32
	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkOk)
	if linkOk == 0 {
		log.Fatal("gl.LinkProgram")
	}

	gl.GenBuffers(1, &vboTriangle)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboTriangle)
	gl.BufferData(gl.ARRAY_BUFFER, floatSize*len(triangleAttributes), gl.Ptr(triangleAttributes), gl.STATIC_DRAW)

	attributeCoord2d = gl.GetAttribLocation(program, gl.Str("coord2d\x00"))
	if attributeCoord2d == -1 {
		log.Fatal("failed to bind attribute")
	}

	attributeVColor = gl.GetAttribLocation(program, gl.Str("v_color\x00"))
	if attributeVColor == -1 {
		log.Fatal("could not bind attribute v_color")
	}

	uniformFade = gl.GetUniformLocation(program, gl.Str("fade\x00"))
	if uniformFade == -1 {
		log.Fatal("could not bind uniform fade")
	}

	return program
}

func onDisplay(program uint32) {
	coords := uint32(attributeCoord2d)
	vcolor := uint32(attributeVColor)

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)
	gl.Uniform1f(uniformFade, curFade)

	gl.EnableVertexAttribArray(coords)
	gl.EnableVertexAttribArray(vcolor)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboTriangle)
	gl.VertexAttribPointer(coords, 2, gl.FLOAT, false, 5*floatSize, nil)
	gl.VertexAttribPointer(vcolor, 3, gl.FLOAT, false, 5*floatSize, gl.PtrOffset(2*floatSize))

	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.DisableVertexAttribArray(vcolor)
	gl.DisableVertexAttribArray(coords)
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(640, 480, "tut03", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	t0 := time.Now()
	program := initResources()
	for {
		curFade = float32(math.Sin(time.Now().Sub(t0).Seconds()*2*math.Pi/5)/2 + 0.5)
		onDisplay(program)
		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteProgram(program)
	gl.DeleteBuffers(1, &vboTriangle)
}
