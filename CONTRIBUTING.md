# Contributing to Marv

Marv was designed as a tool that could be easily extended by users, and as such new contributions to Marv are welcome.
Please read the rest of this documented to best understand how one can contribute, and more importantly, what not to
try contributing.

## What Will Not Be Accepted

The below is a non-exhaustive list of contributions that will be rejected without review.

* AI Generated PRs that touch irrelevant files and make useless changes
* Large changes to the user interface
* Changes that include copyrighted or trademarked material

Changes to the underlying architecture or the Marv mutations schema will **always** be rejected without review unless
they have been thoroughly discussed in an issue beforehand.

## Frameworks

New framework implementation are always welcome. To add a framework to Marv one must implement the `fwlib.Framework`
interface. All framework implementations should be in a folder named after the framework, using underscores for spaces
and committed in the `fws` directory. The `fwlib.Framework` implementation should also be added to the `fws/fws.go`
file in the `Frameworks` function. The `fws` directory contains many examples of already implemented frameworks which
should act as a good guideline for a new implementation. It is preferable that structures like `YamlConfig` and 
`YamlWrapper` retain their standardized names in any new framework implementation, although this can be changed after
an implementation pull request has been created.

## Themes

New themes are fun and easy contributions. A theme can be created by copying and renaming one of the existing JSON
themes files found in the [web/themes](web/themes) directory. Marv's theming system determines whether to use dark
or light branding based on the color of the main background, and the interface icons are set to the main text foreground
color.