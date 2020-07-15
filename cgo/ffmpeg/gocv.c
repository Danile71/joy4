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
    struct SwsContext *img_convert_ctx = sws_getCachedContext( NULL, pCodecCtx->width, pCodecCtx->height, pCodecCtx->pix_fmt, pFrame->width, pFrame->height, AV_PIX_FMT_BGR24, SWS_BICUBIC, NULL, NULL, NULL );
    nFrame->format = AV_PIX_FMT_BGR24;
    nFrame->width = pFrame->width;
    nFrame->height = pFrame->height;
    *size = av_image_get_buffer_size( AV_PIX_FMT_BGR24, pFrame->width, pFrame->height, 1);
    uint8_t *tmp_picture_buf = (uint8_t *)malloc(*size);    
    av_image_fill_arrays(nFrame->data, nFrame->linesize, tmp_picture_buf, AV_PIX_FMT_BGR24, pFrame->width, pFrame->height, 1);
    sws_scale(img_convert_ctx, (const uint8_t* const*)pFrame->data, pFrame->linesize, 0, nFrame->height, nFrame->data, nFrame->linesize);
    return tmp_picture_buf;
}
