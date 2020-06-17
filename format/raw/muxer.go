package raw

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Danile71/joy4/av"
	"github.com/Danile71/joy4/codec/h264parser"
	"github.com/Danile71/joy4/utils/bits/pio"
)

type Muxer struct {
	w       io.WriteSeeker
	bufw    *bufio.Writer
	wpos    int64
	streams []*Stream
}

func NewMuxer(w io.WriteSeeker) *Muxer {
	return &Muxer{
		w:    w,
		bufw: bufio.NewWriterSize(w, pio.RecommendBufioSize),
	}
}

///* AvRational */
//typedef struct _PACKTED_ RawRational {
//	int32_t num;  ///< Numerator
//	int32_t den;  ///< Denominator
//} RawRational;
//
//typedef struct _PACKTED_ RawHeader {
//	/* codec context */
//	struct {
//		RawRational sample_aspect_ratio;
//		int32_t width;
//		int32_t height;
//		int32_t format;
//		int32_t codec_type; /*enum AVMediaType*/
//		int32_t codec_id;   /*enum AVCodecID*/
//	} ccx;
//
//	/* stream */
//	struct {
//		int64_t start_time;
//		RawRational time_base;
//		int64_t nb_frames;
//		RawRational sample_aspect_ratio;
//		RawRational r_frame_rate;
//		RawRational avg_frame_rate;
//	} st;
//} RawHeader;

func (m *Muxer) newStream(codec av.CodecData) (err error) {
	switch codec.Type() {
	case av.H264:
		// videoCodec := codec.(av.VideoCodecData)
		// videoCodec.Width()
		// videoCodec.Height()

	case av.AAC:
	default:
		err = fmt.Errorf("raw: codec type=%v is not supported", codec.Type())
		return
	}
	stream := &Stream{CodecData: codec, muxer: m}

	m.streams = append(m.streams, stream)
	return
}

func (m *Muxer) writeBytes(data []byte) (err error) {
	size := make([]byte, 4)
	pio.PutU32BE(size, uint32(len(data)))

	// TODO: check write count
	if _, err = m.bufw.Write(size); err != nil {
		return
	}
	if _, err = m.bufw.Write(data); err != nil {
		return
	}

	return nil
}

func (m *Muxer) writeI64(data int64) (err error) {
	data8 := make([]byte, 8)
	pio.PutI64BE(data8, data)
	// TODO: check write count
	if _, err = m.bufw.Write(data8); err != nil {
		return
	}

	return nil
}

func (m *Muxer) writeI32(data int32) (err error) {
	data4 := make([]byte, 4)
	pio.PutI32BE(data4, data)
	// TODO: check write count
	if _, err = m.bufw.Write(data4); err != nil {
		return
	}

	return nil
}

func (m *Muxer) writeU32(data uint32) (err error) {
	data4 := make([]byte, 4)
	pio.PutU32BE(data4, data)
	// TODO: check write count
	if _, err = m.bufw.Write(data4); err != nil {
		return
	}

	return nil
}

func (m *Muxer) writeI8(data int8) (err error) {
	// TODO: check write count
	if _, err = m.bufw.Write([]byte{byte(data)}); err != nil {
		return
	}

	return nil
}

func (m *Muxer) writeU8(data uint8) (err error) {
	// TODO: check write count
	if _, err = m.bufw.Write([]byte{data}); err != nil {
		return
	}

	return nil
}

func (m *Muxer) writeBool(data bool) (err error) {
	if data {
		return m.writeU8(0xFF)
	}

	return m.writeU8(0x00)
}

func (m *Muxer) WriteHeader(codecs []av.CodecData) (err error) {
	// m.streams = []*Stream{}
	m.writeU32(uint32(len(codecs)))
	for _, codec := range codecs {
		if err = m.writeU32(uint32(codec.Type())); err != nil {
			return
		}

		switch codec.Type() {
		case av.H264:
			h264 := codec.(h264parser.CodecData)
			// write Record data
			if err = m.writeBytes(h264.Record); err != nil {
				return
			}

			//{ // write RecordInfo
			//	size := uint32(h264.RecordInfo.Len())
			//	if err = m.writeU32(size); err != nil {
			//		return
			//	}
			//
			//	buf := make([]byte, size)
			//	h264.RecordInfo.Marshal(buf)
			//	if _, err = m.bufw.Write(buf); err != nil {
			//		return
			//	}
			//}

		case av.AAC:
			err = fmt.Errorf("mp4: codec type=%v is not implement", codec.Type())
			return
		default:
			err = fmt.Errorf("mp4: codec type=%v is not supported", codec.Type())
			return
		}
	}

	return nil
}

///* AvPacket */
//typedef struct _PACKTED_ RawPacket {
//	int64_t pts;
//	int64_t dts;
//	int32_t flags;
//	int64_t duration;
//	int64_t pos;
//	int32_t size;
//} RawPacket;
func (m *Muxer) WritePacket(pkt av.Packet) (err error) {

	//IsKeyFrame      bool // video packet is key frame
	//Idx             int8 // stream index in container format
	//CompositionTime time.Duration // packet presentation time minus decode time for H264 B-Frame
	//Time time.Duration // packet decode time
	//Data            []byte // packet data

	m.writeBool(pkt.IsKeyFrame)
	m.writeI8(pkt.Idx)
	m.writeI64(int64(pkt.CompositionTime))
	m.writeI64(int64(pkt.Time))
	m.writeBytes(pkt.Data)
	//stream := m.streams[pkt.Idx]
	//if stream.lastpkt != nil {
	//	if err = stream.writePacket(*stream.lastpkt, pkt.Time-stream.lastpkt.Time); err != nil {
	//		return
	//	}
	//}
	//stream.lastpkt = &pkt
	return
}

func (m *Muxer) WriteTrailer() (err error) {
	if err = m.bufw.Flush(); err != nil {
		return
	}

	return nil
}
