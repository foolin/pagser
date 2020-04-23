package pagser

const ignoreSymbol = "-"

// Config configuration
type Config struct {
	TagerName  string //struct tag name, default is `pagser`
	FuncSymbol string //Function symbol, default is `->`
	CastError  bool   //Returns an error when the type cannot be converted, default is `false`
	Debug      bool   //Debug mode, debug will print some log, default is `false`
}

var defaultCfg = Config{
	TagerName:  "pagser",
	FuncSymbol: "->",
	CastError:  false,
	Debug:      false,
}
