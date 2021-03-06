package bsondump

type BSONDumpOptions struct {
	Type     string `long:"type" default:"json" default-mask:"-" description:"type of output: debug, json (default 'json')"`
	ObjCheck bool   `long:"objcheck" description:"validate BSON during processing"`
	Pretty   bool   `long:"pretty" description:"output JSON formatted to be human-readable"`
}

func (_ *BSONDumpOptions) Name() string {
	return "output"
}

func (_ *BSONDumpOptions) PostParse() error {
	return nil
}

func (_ *BSONDumpOptions) Validate() error {
	return nil
}
