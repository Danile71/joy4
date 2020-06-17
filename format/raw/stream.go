package raw

import (
	"fmt"
	"time"

	"github.com/Danile71/joy4/av"
)

type Stream struct {
	av.CodecData
	idx       int
	timeScale int64
	muxer     *Muxer
	lastpkt   *av.Packet
}

func (self *Stream) writePacket(pkt av.Packet, rawdur time.Duration) (err error) {
	if rawdur < 0 {
		err = fmt.Errorf("mp4: stream#%d time=%v < lasttime=%v", pkt.Idx, pkt.Time, self.lastpkt.Time)
		return
	}

	//writeFn := func(data interface{}) error {
	//	return binary.Write(self.muxer.w, binary.BigEndian, data)
	//}

	//	int64_t pts;
	//	int64_t dts;
	//	int32_t flags;
	//	int64_t duration;
	//	int64_t pos;
	//	int32_t size;

	//	pts := self.timeToTs(rawdur)
	//dts := pts
	//writeFn(pts)
	//writeFn(dts)
	//writeFn(int32(0))
	//writeFn(pts)
	//writeFn(pts)
	//writeFn(int32(0))
	//writeFn(int32(len(pkt.Data)))
	if _, err = self.muxer.bufw.Write(pkt.Data); err != nil {
		return
	}

	//if pkt.IsKeyFrame && self.sample.SyncSample != nil {
	//	self.sample.SyncSample.Entries = append(self.sample.SyncSample.Entries, uint32(self.sampleIndex+1))
	//}

	//duration := uint32(self.timeToTs(rawdur))
	//if self.sttsEntry == nil || duration != self.sttsEntry.Duration {
	//	self.sample.TimeToSample.Entries = append(self.sample.TimeToSample.Entries, mp4io.TimeToSampleEntry{Duration: duration})
	//	self.sttsEntry = &self.sample.TimeToSample.Entries[len(self.sample.TimeToSample.Entries)-1]
	//}
	//self.sttsEntry.Count++

	//if self.sample.CompositionOffset != nil {
	//	offset := uint32(self.timeToTs(pkt.CompositionTime))
	//	if self.cttsEntry == nil || offset != self.cttsEntry.Offset {
	//		table := self.sample.CompositionOffset
	//		table.Entries = append(table.Entries, mp4io.CompositionOffsetEntry{Offset: offset})
	//		self.cttsEntry = &table.Entries[len(table.Entries)-1]
	//	}
	//	self.cttsEntry.Count++
	//}

	//self.duration += int64(duration)
	//self.sampleIndex++
	//self.sample.ChunkOffset.Entries = append(self.sample.ChunkOffset.Entries, uint32(self.muxer.wpos))
	//self.sample.SampleSize.Entries = append(self.sample.SampleSize.Entries, uint32(len(pkt.Data)))
	//
	//self.muxer.wpos += int64(len(pkt.Data))
	return
}

func (self *Stream) timeToTs(tm time.Duration) int64 {
	return int64(tm * time.Duration(self.timeScale) / time.Second)
}
