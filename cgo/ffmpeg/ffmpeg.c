#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libavutil/imgutils.h>
#include <libswresample/swresample.h>
#include <libavutil/opt.h>
#include <string.h>
#include <libswscale/swscale.h>
#include "ffmpeg.h"

int decode(AVCodecContext *avctx, AVFrame *frame, uint8_t *data, int size, int *got_frame)
{
    int ret;
	struct AVPacket pkt = {.data = data, .size = size};

    *got_frame = 0;
    
    ret = avcodec_send_packet(avctx, &pkt);
    
    av_packet_unref(&pkt);
    
    if (ret < 0)
      return ret == AVERROR_EOF ? 0 : ret;
    

    ret = avcodec_receive_frame(avctx, frame);
    if (ret < 0 && ret != AVERROR(EAGAIN) && ret != AVERROR_EOF)
        return ret;
    if (ret >= 0)
        *got_frame = 1;

    return 0;
}

int encode(AVCodecContext *avctx, AVPacket *pkt, int *got_packet, AVFrame *frame)
{
    int ret;

    *got_packet = 0;

    ret = avcodec_send_frame(avctx, frame);
    if (ret < 0)
        return ret;

    ret = avcodec_receive_packet(avctx, pkt);
    if (!ret)
        *got_packet = 1;
    if (ret == AVERROR(EAGAIN))
        return 0;

    return ret;
}

int avcodec_encode_jpeg(AVCodecContext *pCodecCtx, AVFrame *pFrame,AVPacket *packet) {
    AVCodec *jpegCodec = avcodec_find_encoder(AV_CODEC_ID_MJPEG);
    int ret = -1;

    if (!jpegCodec) {
        return ret;
    }

    AVCodecContext *jpegContext = avcodec_alloc_context3(jpegCodec);
    if (!jpegContext) {
        jpegCodec = NULL;
        return ret;
    }

    jpegContext->pix_fmt = AV_PIX_FMT_YUVJ420P;
    jpegContext->height = pFrame->height;
    jpegContext->width = pFrame->width;
    jpegContext->time_base= (AVRational){1,25};

    ret = avcodec_open2(jpegContext, jpegCodec, NULL);
    if (ret < 0) {
        goto error;
    }
    
    int gotFrame;

    ret = encode(jpegContext, packet, &gotFrame, pFrame);
    if (ret < 0) {
        goto error;
    }
    
    error:
        avcodec_close(jpegContext);
        avcodec_free_context(&jpegContext);
        jpegCodec = NULL;
    return ret;
}



int avcodec_encode_jpeg_nv12(AVCodecContext *pCodecCtx, AVFrame *pFrame,AVPacket *packet) {
    struct SwsContext *img_convert_ctx = sws_getCachedContext( NULL, pCodecCtx->width, pCodecCtx->height, pCodecCtx->pix_fmt, pFrame->width, pFrame->height, AV_PIX_FMT_BGR24, SWS_BICUBIC, NULL, NULL, NULL );
    AVFrame *nFrame = av_frame_alloc();
    nFrame->format = AV_PIX_FMT_YUVJ420P;
    nFrame->width = pFrame->width;
    nFrame->height = pFrame->height;
    
    int size = av_image_get_buffer_size( AV_PIX_FMT_YUVJ420P, pFrame->width, pFrame->height, 1);
    
    uint8_t *tmp_picture_buf = (uint8_t *)malloc(size);    
    
    av_image_fill_arrays(nFrame->data, nFrame->linesize, tmp_picture_buf, AV_PIX_FMT_YUVJ420P, pFrame->width, pFrame->height, 1);
    
    sws_scale(img_convert_ctx, (const uint8_t* const*)pFrame->data, pFrame->linesize, 0, nFrame->height, nFrame->data, nFrame->linesize);        
    
    return avcodec_encode_jpeg(pCodecCtx,nFrame,packet);
}

