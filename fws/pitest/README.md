### Pitest

Pitest must be run with the `-Dfeatures="+EXPORT"` flag which exports the mutated class files. This is required because
Marv will decompile these class files to construct each mutants replacement string.

> [!NOTE]
> The replacement strings (inserted lines) that Marv produces are correct, however they are occasionally flanked by
> incorrectly formatted deleted lines due to formatting differences between the source code and decompiled class code.

A new Marv Pitest configuration can be created by running the `marv init -f Pitest` command.

#### Decompilers

Marv has a range of decompiler options that can be used with to construct the Pitest mutant replacement strings. They
are listed below.

> [!CAUTION]
> The `garlic` decompiler is currently unstable and using it could cause some mutants to be skipped due to a
> segmentation fault that occurs when running `garlic` on some class files.

* [vineflower-server](https://github.com/SecretSheppy/vineflower-server) (recommended)
* [vineflower](https://github.com/Vineflower/vineflower)
* [garlic](https://github.com/neocanable/garlic)

For installation location see [Installation - Libraries](#libraries)
