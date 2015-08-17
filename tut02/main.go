package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var vboTriangle uint32
var triangleVertices = []float32{0.0, 0.8, -0.8, -0.8, 0.8, -0.8}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func initResources() (uint32, int32) {
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
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangleVertices), gl.Ptr(triangleVertices), gl.STATIC_DRAW)

	attribName := "coord2d\x00"
	attribCoord2d := gl.GetAttribLocation(program, gl.Str(attribName))
	if attribCoord2d == -1 {
		log.Fatal("failed to bind attribute")
	}
	return program, attribCoord2d
}

func onDisplay(program uint32, coords uint32) {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.UseProgram(program)
	gl.EnableVertexAttribArray(coords)

	gl.VertexAttribPointer(coords, 2, gl.FLOAT, false, 0, nil)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
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
	window, err := glfw.CreateWindow(640, 480, "tut02", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	program, coords := initResources()
	for {
		onDisplay(program, uint32(coords))
		window.SwapBuffers()
	}

	gl.DeleteProgram(program)
	gl.DeleteBuffers(1, &vboTriangle)
}
