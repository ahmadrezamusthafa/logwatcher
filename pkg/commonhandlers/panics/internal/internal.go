package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/logger"
	"os"
)

var (
	env           string
	filepath      string
	tagString     string
	customMessage string

	title      string
	errUnknown = errors.New("Unknown error")
)

func getError(msg interface{}) error {
	return parseError(msg)
}

func parseError(msg interface{}) error {
	if msg != nil {
		switch err := msg.(type) {
		case string:
			return errors.New(err)
		case error:
			return err
		default:
			return errUnknown
		}
	}
	return nil
}

func HandlePanic(action func(error)) {
	if err := getError(recover()); err != nil {
		action(err)
	}
}

func PublishError(errs error, detail []byte, dump []byte) {
	var message bytes.Buffer

	message.WriteString(title)
	message.WriteString(errs.Error())
	message.WriteString(tagString)
	message.WriteString(customMessage)
	message.WriteString("\n")

	if len(detail) != 0 {
		message.WriteString("```")
		message.Write(detail)
		message.WriteString("```\n")
	}

	if filepath != "" {
		go func() {
			fp := fmt.Sprintf("%s/panics.log", filepath)
			file, err := os.OpenFile(fp, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				logger.Err("[Interceptor][PublishError] Failed to open file %s", fp)
				return
			}
			file.Write(message.Bytes())
			if len(dump) != 0 {
				file.Write(dump)
				file.WriteString("\r\n")
			}
			file.Close()
		}()
	}
}

var (
	panicKeyword   = []byte("src/runtime/panic.go")
	sorabelKeyword = []byte("github.com/salestock")
	captureKeyword = []byte("/github.com/salestock/ssource/segmentation/ss-libs/commonhandlers/panics/")
)

func TrimStackTrace(stackTrace []byte) []byte {
	panicIdx := bytes.Index(stackTrace, panicKeyword)
	if panicIdx == -1 {
		panicIdx = 0
	}
	stack := stackTrace[panicIdx:]
	panicIdx = bytes.Index(stack, sorabelKeyword)
	if panicIdx == -1 {
		panicIdx = 0
	}
	stack = stack[panicIdx:]
	captureIdx := bytes.Index(stack, captureKeyword)
	if captureIdx != -1 {
		newlineIdx := bytes.Index(stack[captureIdx:], []byte("\n"))
		if newlineIdx != -1 {
			captureIdx += newlineIdx + 1
			return stack[:captureIdx]
		}
	}
	return stack
}
