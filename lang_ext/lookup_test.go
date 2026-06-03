package lang_ext

import (
	"strings"
	"sync"
	"testing"
)

func TestLanguageOf(t *testing.T) {
	cases := []struct {
		ext      string
		wantLang string
		wantOk   bool
	}{
		{".go", "Golang", true},
		{".py", "Python", true},
		{".ts", "TypeScript", true},
		{".js", "JavaScript", true},
		{".tsx", "TS React", true},
		{".unknown_ext_xx", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		t.Run(c.ext, func(t *testing.T) {
			got, ok := LanguageOf(c.ext)
			if ok != c.wantOk || got != c.wantLang {
				t.Errorf("LanguageOf(%q) = (%q, %v), want (%q, %v)",
					c.ext, got, ok, c.wantLang, c.wantOk)
			}
		})
	}
}

func TestIsExcludeDir(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"node_modules", true},
		{"vendor", true},
		{"dist", true},
		{"public", true},
		{"src", false},
		{"cmd", false},
		{"", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsExcludeDir(c.name); got != c.want {
				t.Errorf("IsExcludeDir(%q) = %v, want %v", c.name, got, c.want)
			}
		})
	}
}

func TestIsExcludeFileExt(t *testing.T) {
	cases := []struct {
		ext  string
		want bool
	}{
		{".exe", true},
		{".so", true},
		{".pyc", true},
		{".dll", true},
		{".go", false},
		{".py", false},
		{"", false},
	}
	for _, c := range cases {
		t.Run(c.ext, func(t *testing.T) {
			if got := IsExcludeFileExt(c.ext); got != c.want {
				t.Errorf("IsExcludeFileExt(%q) = %v, want %v", c.ext, got, c.want)
			}
		})
	}
}

func TestIsExcludeFileName(t *testing.T) {
	if !IsExcludeFileName("package.json") {
		t.Errorf("package.json should be excluded")
	}
	if !IsExcludeFileName("package-lock.json") {
		t.Errorf("package-lock.json should be excluded")
	}
	if IsExcludeFileName("main.go") {
		t.Errorf("main.go should not be excluded")
	}
}

func TestIsSkipLineCountExt(t *testing.T) {
	for _, ext := range []string{".png", ".jpg", ".mp4", ".zip", ".pdf"} {
		if !IsSkipLineCountExt(ext) {
			t.Errorf("%q should skip line count", ext)
		}
	}
	for _, ext := range []string{".go", ".py", ".md"} {
		if IsSkipLineCountExt(ext) {
			t.Errorf("%q should not skip line count", ext)
		}
	}
}

func TestExtOfLanguage(t *testing.T) {
	ext, ok := ExtOfLanguage("Golang")
	if !ok || ext != ".go" {
		t.Errorf(`ExtOfLanguage("Golang") = (%q, %v), want (".go", true)`, ext, ok)
	}
	if _, ok := ExtOfLanguage("NotALanguage"); ok {
		t.Errorf(`ExtOfLanguage("NotALanguage") should miss`)
	}
}

func TestExtOfLanguageLower(t *testing.T) {
	ext, ok := ExtOfLanguageLower("golang")
	if !ok || ext != ".go" {
		t.Errorf(`ExtOfLanguageLower("golang") = (%q, %v), want (".go", true)`, ext, ok)
	}
	if _, ok := ExtOfLanguageLower("Golang"); ok {
		t.Errorf(`ExtOfLanguageLower("Golang") should miss (key is lowercased)`)
	}
}

func TestKnownFileExtOf(t *testing.T) {
	for _, ext := range []string{".go", ".exe", ".png"} {
		if _, ok := KnownFileExtOf(ext); !ok {
			t.Errorf("KnownFileExtOf(%q) should be true (合并源覆盖语言+排除+二进制)", ext)
		}
	}
	if _, ok := KnownFileExtOf(".unknown_xx"); ok {
		t.Errorf("unknown ext should miss")
	}
}

func TestCommonFileNameOf(t *testing.T) {
	for _, name := range []string{"Makefile", "makefile", "Dockerfile", "README.md"} {
		if _, ok := CommonFileNameOf(name); !ok {
			t.Errorf("CommonFileNameOf(%q) should be true", name)
		}
	}
	if _, ok := CommonFileNameOf("random.txt"); ok {
		t.Errorf("random.txt should not be a common file name")
	}
}

func TestFrontLanguageExts_ReturnsIsolatedCopy(t *testing.T) {
	first := FrontLanguageExts()
	if _, ok := first[".html"]; !ok {
		t.Fatalf("FrontLanguageExts should contain .html")
	}

	first[".html"] = "hacked"
	delete(first, ".vue")

	second := FrontLanguageExts()
	if second[".html"] == "hacked" {
		t.Errorf("copy not isolated: second[.html]=%q", second[".html"])
	}
	if _, ok := second[".vue"]; !ok {
		t.Errorf("copy not isolated: second call missing .vue")
	}
}

func TestBackLanguageExts_ReturnsIsolatedCopy(t *testing.T) {
	first := BackLanguageExts()
	if _, ok := first[".go"]; !ok {
		t.Fatalf("BackLanguageExts should contain .go")
	}

	first[".go"] = "hacked"
	delete(first, ".java")

	second := BackLanguageExts()
	if second[".go"] == "hacked" {
		t.Errorf("copy not isolated: second[.go]=%q", second[".go"])
	}
	if _, ok := second[".java"]; !ok {
		t.Errorf("copy not isolated: second call missing .java")
	}
}

// init 把 commonLanguageExt 反向构建到两个 reverse map,
// round-trip 验证:扩展名 → 语言 → reverse 出的扩展名 → 仍解析到同语言。
// 多个扩展名映射到同一语言时(.cc/.cpp 都是 C++),reverse 只保留最后一次写入,
// 故仅要求"reverse 出的扩展名"能正向回到原语言名。
func TestReverseMapsConsistentWithInit(t *testing.T) {
	for ext, lang := range commonLanguageExt {
		gotExt, ok := ExtOfLanguage(lang)
		if !ok {
			t.Errorf("reverse missing for language %q (ext=%q)", lang, ext)
			continue
		}
		if back, ok := LanguageOf(gotExt); !ok || back != lang {
			t.Errorf("round-trip mismatch: %q -> %q -> (%q, %v)", ext, lang, back, ok)
		}

		lowerExt, ok := ExtOfLanguageLower(strings.ToLower(lang))
		if !ok {
			t.Errorf("lower-reverse missing for language %q", lang)
			continue
		}
		if back, ok := LanguageOf(lowerExt); !ok || back != lang {
			t.Errorf("lower round-trip mismatch: %q -> %q -> (%q, %v)", ext, lang, back, ok)
		}
	}
}

// 并发只读必须无 race(init 后包内 map 仅被读)。
// 配合 `go test -race` 检测。
func TestConcurrentReadNoRace(t *testing.T) {
	const workers = 8
	const iterations = 500

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for range iterations {
				_, _ = LanguageOf(".go")
				_ = IsExcludeDir("node_modules")
				_ = IsExcludeFileExt(".exe")
				_ = IsExcludeFileName("package.json")
				_ = IsSkipLineCountExt(".png")
				_, _ = ExtOfLanguage("Golang")
				_, _ = ExtOfLanguageLower("golang")
				_, _ = KnownFileExtOf(".go")
				_, _ = CommonFileNameOf("Makefile")
				_ = FrontLanguageExts()
				_ = BackLanguageExts()
			}
		}()
	}
	wg.Wait()
}

func TestCommonDepLockFile(t *testing.T) {
	m := CommonDepLockFile()
	for _, key := range []string{".sum", ".lock", ".gitignore"} {
		if _, ok := m[key]; !ok {
			t.Errorf("CommonDepLockFile missing %q", key)
		}
	}
}
