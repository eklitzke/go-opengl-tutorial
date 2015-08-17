package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

// PrintLog prints the error log for an object
func PrintLog(object uint32) {
	var logLength int32
	if gl.IsShader(object) {
		gl.GetShaderiv(object, gl.INFO_LOG_LENGTH, &logLength)
	} else if gl.IsProgram(object) {
		gl.GetProgramiv(object, gl.INFO_LOG_LENGTH, &logLength)
	} else {
		log.Fatal("PrintLog: not a shader or program")
	}

	infoLog := strings.Repeat("\x00", int(logLength+1))
	if gl.IsShader(object) {
		gl.GetShaderInfoLog(object, logLength, nil, gl.Str(infoLog))
	} else if gl.IsProgram(object) {
		gl.GetProgramInfoLog(object, logLength, nil, gl.Str(infoLog))
	}
	log.Fatal(infoLog)
}

// CreateShader creates a shader object from a source file.
func CreateShader(filename string, xtype uint32) uint32 {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	dataString := "#version 120\n" + string(data) + "\x00"
	dataSrc := gl.Str(dataString)

	res := gl.CreateShader(xtype)
	gl.ShaderSource(res, 1, &dataSrc, nil)
	gl.CompileShader(res)

	var compileOk int32
	gl.GetShaderiv(res, gl.COMPILE_STATUS, &compileOk)
	if compileOk == gl.FALSE {
		PrintLog(res)
		gl.DeleteShader(res)
	}
	return res
}
