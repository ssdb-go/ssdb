package ssdb

import "time"

// NewCmdResult returns a Cmd initialised with val and err for testing.
func NewCmdResult(val interface{}, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewSliceResult returns a Cmd initialised with val and err for testing.
func NewSliceResult(val []interface{}, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewStatusResult returns a Cmd initialised with val and err for testing.
func NewStatusResult(val string, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewIntResult returns an Cmd initialised with val and err for testing.
func NewIntResult(val int64, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewDurationResult returns a Cmd initialised with val and err for testing.
func NewDurationResult(val time.Duration, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewBoolResult returns a Cmd initialised with val and err for testing.
func NewBoolResult(val bool, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewStringResult returns a Cmd initialised with val and err for testing.
func NewStringResult(val string, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewFloatResult returns a Cmd initialised with val and err for testing.
func NewFloatResult(val float64, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewStringSliceResult returns a Cmd initialised with val and err for testing.
func NewStringSliceResult(val []string, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewBoolSliceResult returns a Cmd initialised with val and err for testing.
func NewBoolSliceResult(val []bool, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewStringStringMapResult returns a StringStringMapCmd initialised with val and err for testing.
func NewMapStringStringResult(val map[string]string, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewMapStringIntCmdResult returns a Cmd initialised with val and err for testing.
func NewMapStringIntCmdResult(val map[string]int64, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}

// NewTimeCmdResult returns a Cmd initialised with val and err for testing.
func NewTimeCmdResult(val time.Time, err error) *Cmd {
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
}
