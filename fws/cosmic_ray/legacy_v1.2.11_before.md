# Cosmic Ray Legacy Documentation

> [!IMPORTANT]
> The Cosmic Ray framework implementation in Marv used to produce visual errors when Marv encountered an edge case in
> Cosmic Ray's output. This was present in all versions of Marv before and including v1.2.11. Extensive edge case handling was added in
> Marv v1.3.0 and all Marv versions v1.3.0+ will not produce visual errors with Cosmic Ray.

## Rare Formatting Errors

As Cosmic Ray only specifies mutations location in the source code and not the actual source code mutation string, Marv
has to extract it from the provided diffs. In most cases this will work flawlessly; however, occasionally Cosmic Ray
will output an imperfect diff, an example of which is shown below.

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

**Caption:** Imperfect diff produced by Cosmic Ray. The second lines marked `-` and `+` are not actually changed,
and Cosmic Ray itself is only expecting one line to be mutated.

These imperfect diffs cause the Marv process that extracts the source code mutations to return slightly incorrect
values. How this is reflected in the Marv interface is shown in the below screenshot.

![formatting error in Marv interface](docs/formatting_error.png)

**Caption:** Incorrect formatting in the Marv interface.

> [!IMPORTANT]
> When formatting errors occur it is still obvious what the mutation is, it is just good to be aware that formatting
errors may occur.
