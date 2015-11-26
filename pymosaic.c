#define Py_LIMITED_API
#include <Python.h>

PyObject * download(PyObject *, PyObject *);
PyObject * analyze_images(PyObject *, PyObject *);
PyObject * generate_mosaic(PyObject *, PyObject *);


int PyArg_ParseTuple_SISS(PyObject * args, char ** str, int * count, char ** inp, char ** key) {  
    return PyArg_ParseTuple(args, "siss", str, count, inp, key);
}

int PyArg_ParseTuple_S(PyObject * args, char ** inputDir) {  
    return PyArg_ParseTuple(args, "s", inputDir);
}

int PyArg_ParseTuple_SSS(PyObject * args, char ** inputImg, char ** outImg, char ** analysis) {  
    return PyArg_ParseTuple(args, "sss", inputImg, outImg, analysis);
}


static PyMethodDef PymosaicMethods[] = {  
    {"download", download, METH_VARARGS, "flickr client"},
    {"analyze_images", analyze_images, METH_VARARGS, "analyze images and generate statistics"},
    {"generate_mosaic", generate_mosaic, METH_VARARGS, "generate mosaic from target image"},
    {NULL, NULL, 0, NULL}
};

static struct PyModuleDef PymosaicModule = {  
   PyModuleDef_HEAD_INIT, "pymosaic", NULL, -1, PymosaicMethods
};

PyMODINIT_FUNC  
PyInit_pymosaic(void)  
{
    return PyModule_Create(&PymosaicModule);
}
