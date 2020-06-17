package raw

import (
	"fmt"
	"io"
	"time"

	"github.com/Danile71/joy4/av"
	"github.com/Danile71/joy4/codec/h264parser"
	"github.com/Danile71/joy4/utils/bits/pio"
)

type Demuxer struct {
	r io.ReadSeeker
}

func NewDemuxer(r io.ReadSeeker) *Demuxer {
	return &Demuxer{
		r: r,
	}
}

func (d *Demuxer) readBytes() (data []byte, err error) {
	size4 := make([]byte, 4)

	// TODO: check read count
	if _, err = d.r.Read(size4); err != nil {
		return
	}
	size := pio.U32BE(size4)
	data = make([]byte, size)
	if _, err = d.r.Read(data); err != nil {
		return
	}

	return
}

func (d *Demuxer) readI64() (data int64, err error) {
	d8 := make([]byte, 8)
	// TODO: check read count
	if _, err = d.r.Read(d8); err != nil {
		return
	}
	data = pio.I64BE(d8)
	return
}

func (d *Demuxer) readI32() (data int32, err error) {
	d4 := make([]byte, 4)
	// TODO: check read count
	if _, err = d.r.Read(d4); err != nil {
		return
	}
	data = pio.I32BE(d4)
	return
}

func (d *Demuxer) readU32() (data uint32, err error) {
	d4 := make([]byte, 4)
	// TODO: check read count
	if _, err = d.r.Read(d4); err != nil {
		return
	}
	data = pio.U32BE(d4)
	return
}

func (d *Demuxer) readI8() (data int8, err error) {
	d1 := make([]byte, 1)
	// TODO: check read count
	if _, err = d.r.Read(d1); err != nil {
		return
	}
	data = int8(d1[0])
	return
}

func (d *Demuxer) readU8() (data uint8, err error) {
	d1 := make([]byte, 1)
	// TODO: check read count
	if _, err = d.r.Read(d1); err != nil {
		return
	}
	data = d1[0]
	return
}

func (d *Demuxer) readBool() (data bool, err error) {
	d1 := make([]byte, 1)
	// TODO: check read count
	if _, err = d.r.Read(d1); err != nil {
		return
	}
	data = d1[0] == 0xFF
	return
}

func (d *Demuxer) Streams() (streams []av.CodecData, err error) {
	var codecType uint32
	if codecType, err = d.readU32(); err != nil {
		return
	}

	switch av.CodecType(codecType) {
	case av.H264:
		var record []byte
		if record, err = d.readBytes(); err != nil {
			return
		}
		var codec h264parser.CodecData
		if codec, err = h264parser.NewCodecDataFromAVCDecoderConfRecord(record); err != nil {
			return
		}
		streams = append(streams, codec)
	case av.AAC:
		err = fmt.Errorf("mp4: codec type=%v is not implement", codecType)
		return

	default:
		err = fmt.Errorf("mp4: codec type=%v is not supported", codecType)
		return

	}

	return
}

func (d *Demuxer) ReadPacket() (pkt av.Packet, err error) {
	if pkt.IsKeyFrame, err = d.readBool(); err != nil {
		return
	}

	if pkt.Idx, err = d.readI8(); err != nil {
		return
	}

	var i64 int64
	if i64, err = d.readI64(); err != nil {
		return
	} else {
		pkt.CompositionTime = time.Duration(i64)
	}

	if i64, err = d.readI64(); err != nil {
		return
	} else {
		pkt.Time = time.Duration(i64)
	}

	if pkt.Data, err = d.readBytes(); err != nil {
		return
	}

	return
}
