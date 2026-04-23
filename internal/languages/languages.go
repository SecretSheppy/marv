package languages

var (
	Cpp  = &Language{name: "C++", ext: "cpp", icon: "/resources/languages/cpp_logo.svg"}
	Go   = &Language{name: "Go", ext: "go", icon: "/resources/languages/go_logo_blue.svg"}
	Java = &Language{name: "Java", ext: "java", icon: "/resources/languages/java_duke_icon.svg"}
	Rust = &Language{name: "Rust", ext: "rs", icon: "/resources/languages/rust_logo.svg"}
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
