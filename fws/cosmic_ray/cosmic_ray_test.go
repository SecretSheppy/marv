package cosmic_ray

import "testing"

func TestDiffStringChangeExtraction(t *testing.T) {
	tests := []struct {
		Name, Diff, ExpRem, ExpIns string
	}{
		{
			Name: "Test extracting changes with no padding after - and +",
			Diff: `--- mutation diff ---
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
         f"({chardet_version})/charset_normalizer ({charset_normalizer_version}) "`,
			ExpRem: "except (AssertionError, ValueError):",
			ExpIns: "except (AssertionError, CosmicRayTestingException):",
		},
		{
			Name: "Test extracting changes with padding after - and +",
			Diff: `--- mutation diff ---
--- asrc/requests/__init__.py
+++ bsrc/requests/__init__.py
@@ -73,7 +73,7 @@
     major, minor, patch = urllib3_version_list  # noqa: F811
     major, minor, patch = int(major), int(minor), int(patch)
     # urllib3 >= 1.21.1
-    assert major >= 1
+    assert major != 1
     if major == 1:
         assert minor >= 21`,
			ExpRem: "assert major >= 1",
			ExpIns: "assert major != 1",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			rem, ins := changes(test.Diff)
			if rem != test.ExpRem {
				t.Errorf("expected removed line >> %s\ngot >> %s", rem, test.ExpRem)
			}
			if ins != test.ExpIns {
				t.Errorf("expected inserted line >> %s\ngot >> %s", rem, test.ExpRem)
			}
		})
	}
}

// --- mutation diff ---
//--- asrc/requests/__init__.py
//+++ bsrc/requests/__init__.py
//@@ -216,5 +216,5 @@
// logging.getLogger(__name__).addHandler(NullHandler())
//
// # FileModeWarnings go off per the default.
//-warnings.simplefilter("default", FileModeWarning, append=True)
//-
//+warnings.simplefilter("default", FileModeWarning, append=False)
//+
