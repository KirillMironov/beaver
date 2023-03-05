package mock

type Logger struct{}

func (Logger) Info(...any) {}

func (Logger) Infof(string, ...any) {}

func (Logger) Error(...any) {}

func (Logger) Errorf(string, ...any) {}
