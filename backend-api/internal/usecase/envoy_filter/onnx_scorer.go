package envoy_filter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"gorgonia.org/tensor"
)

type ONNXMLScorer struct {
	mu       sync.RWMutex
	sessions map[string]*onnxModelSession
}

type onnxModelSession struct {
	model *onnx.Model
	graph *gorgonnx.Graph

	inputRank      int
	inputFixedDims []int

	mu sync.Mutex
}

func NewONNXMLScorer() (*ONNXMLScorer, error) {
	return &ONNXMLScorer{
		sessions: make(map[string]*onnxModelSession),
	}, nil
}

func (s *ONNXMLScorer) Score(_ context.Context, modelName string, modelData []byte, featureVector []float32) (float32, error) {
	if len(modelData) == 0 {
		return 0, errors.New("onnx model data is empty")
	}
	if len(featureVector) == 0 {
		return 0, errors.New("feature vector is empty")
	}

	session, err := s.getOrCreateSession(modelName, modelData)
	if err != nil {
		return 0, err
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	inputShape, err := session.buildInputShape(len(featureVector))
	if err != nil {
		return 0, err
	}
	inputTensor := tensor.New(
		tensor.WithShape(inputShape...),
		tensor.WithBacking(featureVector),
	)

	if err := session.model.SetInput(0, inputTensor); err != nil {
		return 0, fmt.Errorf("failed to set model input: %w", err)
	}
	if err := session.graph.Run(); err != nil {
		return 0, fmt.Errorf("onnx inference failed: %w", err)
	}

	outputs, err := session.model.GetOutputTensors()
	if err != nil {
		return 0, fmt.Errorf("failed to read model output tensors: %w", err)
	}
	if len(outputs) != 1 {
		return 0, fmt.Errorf("model %q must return exactly 1 output tensor, got %d", modelName, len(outputs))
	}

	return extractScoreFromTensor(outputs[0])
}

func (s *ONNXMLScorer) getOrCreateSession(modelName string, modelData []byte) (*onnxModelSession, error) {
	hash := hashModelData(modelData)

	s.mu.RLock()
	existing := s.sessions[hash]
	s.mu.RUnlock()
	if existing != nil {
		return existing, nil
	}

	created, err := createONNXModelSession(modelName, modelData)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if cached := s.sessions[hash]; cached != nil {
		return cached, nil
	}
	s.sessions[hash] = created
	return created, nil
}

func createONNXModelSession(modelName string, modelData []byte) (*onnxModelSession, error) {
	graph := gorgonnx.NewGraph()
	model := onnx.NewModel(graph)
	if err := model.UnmarshalBinary(modelData); err != nil {
		return nil, fmt.Errorf("failed to decode onnx model %q: %w", modelName, err)
	}

	inputs := model.GetInputTensors()
	if len(inputs) != 1 {
		return nil, fmt.Errorf("model %q must have exactly 1 input, got %d", modelName, len(inputs))
	}

	inputShape := inputs[0].Shape()
	if len(inputShape) == 0 {
		return nil, fmt.Errorf("model %q has empty input shape", modelName)
	}
	if len(inputShape) > 2 {
		return nil, fmt.Errorf("unsupported input rank %d for model %q", len(inputShape), modelName)
	}

	fixedDims := make([]int, len(inputShape))
	copy(fixedDims, inputShape)

	return &onnxModelSession{
		model:          model,
		graph:          graph,
		inputRank:      len(inputShape),
		inputFixedDims: fixedDims,
	}, nil
}

func (s *onnxModelSession) buildInputShape(featureCount int) ([]int, error) {
	if s.inputRank == 1 {
		dim := s.inputFixedDims[0]
		if dim > 0 && dim != featureCount {
			return nil, fmt.Errorf("model input size mismatch: expected %d, got %d", dim, featureCount)
		}
		return []int{featureCount}, nil
	}

	if s.inputRank == 2 {
		batch := s.inputFixedDims[0]
		width := s.inputFixedDims[1]
		if batch > 0 && batch != 1 {
			return nil, fmt.Errorf("unsupported model batch size %d: only 1 is supported", batch)
		}
		if width > 0 && width != featureCount {
			return nil, fmt.Errorf("model input width mismatch: expected %d, got %d", width, featureCount)
		}
		return []int{1, featureCount}, nil
	}

	return nil, fmt.Errorf("unsupported input rank %d", s.inputRank)
}

func extractScoreFromTensor(t tensor.Tensor) (float32, error) {
	if t == nil {
		return 0, errors.New("onnx output tensor is nil")
	}

	switch data := t.Data().(type) {
	case []float32:
		if len(data) == 0 {
			return 0, errors.New("onnx output tensor is empty")
		}
		return data[0], nil
	case []float64:
		if len(data) == 0 {
			return 0, errors.New("onnx output tensor is empty")
		}
		return float32(data[0]), nil
	default:
		return 0, fmt.Errorf("unsupported onnx output backing type %T", data)
	}
}

func hashModelData(modelData []byte) string {
	sum := sha256.Sum256(modelData)
	return hex.EncodeToString(sum[:])
}
