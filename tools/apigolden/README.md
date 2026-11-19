# apigolden — the public-API golden of gpucontext

`apigolden` prints every exported declaration of `github.com/gogpu/gpucontext` — types, fields, interface
method sets, functions, methods, constants and variables, with full signatures — as deterministic text.
The committed file `testdata/api.txt` is the frozen public shape of the module; a change to it is a change
to a contract every consumer depends on.

```bash
GOWORK=off go -C tools run ./apigolden -out apigolden/testdata/api.txt
git diff --exit-code tools/apigolden/testdata/   # must be empty
```

## Why one context

gpucontext ships no build-tagged files: the module is one package whose file set is identical on every
platform. `go list` returns the same 20 files for `windows/amd64`, `linux/amd64`, `darwin/arm64` and
`js/wasm`, so one golden describes the API everywhere. The tool still loads through a declared build
context and refuses a context without a floor, so a build tag added later cannot silently shrink the
golden: the context list in `main.go` has to be extended deliberately.

## The floor

`contextFloors` pins the package and identifier counts the loader must reach before the golden is
written. A load that loses the package or drops identifiers fails instead of producing a short file that
looks reviewed. The current floor is `1 package / 346 identifiers`.

## Regenerating after an intended change

1. Edit the code.
2. Regenerate into `testdata/api.txt`.
3. Review the diff as a public-API change: removing or narrowing an exported name is a breaking change
   and needs the ecosystem API-CHANGE round before it ships.
4. Commit the code and the regenerated golden together.

The formatter is covered by unit tests (`go -C tools test ./apigolden`) that build packages in memory, so
they cannot be affected by the tree under analysis.
