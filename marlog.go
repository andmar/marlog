package marlog

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// NOTE: Refer to Log package documentation for the meaning of these, these are just poiting to them.
	// FlagLdate ...
	FlagLdate = log.Ldate
	// FlagLtime ...
	FlagLtime = log.Ltime
	// FlagLmicroseconds ...
	FlagLmicroseconds = log.Lmicroseconds
	// FlagLlongfile ...
	FlagLlongfile = log.Llongfile
	// FlagLshortfile ...
	FlagLshortfile = log.Lshortfile
	// FlagLUTC ...
	FlagLUTC = log.LUTC
	// FlagLstdFlags ...
	FlagLstdFlags = log.Ldate | log.Ltime
)

// MarLog ...
var MarLog *MarLogger

// MarLogger ...
type MarLogger struct {
	Prefix        string
	Flags         int
	stamps        map[string]*stamp
	outputHandles map[string]*outputHandle
}

// Stamp ...
type stamp struct {
	Name          string // NOTE: Should be the same as the key in the stamps map
	Active        bool
	MessagePrefix string
	HandleKeys    []string // NOTE: Optionally filter the outputHandlers to use. If this is empty the system will use all output handles to output messages with this Stamp
}

type outputHandle struct {
	Name   string // NOTE: Should be the same as the key in the stamps map
	handle io.Writer
}

// Log ...
func (logger *MarLogger) Log(stampName string, message string, newLine bool, fatal bool) error {

	if _, found := logger.stamps[stampName]; found == false {
		return fmt.Errorf("Stamp named \"%s\" does not exist.", stampName)
	}

	stamp := logger.stamps[stampName]

	if stamp.Active == false {
		return fmt.Errorf("Stamp named \"%s\" is not active.", stampName)
	}

	for _, currentHandleKey := range stamp.HandleKeys {
		outputHandle, found := logger.outputHandles[currentHandleKey]
		if found == false {
			return fmt.Errorf("Output Handle named \"%s\" not found.", currentHandleKey)
		}
		log := log.New(outputHandle.handle, stamp.MessagePrefix, logger.Flags)

		if newLine == true {
			log.Println(message)
		} else {
			log.Printf(message)
		}

	}

	if fatal == true { // NOTE: This is kinda the same as using Golang's log.Fatalf
		os.Exit(-1)
	}

	return nil

}

// AddStamp ...
func (logger *MarLogger) AddStamp(stampName string, outputHandleKeys ...string) error {

	if _, found := logger.stamps[stampName]; found == true {
		return fmt.Errorf("Stamp named \"%s\" already exists.", stampName)
	}

	var handleKeys []string
	for _, currentHandleKey := range outputHandleKeys {
		handleKeys = append(handleKeys, currentHandleKey)
	}

	newStamp := new(stamp)
	newStamp.Name = stampName
	newStamp.Active = true
	newStamp.MessagePrefix = ""
	newStamp.HandleKeys = handleKeys
	logger.stamps[stampName] = newStamp

	return nil
}

// AddOutputHandle ...
func (logger *MarLogger) AddOutputHandle(handleName string, handle io.Writer) error {

	if _, found := logger.outputHandles[handleName]; found == true {
		return fmt.Errorf("Output Handle named \"%s\" already exists.", handleName)
	}

	newOutputHandle := new(outputHandle)
	newOutputHandle.Name = handleName
	newOutputHandle.handle = handle
	logger.outputHandles[handleName] = newOutputHandle

	return nil
}

// ActivateStamps ...
func (logger *MarLogger) ActivateStamps(stampNames ...string) error {

	if len(stampNames) == 0 {
		return fmt.Errorf("No Stamp names provided")
	}

	for _, currentStampName := range stampNames {
		if _, found := logger.stamps[currentStampName]; found == false {
			return fmt.Errorf("Nothing done. Stamp named \"%s\" does not exist.", currentStampName)
		}
	}

	for _, currentStampName := range stampNames {
		logger.stamps[currentStampName].Active = true
	}

	return nil
}

// DeactivateStamps ...
func (logger *MarLogger) DeactivateStamps(stampNames ...string) error {

	if len(stampNames) == 0 {
		return fmt.Errorf("No Stamp names provided")
	}

	for _, currentStampName := range stampNames {
		if _, found := logger.stamps[currentStampName]; found == false {
			return fmt.Errorf("Nothing done. Stamp named \"%s\" does not exist.", currentStampName)
		}
	}

	for _, currentStampName := range stampNames {
		logger.stamps[currentStampName].Active = false
	}

	return nil
}

func init() {
	MarLog = new(MarLogger)
	MarLog.Prefix = ""
	MarLog.Flags = FlagLdate | FlagLtime | FlagLshortfile
	MarLog.stamps = make(map[string]*stamp)
	MarLog.outputHandles = make(map[string]*outputHandle)
}
