package languages

import "path"

type Language struct {
	name, icon, mext string
	exts             []string
}

func (l *Language) Name() string {
	return l.name
}

func (l *Language) MExt() string {
	return l.mext
}

func (l *Language) Icon() string {
	return l.icon
}

func (l *Language) hasExt(ext string) bool {
	for _, e := range l.exts {
		if e == ext {
			return true
		}
	}
	return false
}

var languages = []*Language{
	{
		name: "C#",
		icon: "/resources/languages/dotnet_logo.svg",
		mext: "cs",
		exts: []string{"cs", "csx"},
	},
	{
		name: "C",
		icon: "/resources/languages/c_logo.svg",
		mext: "c",
		exts: []string{"c", "h"},
	},
	{
		name: "C++",
		icon: "/resources/languages/cpp_logo.svg",
		mext: "cpp",
		exts: []string{"cpp", "C", "cc", "cxx", "c++", "cppm", "ixx", "hpp"},
	},
	{
		name: "Go",
		icon: "/resources/languages/go_logo_blue.svg",
		mext: "go",
		exts: []string{"go"},
	},
	{
		name: "Java",
		icon: "/resources/languages/java_duke_icon.svg",
		mext: "java",
		exts: []string{"java"},
	},
	{
		name: "JavaScript",
		icon: "/resources/languages/js_logo.png",
		mext: "js",
		exts: []string{"js", "mjs", "cjs", "jsx"},
	},
	{
		name: "PHP",
		icon: "/resources/languages/php_logo.svg",
		mext: "php",
		exts: []string{"php", "phtml"},
	},
	{
		name: "Rust",
		icon: "/resources/languages/rust_logo.svg",
		mext: "rs",
		exts: []string{"rs"},
	},
	{
		name: "Scala",
		icon: "/resources/languages/marv_scala_icon.png",
		mext: "scala",
		exts: []string{"scala", "sc"},
	},
	{
		name: "TypeScript",
		icon: "/resources/languages/ts_log_128.svg",
		mext: "ts",
		exts: []string{"ts", "tsx"},
	},
	{
		name: "Python",
		icon: "/resources/languages/python_logo.svg",
		mext: "py",
		exts: []string{"py", "pyi", "pyw"},
	},
	{
		name: "Ruby",
		icon: "/resources/languages/ruby_logo.png",
		mext: "rb",
		exts: []string{"rb"},
	},
	{
		name: "Kotlin",
		icon: "/resources/languages/Kotlin_Kodee.svg",
		mext: "kt",
		exts: []string{"kt", "kts"},
	},
	{
		name: "Haskell",
		icon: "/resources/languages/haskell_logo.svg",
		mext: "hs",
		exts: []string{"hs"},
	},
	{
		name: "SQL",
		icon: "/resources/languages/bad_sql_db_logo.svg",
		mext: "sql",
		exts: []string{"sql"},
	},
}

func UnknownLanguage(ext string) *Language {
	return &Language{
		name: "Unknown",
		icon: "/resources/languages/unknown.png",
		mext: ext,
	}
}

func GetLanguageFromFile(filePath string) *Language {
	ext := path.Ext(filePath)[1:]
	for _, language := range languages {
		if language.hasExt(ext) {
			return language
		}
	}
	return UnknownLanguage(ext)
}
