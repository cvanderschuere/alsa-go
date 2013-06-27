package alsa

import (
	"testing"
	"time"
)

//Tests by playing a sine wave for 15 seconds
func TestAlsaInit(t *testing.T){
	controlChan := make(chan bool)
	streamChan := Init(controlChan)
	
	//Make stream
	dataChan := make(chan AudioData, 100)
	aStream := AudioStream{Channels:2, Rate:4410, SampleFormat:INT16_TYPE, DataStream:dataChan}
	
	//Send stream
	streamChan<-aStream
	
	//Create sample to play
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}  //PI
	
	for i:=0;i<5;i++{
		b  = append(b ,b...)
	}
	
	start := time.Now()
	diff := start.Sub(start)

	//Loop for 15 seconds (With terrible noise)
	for diff.Seconds()<3{
		dataChan<-b
		diff = time.Now().Sub(start);
	}
	
}