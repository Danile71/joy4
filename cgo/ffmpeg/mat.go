// +build gocv

package ffmpeg

// #cgo !windows pkg-config: opencv4
//#cgo LDFLAGS: -lavformat -lavutil -lavcodec -lavresample -lswscale
// #cgo CXXFLAGS:   --std=c++1z
// #cgo windows  CPPFLAGS:   -IC:/opencv/build/install/include
// #cgo windows  LDFLAGS:    -LC:/opencv/build/install/x64/mingw/lib -lopencv_core420 -lopencv_face420 -lopencv_videoio420 -lopencv_imgproc420 -lopencv_highgui420 -lopencv_imgcodecs420 -lopencv_objdetect420 -lopencv_features2d420 -lopencv_video420 -lopencv_dnn420 -lopencv_xfeatures2d420 -lopencv_plot420 -lopencv_tracking420 -lopencv_img_hash420 -lopencv_calib3d420
// #include <stdlib.h>
// #include <stdint.h>
// #include "gocv.h"
// #include "ffmpeg.h"
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"gocv.io/x/gocv"
)

func (self *VideoDecoder) DecodeMat(pkt []byte) (img gocv.Mat, err error) {
	ff := &self.ff.ff
	cgotimg := C.int(0)
	frame := C.av_frame_alloc()
	defer C.av_frame_free(&frame)

	cerr := C.decode(ff.codecCtx, frame, (*C.uchar)(unsafe.Pointer(&pkt[0])), C.int(len(pkt)), &cgotimg)

	if cerr < C.int(0) {
		err = fmt.Errorf("ffmpeg: decode failed: %d", cerr)
		return
	}

	if cgotimg != C.int(0) {
		w := int(frame.width)
		h := int(frame.height)
		sz := C.int(0)
		nframe := C.av_frame_alloc()
		defer C.av_frame_free(&nframe)

		cdata := C.avcodec_encode_to_mat(ff.codecCtx, frame, nframe, &sz)
		if cerr != C.int(0) {
			err = fmt.Errorf("ffmpeg: avcodec_encode_jpeg failed: %d", cerr)
			return
		}
		data := make([]byte, int(sz))
		copy(data, *(*[]byte)(unsafe.Pointer(&cdata)))
		img, err = gocv.NewMatFromBytes(h, w, gocv.MatTypeCV8UC3, data)
		return
	}

	return img, errors.New("no image")
}
