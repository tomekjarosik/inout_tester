package testcase

import (
	"io"
	"os"
	"path"
)

// Streams provide data for a test case
type Streams struct {
	Input  io.Reader
	Output io.Reader
	Close  func() error
}

type StreamsProvider func(info Info) (Streams, error)

func DirectoryBasedDataStreamsProvider(dir string) StreamsProvider {
	return func(info Info) (Streams, error) {
		streams := Streams{}
		inFile, err := os.OpenFile(path.Join(dir, info.Name+".in"), os.O_RDONLY, 0755)
		if err != nil {
			return streams, err
		}
		goldenOutFile, err := os.OpenFile(path.Join(dir, info.Name+".out"), os.O_RDONLY, 0755)
		if err != nil {
			return streams, err
		}

		streams.Input = inFile
		streams.Output = goldenOutFile

		streams.Close = func() error {
			defer inFile.Close()
			defer goldenOutFile.Close()
			// TODO: Handle closing better
			return nil
		}
		return streams, nil
	}
}
