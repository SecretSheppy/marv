# Cosmic Ray

## Contents

## Getting Started With Cosmic Ray In Marv

## Rare Formatting Errors

```diff
--- mutation diff ---
--- asrc/requests/__init__.py
+++ bsrc/requests/__init__.py
@@ -216,5 +216,5 @@
 logging.getLogger(__name__).addHandler(NullHandler())
 
 # FileModeWarnings go off per the default.
-warnings.simplefilter("default", FileModeWarning, append=True)
-
+warnings.simplefilter("default", FileModeWarning, append=False)
+
```

![formatting error in Marv interface](docs/formatting_error.png)