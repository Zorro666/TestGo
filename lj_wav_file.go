package lj_wav_file

import "os"
import "log"
import "./bufferOutput"

const LJ_WAV_FORMAT_PCM uint16 = 0x1
const LJ_WAV_FORMAT_IEEE_FLOAT uint16 = 0x3
const LJ_WAV_FORMAT_ALAW uint16 = 0x6
const LJ_WAV_FORMAT_MULAW uint16 = 0x7
const LJ_WAV_FORMAT_EXTENSIBLE uint16 = 0xFFFE

const LJ_WAV_OK int = 0
const LJ_WAV_ERROR int = -1

/*
RIFF_HEADER: little-endian
{
riffID								Bytes 4	: Chunk ID: "RIFF"
chunkSize							Bytes 4	: Chunk size: 4+n
wavID									Bytes 4	: WAVE ID: "WAVE"
}
chunkData							Bytes n : chunk data e.g. WAV chunk
*/

/*
WAV_HEADER: little-endian
{
chunkID								Bytes 4		: Chunk ID: "fmt "
chunkSize							Bytes 4		: Chunk size: 16 or 18 or 40
wFormatTag						Bytes 2		: Format code
nChannels							Bytes 2		: Number of interleaved channels
nSamplesPerSec				Bytes 4		: Sampling rate (blocks per second)
nAvgBytesPerSec				Bytes 4		: Data rate
nBlockAlign						Bytes 2		: Data block size (bytes)
wBitsPerSample				Bytes 2		: Bits per sample
}
// chunkSize = 18
{
cbSize								Bytes 2		: Size of the extension (0 or 22)
}
// chunkSize = 40
{
wValidBitsPerSample		Bytes 2		: Number of valid bits
dwChannelMask					Bytes 4		: Speaker position mask
SubFormat							Bytes 16	: GUID, including the data format code
}

wFormatTag: Format Codes:
0x0001	WAVE_FORMAT_PCM					PCM
0x0003	WAVE_FORMAT_IEEE_FLOAT	IEEE float
0x0006	WAVE_FORMAT_ALAW				8-bit	ITU-T G.711 A-law
0x0007	WAVE_FORMAT_MULAW				8-bit	ITU-T G.711 ?-law
0xFFFE	WAVE_FORMAT_EXTENSIBLE	Determined by SubFormat

PCM Format

The first part of the Format chunk is used to describe PCM data.

For PCM data, the Format chunk in the header declares the number of bits/sample in each sample (wBitsPerSample). The original documentation (Revision 1) specified that the number of bits per sample is to be rounded up to the next multiple of 8 bits. This rounded-up value is the container size. This information is redundant in that the container size (in bytes) for each sample can also be determined from the block size divided by the number of channels (nBlockAlign / nChannels).
This redundancy has been appropriated to define new formats. For instance, Cool Edit uses a format which declares a sample size of 24 bits together with a container size of 4 bytes (32 bits) determined from the block size and number of channels. With this combination, the data is actually stored as 32-bit IEEE floats. The normalization (full scale 223) is however different from the standard float format.
PCM data is two's-complement except for resolutions of 1-8 bits, which are represented as offset binary.

*/

func lj_makefourcc(ch0 byte, ch1 byte, ch2 byte, ch3 byte) (result uint32) {
	result = (uint32(ch0) << 0) | (uint32(ch1) << 8) | (uint32(ch2) << 16) | (uint32(ch3) << 24)
	return
}

type lj_wav_riff_header struct {
	chunkID uint32									// "RIFF"
	chunkSize uint32								// 4+chunkDataSize
	wavID uint32										// "WAVE"
}

func lj_wav_riff_header_size() (size int) {
	size = 4+4+4
	return
}

func (riffHeader *lj_wav_riff_header) write(file *os.File) (ok bool, err os.Error) {
	var buffer bufferOutput.BufferOutput
	buffer.Create(128)
	//chunkID uint32								// "RIFF"
	buffer.AddUint32(riffHeader.chunkID)
	//chunkSize uint32							// 4+chunkDataSize
	buffer.AddUint32(riffHeader.chunkSize)
	//wavID uint32									// "WAVE"
	buffer.AddUint32(riffHeader.wavID)

	ok = false
	n, err := buffer.Write(file)
	if n == lj_wav_riff_header_size() {
		ok = true
	}
	return
}

type lj_wav_format_header struct {
	chunkID uint32									// "fmt "
	chunkSize uint32								// 16 
	wFormatTag uint16								// PCM=0x1, IEEE_FLOAT=0x3, ALAW=0x6, MULAW=0x7, EXTENSIBLE = 0xFFFE
	nChannels uint16								// number of interleaved channels
	nSamplesPerSec uint32						// sampling rate (blocks per second)
	nAvgBytesPerSec uint32					// data rate (bytes per second)
	nBlockAlign uint16							// data block size in bytes
	wBitsPerSample uint16						// bits per sample
}

func lj_wav_format_header_size() (size int) {
	size = (4+4+2+2+4+4+2+2)
	return
}

func (formatHeader *lj_wav_format_header) write(file *os.File) (ok bool, err os.Error) {
	var buffer bufferOutput.BufferOutput
	buffer.Create(128)
	//chunkID uint32								// "fmt "
	buffer.AddUint32(formatHeader.chunkID)
	//chunkSize uint32							// 16
	buffer.AddUint32(formatHeader.chunkSize)
	//wFormatTag uint16							// PCM=0x1, IEEE_FLOAT=0x3, ALAW=0x6, MULAW=0x7, EXTENSIBLE = 0xFFFE
	buffer.AddUint16(formatHeader.wFormatTag)
	//nChannels uint16							// number of interleaved channels
	buffer.AddUint16(formatHeader.nChannels)
	//nSamplesPerSec uint32					// sampling rate (blocks per second)
	buffer.AddUint32(formatHeader.nSamplesPerSec)
	//nAvgBytesPerSec uint32				// data rate (bytes per second)
	buffer.AddUint32(formatHeader.nAvgBytesPerSec)
	//nBlockAlign uint16						// data block size in bytes
	buffer.AddUint16(formatHeader.nBlockAlign)
	//wBitsPerSample uint16					// bits per sample
	buffer.AddUint16(formatHeader.wBitsPerSample)

	ok = false
	n, err := buffer.Write(file)
	if n == lj_wav_format_header_size() {
		ok = true
	}
	return
}

type lj_wav_data_header struct {
	chunkID uint32									// "data"
	chunkSize uint32								// wBitsPerSamples * 8 * nChannels * numSamples
}

func lj_wav_data_header_size() (size int) {
	size = (4+4)
	return
}

func (dataHeader *lj_wav_data_header) write(file *os.File) (ok bool, err os.Error) {
	var buffer bufferOutput.BufferOutput
	buffer.Create(128)
	//chunkID uint32								// "data"
	buffer.AddUint32(dataHeader.chunkID)
	//chunkSize uint32							// wBitsPerSamples * 8 * nChannels * numSamples
	buffer.AddUint32(dataHeader.chunkSize)

	ok = false
	n, err := buffer.Write(file)
	if n == lj_wav_data_header_size() {
		ok = true
	}
	return
}

//////////////////////////////////////////////////////////////////////////////
// 
// External structures
// 
//////////////////////////////////////////////////////////////////////////////

type LJ_WAV_FILE struct {
	file *os.File
	format uint16
	numChannels uint32
	sampleRate uint32
	numBytesPerChannel uint32
	numBytesWritten uint32
}

//////////////////////////////////////////////////////////////////////////////
// 
// External Data and functions
// 
//////////////////////////////////////////////////////////////////////////////

func LJ_WAV_create(filename string, format uint16, numChannels uint32, sampleRate uint32, numBytesPerChannel uint32) (wavFile *LJ_WAV_FILE) {
	var riffHeader lj_wav_riff_header
	var wavHeader lj_wav_format_header
	var dataHeader lj_wav_data_header
	var ok bool = false

	wavFile = &LJ_WAV_FILE{}
	if wavFile == nil {
		log.Printf("LJ_WAV_create: &LJ_WAV_FILE failed '%s'\n", filename)
		return
	}
	file, err := os.Create(filename)

	if file == nil || err != nil {
		log.Printf("LJ_WAV_create: failed to open output file '%s' error:%s\n", filename, err.String())
		file.Close()
		return
	}

	riffHeader.chunkID = lj_makefourcc('R','I','F','F')
	riffHeader.chunkSize = 0 // 4 + (8 + 16) + (8 + (wBitsPerSamples * 8 * nChannels * numSamples))
	riffHeader.wavID = lj_makefourcc('W','A','V','E')
	ok, err = riffHeader.write(file)
	if ok == false || err != nil {
		log.Printf("LJ_WAV_create: failed to write WAV_RIFF_HEADER file '%s' error:%s\n", filename, err.String())
		file.Close()
		return nil
	}

	wavHeader.chunkID = lj_makefourcc('f','m','t',' ')
	wavHeader.chunkSize = 16
	wavHeader.wFormatTag = format
	wavHeader.nChannels = uint16(numChannels)
	wavHeader.nSamplesPerSec = uint32(sampleRate)
	wavHeader.nAvgBytesPerSec = uint32(sampleRate) * uint32(numBytesPerChannel) * uint32(numChannels)
	wavHeader.nBlockAlign = uint16(numBytesPerChannel * numChannels)
	wavHeader.wBitsPerSample = uint16(8 * numBytesPerChannel)
	ok, err = wavHeader.write(file)
	if ok == false || err != nil {
		log.Printf("LJ_WAV_create: failed to write WAV_FORMAT_HEADER file '%s' error:%s\n", filename, err.String())
		file.Close()
		return nil
	}

	dataHeader.chunkID = lj_makefourcc('d','a','t','a')
	dataHeader.chunkSize = 0 // wBitsPerSamples * 8 * nChannels * numSamples
	ok, err = dataHeader.write(file)
	if ok == false || err != nil {
		log.Printf("LJ_WAV_create: failed to write WAV_DATA_HEADER file '%s' error:%s\n", filename, err.String())
		file.Close()
		return nil
	}

	wavFile.format = format
	wavFile.numChannels = numChannels
	wavFile.sampleRate = sampleRate
	wavFile.numBytesPerChannel = numBytesPerChannel
	wavFile.numBytesWritten = 0
	wavFile.file = file

	return wavFile
}

func LJ_WAV_FILE_writeChannel(wavFile *LJ_WAV_FILE, sampleData uint16) (result int) {
	result = LJ_WAV_ERROR

	if wavFile == nil {
		log.Printf("LJ_WAV_writeChannel: wavFile is nil\n")
		return
	}

	if wavFile.numBytesPerChannel != 2 {
		log.Printf("LJ_WAV_FILE_writeChannel ERROR only support 2 bytes per channel not:%d", wavFile.numBytesPerChannel)
		return
	}

	var buffer bufferOutput.BufferOutput
	buffer.Create(2)
	buffer.Reset()
	buffer.AddUint16(sampleData)
	buffer.Write(wavFile.file)
	wavFile.numBytesWritten += wavFile.numBytesPerChannel

	return LJ_WAV_OK
}

func LJ_WAV_close(wavFile *LJ_WAV_FILE) (result int) {
	result = LJ_WAV_ERROR

	if wavFile == nil {
		log.Printf("LJ_WAV_close: wavFile is nil\n")
		return
	}

	defer	wavFile.file.Close()

	var buffer bufferOutput.BufferOutput
	buffer.Create(4)

	var dataChunkSize uint32 = wavFile.numBytesWritten
	var riffChunkSize uint32 = 4 + (8 + 16) + (8 + dataChunkSize)

	//Seek back and update riff chunkSize and data chunkSize 
	var riffChunkSizeOffset int64 = 4
	_, err := wavFile.file.Seek(riffChunkSizeOffset, 0)
	if err != nil {
		log.Printf("LJ_WAV_close: seek riffChunkSize:%d failed error:%s\n", riffChunkSizeOffset, err.String())
		return
	}
	buffer.Reset()
	buffer.AddUint32(riffChunkSize)
	n, err := buffer.Write(wavFile.file)
	if n != 4 {
		log.Printf("LJ_WAV_close: error writing riffChunkSize:%d at offset:%d bad number of bytes written:%d expected:4\n", 
							 riffChunkSize, riffChunkSizeOffset, n)
		return
	}
	if err != nil {
		log.Printf("LJ_WAV_close: error writing riffChunkSize:%d at offset:%d error:%s\n", riffChunkSize, riffChunkSizeOffset, err.String())
		return
	}

	var dataChunkSizeOffset int64 = int64(lj_wav_riff_header_size() + lj_wav_format_header_size() + 4)
	_, err = wavFile.file.Seek(dataChunkSizeOffset, 0)
	if err != nil {
		log.Printf("LJ_WAV_close: seek dataChunkSize:%d failed error:%s\n", dataChunkSizeOffset, err.String())
		return
	}
	buffer.Reset()
	buffer.AddUint32(dataChunkSize)
	n, err = buffer.Write(wavFile.file)
	if n != 4 {
		log.Printf("LJ_WAV_close: error writing dataChunkSize:%d at offset:%d bad number of bytes written:%d expected:4\n", 
							 dataChunkSize, dataChunkSizeOffset, n)
		return
	}
	if err != nil {
		log.Printf("LJ_WAV_close: error writing dataChunkSize:%d at offset:%d error:%s\n", dataChunkSize, dataChunkSizeOffset, err.String())
		return
	}

	return LJ_WAV_OK
}
