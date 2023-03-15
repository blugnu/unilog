package unilog

type nulAdapter struct{}

func noop() {}

func (*nulAdapter) Emit(Level, string)                { noop() }
func (nul *nulAdapter) NewEntry() Adapter             { return nul }
func (nul *nulAdapter) WithField(string, any) Adapter { return nul }

var nul = &logger{Adapter: &nulAdapter{}}

func Nul() Logger {
	return nul
}
