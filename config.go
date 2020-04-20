package pagser

// Config configuration
type Config struct {
	TagerName    string //struct tag name, default is `pagser`
	FuncSymbol   string //Function symbol, default is `->`
	IgnoreSymbol string //Ignore symbol, default is `-`
	Debug        bool   //Debug mode, debug will print some log, default is `false`
}

var defaultCfg = Config{
	TagerName:    "pagser",
	FuncSymbol:   "->",
	IgnoreSymbol: "-",
	Debug:        false,
}
