package internal

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrMissingHistoryFileName = errors.New("file name doesnt exists in history")
)

type HistoryMark uint

const (
	HistoryMarkRemove = iota
	HistoryMarkNew
	HistoryMarkExist
)

type History interface {
	// Start sets all file names with HistoryMarkRemove mark.
	// Should be called on start Parser directory parse
	Start()

	Commit()

	Flush()

	// Add loads new folder data to History
	Add(filePath string)

	GetChanged() (new []string, remove []string)

	LogInfo()
}

type history struct {
	log *zap.Logger

	data   map[string]HistoryMark
	remove []string
}

func (h *history) GetChanged() (new []string, remove []string) {
	new = make([]string, 0)
	remove = make([]string, 0)

	for fileName, fileMark := range h.data {
		switch fileMark {
		case HistoryMarkRemove:
			remove = append(remove, fileName)
		case HistoryMarkNew:
			new = append(new, fileName)
		}
	}

	return new, remove
}

func (h *history) Flush() {
	for _, fileName := range h.remove {
		err := h.removeMark(fileName)
		if err != nil {
			h.log.Error("clear with remove mark err", zap.Error(err))
		}
	}

	h.remove = make([]string, 0)
}

func (h *history) Commit() {
	h.remove = make([]string, 0)

	// adds all HistoryMarkRemove to remove array
	for fileName, historyKey := range h.data {
		if historyKey == HistoryMarkRemove {
			h.remove = append(h.remove, fileName)
		}
	}
}

func (h *history) removeMark(fileName string) error {
	_, ok := h.data[fileName]
	if !ok {
		return errors.Wrapf(ErrMissingHistoryFileName, "file name doesn't exist %s", fileName)
	}

	delete(h.data, fileName)
	return nil
}

func (h *history) LogInfo() {
	h.log.Debug("history", zap.Any("data", h.data))
}

func (h *history) Start() {
	// clear all keys with remove marks
	h.Flush()

	// mark all history keys as removed
	for key := range h.data {
		h.data[key] = HistoryMarkRemove
	}
}

func (h *history) Add(fileName string) {
	// check if file exists
	// do nothing (keep processing)
	_, ok := h.data[fileName]
	if ok {
		h.data[fileName] = HistoryMarkExist
	}

	// check if file doesn't exist
	// add it as new
	_, ok = h.data[fileName]
	if !ok {
		h.data[fileName] = HistoryMarkNew
	}
}

func NewHistory(log *zap.Logger) History {
	return &history{
		log:    log,
		data:   make(map[string]HistoryMark),
		remove: make([]string, 0),
	}
}
