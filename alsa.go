package alsa

//Go Packages
import(
	"errors"
)

/*
	Types
*/

type AudioData []byte

type AudioStream struct{
	Channels int
	Rate int
	SampleFormat SampleType
	DataStream chan AudioData
}

/*
	Functions
*/

//Returns channel to send audio streams on
func Init(control <-chan bool) (chan<- AudioStream){
	
	//Create channel (allow two stream buffer)
	//Must be closed by caller
	stream := make(chan AudioStream)
	
	go start(stream,control)

	return stream
}

//Plays by default..send false on control chan to stop
func start(streamChan <-chan AudioStream, control <-chan bool){
	
	var device alsa_device
	
	//End on close of stream chan
	for stream := range streamChan{
		//Create alsa device (will use existing if suitable)
		err := configDevice(&device,&stream)
		if err != nil{
			//fmt.Println("error configuring alsa device") Write to STDERR maybe?
		}
		
		//Delete alsa device
		defer alsa_close(device.pcm)
							
		//Play all data in this stream
		for data := range stream.DataStream{
			select{
				case shouldPlay := <-control:
					for !shouldPlay{
						shouldPlay = <-control //Blocking
					}
				default:
					alsa_write(&device,data)
			}
		}
	}
}

func configDevice(device *alsa_device, stream *AudioStream)(error){

	//Only make new device if one doesn't exist or for different chan_num or rate value
	if device==nil || device.channels != stream.Channels || device.rate != stream.Rate{
		if device.pcm != nil{
			defer alsa_close(device.pcm)
		}
		
		device.channels = stream.Channels
		device.rate = stream.Rate
		
		switch stream.SampleFormat{
		case INT16_TYPE:
			device.numBytes = 2
		default:
			device.numBytes = 1
		}
				
		pcm, err := alsa_open("default",stream.Channels,stream.Rate)
		if err != nil{
			return errors.New("Error on open")
		}
		
		device.pcm = pcm
	}
	
	return nil
}

	

