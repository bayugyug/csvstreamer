package csvstreamer

//CsvStreamOpt option wrapper
type CsvStreamOpt func(*CsvStream)

//WithSeparator default is comma
func WithSeparator(s string) CsvStreamOpt {
	return func(args *CsvStream) {
		args.sep = s
	}
}

//WithEnclosedBy default is the double-qts
func WithEnclosedBy(s string) CsvStreamOpt {
	return func(args *CsvStream) {
		args.enclosed = s
	}
}

//WithEmpty default is the dash -
func WithEmpty(s string) CsvStreamOpt {
	return func(args *CsvStream) {
		args.emptyval = s
	}
}
