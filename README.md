


## DuetFeature Showcase

This is an experimental feature that captures display and processes its content through OCR to form user's knowledge database.

There are two main components:

1. `internal/usecase/capture` Capture module that captures display content every  n seconds.
2. `internal/usecase/ocr` OCR module that converts image to text.

These two components are communicating via a channel.

The purpose of this project is to showcase coding style, architectural patterns and approaches that are relevant to micro-service architecture.

## Architecture
In recent years, various system architecture patterns or principles like Clean Architecture, SOLID, and GRASP, etc have emerged. They all have been popular among high performant engineering teams and they all have the same objective, which is the separation of concerns. Below I'll focus on SOLID's **OCP** and **DIP** as I believe they are the most important.

```
If the OCP states the goal of OO architecture,
the DIP states the primary mechanism (Robert C. Martin)
```

### Dependency Inversion Principle
For both modules (`internal/usecase/capture` , `internal/usecase/ocr`) I use an interface as an intermediary to decouple underlying components (currently `pkg/gocr` and `pkg/screenshot`). This makes it easier to swap them out or perform testing using mocks/stubs.

Here is an example of **Capture** initializer:
```go
// Screen Capture
capt := capture.NewCapturer(
	capture.WithWorker(&screenshot.ImageGenerator{}),   // This is underlying dependency
	capture.WithOutputChan(imageRequestChan),
	capture.WithInterval(2),
)
```

### Open-Closed Principle
OCP helps to avoid making changes in existing well-tested code, reducing the risk of introducing bugs and issues while developing. Let's say this code is already in production, but we want to add a new feature for **Capture** and **OCR** modules which is saving a data produced into a file. One way to do that is to create SaveToFile methods in for both Capturer and OCR modules. However, modifying the Capturer struct/class itself would violate the Open-Closed Principle, and we might need to avoid it.  *(it will as well violate The Single Responsibility Principle as a class should have only one reason to change)*.

Two common ways to extend the functionality of a class without modifying it are **Composition** and **Visitor pattern**. In this case, I'll use the **Visitor pattern**.

First, let's  implement file saving module `internal/datasaver`:
```go
type FileSaver struct {
	Filename string
}

func (fs *FileSaver) VisitCapture(c domain.CapturerInterface) {
	fmt.Println("file saved to disk...")
}
```

Modifications to the original code will be minimal; we need to define a new Visitor interface and implement Accept method:

```go
type CapturerInterface interface {
	Accept(Visitor)
}

type Visitor interface {
	VisitCapture(CapturerInterface)
}
...

func (s *Capturer) Accept(v domain.Visitor) {
	v.VisitCapture(s)
}
```

Now we can call the new functional the following way:
```go
capt.Accept(&datasaver.FileSaver{Filename: "test-filename.jpg"})
ocrInst.Accept(&datasaver.FileSaver{Filename: "test-filename.jpg"})
```

### Functional Options Pattern
Config struct or params are the most common way to pass configuration to a class instance. However, it may introduce breaking changes to the Config struct when new options are added or old ones are being removed. A better alternative is functional options design pattern.

In this app the OCR module is initialized by passing params and screen capture is initialized by passing functional options.

Example without functional options:
```go
// OCR
ocrInst := ocr.NewProcessor(gocr.NewGoCr(), imageRequestChan)
ocrInst.Start()
```

Example with functional options:
```
// Screen Capture
capt := capture.NewCapturer(
	capture.WithWorker(&screenshot.ImageGenerator{}),
	capture.WithOutputChan(imageRequestChan),
	capture.WithInterval(2),
)
```

## Testing

### Capture module

Let's look at how we can test `internal/usecase/capture` module which function is to periodically submit user's screen data.  Since it relies on `screnshot` package to make scherenshots, i'll mock it using `mockgen` package and test `internal/usecase/capture` and `screnshot` independently. 

Here is a test case from `internal/usecase/capture/capture_test.go`. It simply verifies that data is supplied to the selected channel.

```go
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
```
### Screenshot module
While this module has been mocked in previous testing to isolate it from `Capture` it still needs to be tested to make sure screenshots are successfully created. If this was a real application, we would have to use Use VMs and containers to simulate different OS environments, because this functional might fail on certain OS configurations.

Here is a test code example from `pkg/screenshot/screenshot_test.go`:

```go
package screenshot

import (
	"github.com/kbinani/screenshot"
	"testing"
)

func TestCaptureNotEmpty(t *testing.T) {
	gen := ImageGenerator{}

	capture, err := gen.Capture()
	if err != nil {
		t.Errorf("Failed to capture screen: %s", err)
		return
	}

	for i, imgBytes := range capture {
		if len(imgBytes) == 0 {
			t.Errorf("Capture for display %d is empty", i)
		}
	}

}

func TestIfActiveDisplaysFound(t *testing.T) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		t.Error("No active displays found")
		return
	}
}
```


## Future work
- Add swagger library for docs generation
- Explain how to implement logging and monitoring