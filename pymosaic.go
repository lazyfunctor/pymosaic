package main

// #cgo pkg-config: python3
// #define Py_LIMITED_API
// #include <Python.h>
// int PyArg_ParseTuple_SISS(PyObject *, char **, int *, char **, char **);
// int PyArg_ParseTuple_S(PyObject *, char **);
// int PyArg_ParseTuple_SSS(PyObject *, char **, char **, char **);
import "C"

import "github.com/lazyfunctor/pymosaic/client"
import "github.com/lazyfunctor/pymosaic/mosaic"
import "fmt"



//export download
func download(self, args *C.PyObject) *C.PyObject {
    // fmt.Println(args)
    var tag, dir, apiKey *C.char
    var count C.int
    if C.PyArg_ParseTuple_SISS(args, &tag, &count, &dir, &apiKey) == 0 {
        return nil
    }
    client.Download(C.GoString(tag), int(count), C.GoString(dir), C.GoString(apiKey))
    return C.PyLong_FromLongLong(1)
}

//export analyze_images
func analyze_images(self, args *C.PyObject) *C.PyObject {
    var inputDir *C.char
    if C.PyArg_ParseTuple_S(args, &inputDir) == 0 {
        return nil
    }
    analysis := mosaic.AnalyzeAll(C.GoString(inputDir))
    return C.PyUnicode_FromString(C.CString(analysis))
}

//export generate_mosaic
func generate_mosaic(self, args *C.PyObject) *C.PyObject {
    var inputImg, outImg, analysis *C.char
    if C.PyArg_ParseTuple_SSS(args, &inputImg, &outImg, &analysis) == 0 {
        return nil
    }
    err := mosaic.AnalyzeTarget(C.GoString(inputImg), C.GoString(analysis), C.GoString(outImg))
    fmt.Println(err)
    return C.PyLong_FromLongLong(1)
}

func main() {}  

