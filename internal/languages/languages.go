package languages

var (
	CSharp     = &Language{name: "C#", ext: "cs", icon: "/resources/languages/dotnet_logo.svg"}
	Cpp        = &Language{name: "C++", ext: "cpp", icon: "/resources/languages/cpp_logo.svg"}
	Go         = &Language{name: "Go", ext: "go", icon: "/resources/languages/go_logo_blue.svg"}
	Java       = &Language{name: "Java", ext: "java", icon: "/resources/languages/java_duke_icon.svg"}
	JavaScript = &Language{name: "JavaScript", ext: "js", icon: "/resources/languages/js_logo.png"}
	Php        = &Language{name: "PHP", ext: "php", icon: "/resources/languages/php_logo.svg"}
	Rust       = &Language{name: "Rust", ext: "rs", icon: "/resources/languages/rust_logo.svg"}
	Scala      = &Language{name: "Scala", ext: "scala", icon: "/resources/languages/marv_scala_icon.png"}
	TypeScript = &Language{name: "TypeScript", ext: "ts", icon: "/resources/languages/ts_log_128.svg"}
)

type Language struct {
	name, ext, icon string
}

func (l *Language) Name() string {
	return l.name
}

func (l *Language) Ext() string {
	return l.ext
}

func (l *Language) Icon() string {
	return l.icon
}
