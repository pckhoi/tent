package storage

import "log"

type Write struct {
	Row       DataRow
	Overwrite bool
}

func CreateWriteChannel() chan Write {
	channel := make(chan Write)
	go listenForWrite(channel)
	log.Println("Created new write channel")
	return channel
}

func listenForWrite(channel chan Write) {
	for write := range channel {
		log.Printf("Received %+v\n", write)
		go func(_write Write) {
			err := processWrite(_write)
			if err != nil {
				log.Fatal(err)
			}
		}(write)
	}
}

func processWrite(write Write) error {
	fileName := write.Row.getFileName()
	if !fileExist(fileName) || write.Overwrite {
		fieldNames, err := updateFieldNames(write.Row)
		if err != nil {
			return err
		}
		log.Printf("Attempt to write %s\n", fileName)
		err = writeRow(fileName, write.Row.getFileContent(fieldNames))
		if err != nil {
			return err
		}
	}
	return nil
}
