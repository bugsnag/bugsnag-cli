#include <jni.h>

// Write C++ code here.
//
// Do not forget to dynamically load the C++ library into your application.
//
// For instance,
//
// In MainActivity.java:
//    static {
//       System.loadLibrary("picoapp");
//    }
//
// Or, in MainActivity.kt:
//    companion object {
//      init {
//         System.loadLibrary("picoapp")
//      }
//    }
extern "C"
JNIEXPORT jint JNICALL
Java_com_example_picoapp_MainActivity_add(JNIEnv *env, jobject thiz, jint x, jint y) {
    return x + y;
}