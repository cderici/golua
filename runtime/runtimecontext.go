package runtime

type RuntimeContextStatus uint8

const (
	RCS_Live RuntimeContextStatus = iota
	RCS_Done
	RCS_Killed
)

type RuntimeContext interface {
	CpuLimit() uint64
	CpuUsed() uint64

	MemLimit() uint64
	MemUsed() uint64

	Status() RuntimeContextStatus
	Parent() RuntimeContext

	Flags() RuntimeContextFlags
}

type RuntimeContextFlags uint8

const (
	RCF_Empty RuntimeContextFlags = 1 << iota
	RCF_NoIO
)

func (f RuntimeContextFlags) IsSet(ctx RuntimeContext) bool {
	return f&ctx.Flags() != 0
}

var ErrIODisabled = NewErrorS("io disabled")

func (r *Runtime) CheckIO() *Error {
	if RCF_NoIO.IsSet(r) {
		return ErrIODisabled
	}
	return nil
}