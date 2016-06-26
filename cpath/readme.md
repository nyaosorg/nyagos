cpath
=====

Package to operate path-string for some special use.

Join
----

	func Join(paths ...string) string

- Compatible Join with MFC / altpath.h / CPath::Combine
    - (ex:`C:\foo` + `\bar` -> `c:\bar`)
- Do not clean path (keep `./` on arguments)

GetHome
-------
	func cpath.GetHome()

- Get %HOME% or %USERPROFILE%

ReplaceHomeToTilde
------------------
	func ReplaceHomeToTilde(wd string) string

- C:\users\name\foo\bar -> ~\foo\bar

ReplaceHomeToTildeSlash
-----------------------
	func ReplaceHomeToTildeSlash(wd string) string

- C:\users\name\foo\bar -> ~/foo/bar

IsExecutableSuffix
==================
	func IsExecutableSuffix(suffix string) bool

- returns true if suffix exists in %PATHEXT%
