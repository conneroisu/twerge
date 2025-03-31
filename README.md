# twerge

generates tailwind merges and classes from go templ sources

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# tmplmerge

```go
import "github.com/conneroisu/tmplmerge"
```

package twerge provides a tailwind merger for go\-templ and code generation.

## Index

- [Constants](#constants)
- [Variables](#variables)
- [func GetIsArbitraryValue\(val string, label any, testValue func\(string\) bool\) bool](#GetIsArbitraryValue)
- [func IsArbitraryNumber\(val string\) bool](#IsArbitraryNumber)
- [func MakeMergeClassList\(conf \*Config, splitModifiers SplitModifiersFn, getClassGroupID GetClassGroupIdfn\) func\(classList string\) string](#MakeMergeClassList)
- [func SortModifiers\(modifiers \[\]string\) \[\]string](#SortModifiers)
- [type ClassGroupValidator](#ClassGroupValidator)
- [type ClassPart](#ClassPart)
- [type Config](#Config)
  - [func MakeDefaultConfig\(\) \*Config](#MakeDefaultConfig)
- [type ConflictingClassGroups](#ConflictingClassGroups)
- [type GetClassGroupIdfn](#GetClassGroupIdfn)
  - [func MakeGetClassGroupID\(conf \*Config\) GetClassGroupIdfn](#MakeGetClassGroupID)
- [type ICache](#ICache)
  - [func Make\(maxCapacity int\) ICache](#Make)
- [type SplitModifiersFn](#SplitModifiersFn)
  - [func MakeSplitModifiers\(conf \*Config\) SplitModifiersFn](#MakeSplitModifiers)
- [type TwMergeFn](#TwMergeFn)
  - [func CreateTwMerge\(config \*Config, cache ICache\) TwMergeFn](#CreateTwMerge)

## Constants

<a name="SplitClassesRegex"></a>SplitClassesRegex is the regex used to split classes

```go
const SplitClassesRegex = `\s+`
```

## Variables

<a name="Merge"></a>Merge is the default template merger

```go
var Merge = CreateTwMerge(nil, nil)
```

<a name="GetIsArbitraryValue"></a>

## func [GetIsArbitraryValue](https://github.com/conneroisu/tmplmerge/blob/main/config.go#L200-L204)

```go
func GetIsArbitraryValue(val string, label any, testValue func(string) bool) bool
```

GetIsArbitraryValue returns true if the given value is an arbitrary value with the given label. The label can be a string, a map\[string\]bool or a function that takes a string and returns a bool.

<a name="IsArbitraryNumber"></a>

## func [IsArbitraryNumber](https://github.com/conneroisu/tmplmerge/blob/main/config.go#L122)

```go
func IsArbitraryNumber(val string) bool
```

IsArbitraryNumber returns true if the given value is an arbitrary number

<a name="MakeMergeClassList"></a>

## func [MakeMergeClassList](https://github.com/conneroisu/tmplmerge/blob/main/merge.go#L74-L78)

```go
func MakeMergeClassList(conf *Config, splitModifiers SplitModifiersFn, getClassGroupID GetClassGroupIdfn) func(classList string) string
```

MakeMergeClassList creates a function that merges a class list

<a name="SortModifiers"></a>

## func [SortModifiers](https://github.com/conneroisu/tmplmerge/blob/main/merge.go#L127)

```go
func SortModifiers(modifiers []string) []string
```

SortModifiers Sorts modifiers according to following schema: \- Predefined modifiers are sorted alphabetically \- When an arbitrary variant appears, it must be preserved which modifiers are before and after it

<a name="ClassGroupValidator"></a>

## type [ClassGroupValidator](https://github.com/conneroisu/tmplmerge/blob/main/config.go#L45-L48)

ClassGroupValidator is a validator for a class group

```go
type ClassGroupValidator struct {
    Fn           func(string) bool
    ClassGroupID string
}
```

<a name="ClassPart"></a>

## type [ClassPart](https://github.com/conneroisu/tmplmerge/blob/main/config.go#L51-L55)

ClassPart is a part of a class group

```go
type ClassPart struct {
    NextPart     map[string]ClassPart
    Validators   []ClassGroupValidator
    ClassGroupID string
}
```

<a name="Config"></a>

## type [Config](https://github.com/conneroisu/tmplmerge/blob/main/config.go#L18-L42)

Config is the configuration for the template merger

```go
type Config struct {
    // defaults should be good enough
    // hover:bg-red-500 -> :
    ModifierSeparator rune
    // bg-red-500 -> -
    ClassSeparator rune
    // !bg-red-500 -> !
    ImportantModifier rune
    // used for bg-red-500/50 (50% opacity) -> /
    PostfixModifier rune
    // optional
    Prefix string

    // CACHE
    MaxCacheSize int

    // This is a large map of all the classes and their validators -> see default-config.go
    ClassGroups ClassPart

    // class group with conflict + conflicting groups -> if "p" is set all others are removed
    // p: ['px', 'py', 'ps', 'pe', 'pt', 'pr', 'pb', 'pl']
    ConflictingClassGroups ConflictingClassGroups
}
```

<a name="MakeDefaultConfig"></a>

### func [MakeDefaultConfig](https://github.com/conneroisu/tmplmerge/blob/main/config.go#L229)

```go
func MakeDefaultConfig() *Config
```

MakeDefaultConfig returns a default TwMergeConfig

<a name="ConflictingClassGroups"></a>

## type [ConflictingClassGroups](https://github.com/conneroisu/tmplmerge/blob/main/config.go#L58)

ConflictingClassGroups is a map of class groups that conflict with each other

```go
type ConflictingClassGroups map[string][]string
```

<a name="GetClassGroupIdfn"></a>

## type [GetClassGroupIdfn](https://github.com/conneroisu/tmplmerge/blob/main/class.go#L9)

GetClassGroupIdfn returns the class group id for a given class

```go
type GetClassGroupIdfn func(string) (isTwClass bool, groupId string)
```

<a name="MakeGetClassGroupID"></a>

### func [MakeGetClassGroupID](https://github.com/conneroisu/tmplmerge/blob/main/class.go#L12)

```go
func MakeGetClassGroupID(conf *Config) GetClassGroupIdfn
```

MakeGetClassGroupID returns a GetClassGroupIdfn

<a name="ICache"></a>

## type [ICache](https://github.com/conneroisu/tmplmerge/blob/main/lru.go#L23-L26)

ICache is the interface for a LRU cache

```go
type ICache interface {
    Get(string) string
    Set(string, string)
}
```

<a name="Make"></a>

### func [Make](https://github.com/conneroisu/tmplmerge/blob/main/lru.go#L8)

```go
func Make(maxCapacity int) ICache
```

Make creates a new LRU cache

<a name="SplitModifiersFn"></a>

## type [SplitModifiersFn](https://github.com/conneroisu/tmplmerge/blob/main/merge.go#L21)

SplitModifiersFn is the type of the function used to split modifiers

```go
type SplitModifiersFn = func(string) (baseClass string, modifiers []string, hasImportant bool, maybePostfixModPosition int)
```

<a name="MakeSplitModifiers"></a>

### func [MakeSplitModifiers](https://github.com/conneroisu/tmplmerge/blob/main/merge.go#L154)

```go
func MakeSplitModifiers(conf *Config) SplitModifiersFn
```

MakeSplitModifiers creates a function that splits modifiers

<a name="TwMergeFn"></a>

## type [TwMergeFn](https://github.com/conneroisu/tmplmerge/blob/main/merge.go#L18)

TwMergeFn is the type of the template merger.

```go
type TwMergeFn func(args ...string) string
```

<a name="CreateTwMerge"></a>

### func [CreateTwMerge](https://github.com/conneroisu/tmplmerge/blob/main/merge.go#L24-L27)

```go
func CreateTwMerge(config *Config, cache ICache) TwMergeFn
```

CreateTwMerge creates a new template merger

Generated by [gomarkdoc](https://github.com/princjef/gomarkdoc)

<!-- gomarkdoc:embed:end -->
