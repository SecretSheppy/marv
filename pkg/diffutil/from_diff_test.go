package diffutil

import "testing"

func TestDiffLinesExtraction(t *testing.T) {
	diff := `--- mutation diff ---
--- asrc/requests/__init__.py
+++ bsrc/requests/__init__.py
@@ -115,7 +115,7 @@
         chardet_version,
         charset_normalizer_version,
     )
-except (AssertionError, ValueError):
+except (AssertionError, CosmicRayTestingException):
     warnings.warn(
         f"urllib3 ({urllib3.__version__}) or chardet "
         f"({chardet_version})/charset_normalizer ({charset_normalizer_version}) "`
	fdiff := FromFormattedDiff(diff, &DiffConfig{
		PrefixLines:            4,
		SuffixLines:            0,
		FirstRemovedLineNumber: 3,
	})
	if err := fdiff.Number(); err != nil {
		t.Error(err)
	}
	lines := fdiff.Lines()
	if len(lines) != 8 {
		t.Errorf("expected 8 lines but got %d", len(lines))
	}
	if lines[0].Number != 0 {
		t.Errorf("expected first line number to be 0 but got %d", lines[0].Number)
	}
	if lines[3].Type != Removed {
		t.Errorf("expected line 3 to be REMOVED but got %v", lines[3].Type)
	}
	if lines[4].Type != Inserted {
		t.Errorf("expected line 4 to be INSERTED but got %v", lines[4].Type)
	}
}

func TestDiffLinesPadding(t *testing.T) {
	diff := `--- mutation diff ---
--- asrc/requests/__init__.py
+++ bsrc/requests/__init__.py
@@ -73,7 +73,7 @@
     major, minor, patch = urllib3_version_list  # noqa: F811
     major, minor, patch = int(major), int(minor), int(patch)
     # urllib3 >= 1.21.1
-    assert major >= 1
+    assert major != 1
     if major == 1:
         assert minor >= 21`
	fdiff := FromFormattedDiff(diff, &DiffConfig{
		PrefixLines:            4,
		SuffixLines:            0,
		FirstRemovedLineNumber: 3,
	})
	if err := fdiff.Number(); err != nil {
		t.Error(err)
	}
	// NOTE: we want to sync an extra 4 spaces into the diff lines
	fdiff.SyncLineFormatting([]string{
		"        major, minor, patch = urllib3_version_list  # noqa: F811",
		"        major, minor, patch = int(major), int(minor), int(patch)",
		"        # urllib3 >= 1.21.1",
		"        assert major >= 1",
		"        if major == 1:",
		"            assert minor >= 21",
	})
	lines := fdiff.Lines()
	if len(lines) != 7 {
		t.Errorf("expected 7 lines but got %d", len(lines))
	}
	if lines[3].Text != "        assert major >= 1" {
		t.Errorf("expected line:\n\n        assert major >= 1\n\ngot line:\n\n%s\n\n", lines[3].Text)
	}
	if lines[4].Text != "        assert major != 1" {
		t.Errorf("expected line:\n\n        assert major != 1\n\ngot line:\n\n%s\n\n", lines[3].Text)
	}
}

func TestExtractingBlankLines(t *testing.T) {
	diff := `--- mutation diff ---
--- asrc/requests/__init__.py
+++ bsrc/requests/__init__.py
@@ -216,5 +216,5 @@
 logging.getLogger(__name__).addHandler(NullHandler())
 
 # FileModeWarnings go off per the default.
-warnings.simplefilter("default", FileModeWarning, append=True)
-
+warnings.simplefilter("default", FileModeWarning, append=False)
+`
	fdiff := FromFormattedDiff(diff, &DiffConfig{
		PrefixLines: 4,
	})
	lines := fdiff.Lines()
	if len(lines) != 7 {
		t.Errorf("expected 7 lines but got %d", len(lines))
	}
	if lines[4].Type != Removed {
		t.Errorf("expected line 4 to be REMOVED but got %v", lines[4].Type)
	}
	if lines[4].Text != "" {
		t.Errorf("expected line 4 to be \"\" but got:\n\n%s\n\n", lines[4].Text)
	}
}
