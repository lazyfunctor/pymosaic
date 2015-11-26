/* Created by "go tool cgo" - DO NOT EDIT. */

/* package _/home/abhinav/dev/pymosaic */

/* Start of preamble from import "C" comments.  */


#line 3 "/home/abhinav/dev/pymosaic/pymosaic.go"

 #define Py_LIMITED_API
 #include <Python.h>
 int PyArg_ParseTuple_SISS(PyObject *, char **, int *, char **, char **);
 int PyArg_ParseTuple_S(PyObject *, char **);
 int PyArg_ParseTuple_SSS(PyObject *, char **, char **, char **);



/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef __complex float GoComplex64;
typedef __complex double GoComplex128;

// static assertion to make sure the file is being used on architecture
// at least with matching size of GoInt.
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

typedef struct { char *p; GoInt n; } GoString;
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif


extern PyObject* download(PyObject* p0, PyObject* p1);

extern PyObject* analyze_images(PyObject* p0, PyObject* p1);

extern PyObject* generate_mosaic(PyObject* p0, PyObject* p1);

#ifdef __cplusplus
}
#endif
