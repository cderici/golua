package runtime

type RuntimeContext interface {
	CpuLimit() uint64
	CpuUsed() uint64

	MemLimit() uint64
	MemUsed() uint64

	Status() RuntimeContextStatus
	Parent() RuntimeContext

	Flags() RuntimeContextFlags

	SafetyFlags() RuntimeSafetyFlags
}

type RuntimeContextStatus uint8

const (
	RCS_Live RuntimeContextStatus = iota
	RCS_Done
	RCS_Killed
)

type RuntimeContextFlags uint8

const (
	RCF_Empty RuntimeContextFlags = 1 << iota
	RCF_NoIO
	RCF_NoGoLib
)

func (f RuntimeContextFlags) IsSet(ctx RuntimeContext) bool {
	return f&ctx.Flags() != 0
}

var ErrIODisabled = NewErrorS("io disabled")
var ErrGoBridgeDisabled = NewErrorS("go disabled")
var ErrNotSafe = NewErrorS("not safe")

func (r *Runtime) CheckIO() *Error {
	if RCF_NoIO.IsSet(r) {
		return ErrIODisabled
	}
	return nil
}

func (r *Runtime) CheckGoLib() *Error {
	if RCF_NoGoLib.IsSet(r) {
		return ErrGoBridgeDisabled
	}
	return nil
}

type RuntimeSafetyFlags uint16

const (
	RCS_MemSafe RuntimeSafetyFlags = 1 << iota
	RCS_CpuSafe
	RCS_IOSafe
)

type RuntimeContextDef struct {
	CpuLimit    uint64
	MemLimit    uint64
	Flags       RuntimeContextFlags
	SafetyFlags RuntimeSafetyFlags
}
