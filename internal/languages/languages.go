package languages

var (
	Go   = Language{name: "Go", ext: "go", icon: "/resources/languages/golang-brands-solid.svg"}
	Java = Language{name: "Java", ext: "java", icon: "/resources/languages/java-brands-solid.svg"}
	Rust = Language{name: "Rust", ext: "rs", icon: "/resources/languages/rust-brands-solid.svg"}
)

type Language struct {
	name, ext, icon string
}

func (l Language) Name() string {
	return l.name
}

func (l Language) Ext() string {
	return l.ext
}

func (l Language) Icon() string {
	return l.icon
}
