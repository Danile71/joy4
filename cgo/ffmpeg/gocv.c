#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libavutil/imgutils.h>
#include <libswresample/swresample.h>
#include <libavutil/opt.h>
#include <string.h>
#include <libswscale/swscale.h>
#include "ffmpeg.h"

uint8_t *avcodec_encode_to_mat(AVCodecContext *pCodecCtx, AVFrame *pFrame,AVFrame *nFrame,int *size) {
    return convert(pCodecCtx, pFrame, nFrame, size, AV_PIX_FMT_BGR24);
}
