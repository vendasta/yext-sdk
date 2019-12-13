package yext

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrorType string

const (
	ErrorTypeFatal    = "FATAL_ERROR"
	ErrorTypeNonFatal = "NON_FATAL_ERROR"
	ErrorTypeWarning  = "WARNING"
)

type Error struct {
	Message     string `json:"message"`
	Code        int    `json:"code"`
	Type        string `json:"type"`
	RequestUUID string `json:"request_uuid"`
}

func (e Error) Error() string {
	return fmt.Sprintf("type: %s code: %d message: %s, request uuid: %s", e.Type, e.Code, e.Message, e.RequestUUID)
}

func (e Error) ErrorWithoutUUID() string {
	return fmt.Sprintf("type: %s code: %d message: %s", e.Type, e.Code, e.Message)
}

func (e Error) IsError() bool {
	return e.Type == ErrorTypeFatal || e.Type == ErrorTypeNonFatal
}

func (e Error) IsWarning() bool {
	return e.Type == ErrorTypeWarning
}

type Errors []*Error

func (e Errors) Error() string {
	var (
		errs = make([]string, len(e))
		uuid = ""
	)

	for i, err := range e {
		errs[i] = err.ErrorWithoutUUID()
		uuid = err.RequestUUID
	}

	return fmt.Sprintf("%s; request uuid: %s", strings.Join(errs, "; "), uuid)
}

func (e Errors) Errors() []*Error {
	var errors []*Error
	for _, err := range e {
		if err.IsError() {
			errors = append(errors, err)
		}
	}
	return errors
}

func (e Errors) Warnings() []*Error {
	var warnings []*Error
	for _, err := range e {
		if err.IsWarning() {
			warnings = append(warnings, err)
		}
	}
	return warnings
}

// GetNumErrors returns the number of errors (warnings are excluded)
func GetNumErrors(err error) int {
	if err == nil {
		return 0
	}

	if e, ok := err.(Errors); ok {
		return len(e.Errors())
	}

	if e, ok := err.(Error); ok {
		if e.IsError() {
			return 1
		}
		return 0
	}

	return 1
}

//ToUserFriendlyMessage Will return a string describing the error that can be displayed to end users
func ToUserFriendlyMessage(err error) string {
	var message string
	if e, ok := err.(Errors); ok {
		for _, innerError := range e {
			if message != "" {
				message = message + ", "
			}
			message = message + ToUserFriendlyMessage(innerError)
		}
	} else if e, ok := err.(*Error); ok {
		message = e.Message
	} else if e, ok := err.(Error); ok {
		message = e.Message
	} else {
		message = err.Error()
	}
	return message
}

func IsNotFoundError(err error) bool {
	if e, ok := err.(Errors); ok {
		for _, innerError := range e {
			if IsNotFoundError(innerError) {
				return true
			}
		}
	} else if e, ok := err.(*Error); ok {
		if e.Code == 2000 || e.Code == 6004 || e.Code == 2238 {
			return true
		}
	}
	return false
}

//IsBusinessError Returns true if the validation/processing that takes place on Yext's servers failed
// It will return false if the server could not be reached or other protocol error occurs
func IsBusinessError(err error) bool {
	if e, ok := err.(Errors); ok {
		for _, innerError := range e {
			if IsBusinessError(innerError) {
				return true
			}
		}
	} else if _, ok := err.(*Error); ok {
		return true
	}
	return false
}

func IsFatalBusinessError(err error) bool {
	if e, ok := err.(Errors); ok {
		for _, innerError := range e {
			if IsFatalBusinessError(innerError) {
				return true
			}
		}
	} else if e, ok := err.(*Error); ok {
		if e.Type == ErrorTypeFatal {
			return true
		}
	}
	return false
}

func IsErrorCode(err error, code int) bool {
	if e, ok := err.(Errors); ok {
		for _, innerError := range e {
			if IsErrorCode(innerError, code) {
				return true
			}
		}
	} else if e, ok := err.(*Error); ok {
		if e.Code == code {
			return true
		}
	}
	return false
}

func splitStrAtWord(str string, word string) (string, string) {
	var (
		words  = strings.Split(str, " ")
		found  = false
		before = ""
		after  = ""
	)
	for _, w := range words {
		if w == word {
			found = true
		} else if found {
			if after != "" {
				after += " "
			}
			after += w
		} else {
			if before != "" {
				before += " "
			}
			before += w
		}
	}
	return before, after
}

func errorFromString(str string) (*Error, error) {
	strRemaining := strings.TrimLeft(str, "type: ")
	typ, strRemaining := splitStrAtWord(strRemaining, "code:")
	code, message := splitStrAtWord(strRemaining, "message:")
	codeInt, errConv := strconv.Atoi(code)
	if errConv != nil {
		return nil, errConv
	}
	return &Error{Type: typ, Code: codeInt, Message: message}, nil
}

func ErrorsFromString(errorStr string) ([]*Error, error) {
	errStrList := strings.Split(errorStr, "; ")
	var errors []*Error
	uuid := strings.TrimLeft(errStrList[len(errStrList)-1], "request uuid: ")
	for i := 0; i < len(errStrList)-1; i++ {
		errObj, err := errorFromString(errStrList[i])
		if err != nil {
			return nil, err
		}
		errObj.RequestUUID = uuid
		errors = append(errors, errObj)
	}
	return errors, nil
}
