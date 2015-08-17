package main

import (
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func fixupSrc(src string) string {
	return strings.TrimSpace(src) + string([]byte{0})
}

var vShaderSrc, fShaderSrc string

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	vShaderSrc = fixupSrc(
		`#version 120
attribute vec2 coord2d;
void main(void) {
  gl_Position = vec4(coord2d, 0.0, 1.0);
}`)

	fShaderSrc = fixupSrc(`
#version 120
void main(void) {
  gl_FragColor[0] = gl_FragCoord.x/640.0;
  gl_FragColor[1] = gl_FragCoord.y/480.0;
  gl_FragColor[2] = 0.5;
}`)
}

func initResources() (uint32, int32) {
	vSrc0 := fixupSrc(vShaderSrc)
	vSrc := gl.Str(vSrc0)
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vs, 1, &vSrc, nil)
	gl.CompileShader(vs)
	var vok int32
	gl.GetShaderiv(vs, gl.COMPILE_STATUS, &vok)
	if vok == 0 {
		log.Fatal("error in vertex shader")

	}

	fSrc := gl.Str(fShaderSrc)
	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fs, 1, &fSrc, nil)
	gl.CompileShader(fs)
	var fok int32
	gl.GetShaderiv(vs, gl.COMPILE_STATUS, &fok)
	if fok == 0 {
		log.Fatal("error in fragment shader")

	}

	var linkOk int32
	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkOk)
	if linkOk == 0 {
		log.Fatal("gl.LinkProgram")
	}

	attribName := fixupSrc("coord2d")
	attribCoord2d := gl.GetAttribLocation(program, gl.Str(attribName))
	if attribCoord2d == -1 {
		log.Fatal("failed to bind attribute")
	}
	return program, attribCoord2d
}

func onDisplay(program uint32, coords uint32) {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)
	gl.EnableVertexAttribArray(coords)
	triangleVertices := []float32{
		0.0, 0.8,
		-0.8, -0.8,
		0.8, -0.8}
	gl.VertexAttribPointer(coords, 2, gl.FLOAT, false, 0, gl.Ptr(triangleVertices))
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
	window, err := glfw.CreateWindow(640, 480, "tut01", nil, nil)
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
}
