# Experiment aiming at speeding up golang.org/x/net/http/httpguts

## What is this?

This repo contains alternative implementations of
[`httpguts.ValidHeaderFieldName`][ValidHeaderFieldName] and
[`httpguts.IsTokenRune`][IsTokenRune], as well as 
benchmarks that pit those alternative implementations against the current ones.

## What's the point?

Go's [net/http package][net-http] relies on `httpguts.ValidHeaderFieldName`
extensively. However, `httpguts.ValidHeaderFieldName` is currently implemented
in terms of `httpguts.IsTokenRune`, which systematically incurs a
[bounds check][bc]. Therefore, `httpguts.ValidHeaderFieldName` is not as fast
as it ideally could be.

The alternative implementations not only eliminate those bounds checks,
but also obviate the need for decoding strings into runes.
As a result, they promise a significant speedup.

### Validating methods and cookie names in addition to header-field names

The http/net package also ([indirectly][isNotToken]) relies on
`httpguts.IsTokenRune` to validate [HTTP methods][method-val] and
[cookie names][cookie-val]. However,
- according to [RFC 9110][rfc-9110], [header-field names][field-names] and
  [methods][methods] share the same production: [_token_][token];
- according to [RFC 6265][rfc-6265], [cookie names][cookies] too are _tokens_.

Therefore, net/http could use the faster implementation of
`ValidHeaderFieldName` to validate all three: header names, methods, and
cookies. Introducing a more generically named `httpguts.IsToken` function
(and implementing `httpguts.ValidHeaderFieldName` in terms of it) would
likely make more sense, though.

## Running the benchmarks

1. Make sure that [benchstat][benchstat] is installed:
    ```shell
    go install golang.org/x/perf/cmd/benchstat@latest
    ```
2. Clone this repo and `cd` into it:
    ```shell
    git clone https://github.com/jub0bs/httpguts-perf-exp
    cd httpguts-perf-exp
    ```
3. Run the following commands (preferably on an idle machine):
    ```shell
    go test -run ^$ -bench . -benchtime 3s -count 20 > new.txt
    benchstat -col "/v@(std jub0bs)" new.txt          
    ```

## Some results

```txt
goos: darwin
goarch: amd64
pkg: github.com/jub0bs/httpguts-perf-exp
cpu: Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
                       │     std      │               jub0bs                │
                       │    sec/op    │   sec/op     vs base                │
IsCookieNameValid-8      1655.5n ± 0%   371.9n ± 2%  -77.54% (p=0.000 n=20)
ValidHeaderFieldName-8    672.1n ± 1%   330.9n ± 1%  -50.77% (p=0.000 n=20)
IsTokenRune-8             594.0n ± 1%   594.8n ± 0%        ~ (p=0.836 n=20)
geomean                   871.1n        418.3n       -51.98%
```

## TODO

- Write fuzz tests for `ValidHeaderFieldName`.

[IsTokenRune]: https://pkg.go.dev/golang.org/x/net/http/httpguts#IsTokenRune
[ValidHeaderFieldName]: https://pkg.go.dev/golang.org/x/net/http/httpguts#ValidHeaderFieldName
[bc]: https://en.wikipedia.org/wiki/Bounds_checking
[benchstat]: https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
[cookie-val]: https://github.com/golang/go/blob/2e064cf14441460290fd25d9d61f02a9d0bae671/src/net/http/cookie.go#L463
[cookies]: https://www.rfc-editor.org/rfc/rfc6265.html#section-4.1.1
[field-names]: https://httpwg.org/specs/rfc9110.html#fields.names
[isNotToken]: https://github.com/golang/go/blob/2e064cf14441460290fd25d9d61f02a9d0bae671/src/net/http/http.go#L61
[method-val]: https://github.com/golang/go/blob/2e064cf14441460290fd25d9d61f02a9d0bae671/src/net/http/request.go#L846
[methods]: https://httpwg.org/specs/rfc9110.html#method.overview
[net-http]: https://pkg.go.dev/net/http
[rfc-6265]: https://www.rfc-editor.org/rfc/rfc6265.html
[rfc-9110]: https://httpwg.org/specs/rfc9110.html
[token]: https://www.rfc-editor.org/rfc/rfc2616#section-2.2
