package capture

import (
	"crypto/rand"
	mock_screenshot "duetfeature/pkg/screenshot/mock"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func TestCapturerRun(t *testing.T) {
	// Create channels for communication
	outputChan := make(chan []byte)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocking the worker
	m := mock_screenshot.NewMockImageGeneratorInterface(ctrl)

	// We expect the worker to be called at least once
	m.
		EXPECT().
		Capture().
		DoAndReturn(func() ([][]byte, error) {
			time.Sleep(1 * time.Second)
			size := 16 // Define the size of the byte array
			bytes := make([]byte, size)
			_, err := rand.Read(bytes)
			if err != nil {
				panic(err)
			}

			return [][]byte{bytes}, nil
		}).
		MinTimes(1)

	// Init capturer module
	captureModule := NewCapturer(WithWorker(
		m),
		WithOutputChan(outputChan),
		WithInterval(2),
	)
	captureModule.Start()

	// Test cases
	select {
	case data := <-outputChan:
		if reflect.TypeOf(data).Kind() != reflect.Slice {
			t.Errorf("Expected a slice, got %v", reflect.TypeOf(data).Kind())
		}
		if len(data) == 0 {
			t.Errorf("Received empty slice")
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timed out waiting for data")
	}

	// Stop the capture module
	captureModule.Stop()
	close(outputChan)
}
