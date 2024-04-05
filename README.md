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
cookie names. Introducing a more generically named `httpguts.IsToken` function
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
    go test -run ^$ -bench . -benchmem -benchtime 3s -count 20 > new.txt
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
IsCookieNameValid-8      1655.5n ± 0%   367.2n ± 1%  -77.82% (p=0.000 n=20)
ValidHeaderFieldName-8    673.4n ± 1%   330.4n ± 0%  -50.93% (p=0.000 n=20)
IsTokenRune-8             595.2n ± 1%   594.9n ± 1%        ~ (p=0.805 n=20)
geomean                   872.2n        416.4n       -52.26%

                       │     std      │               jub0bs                │
                       │     B/op     │    B/op     vs base                 │
IsCookieNameValid-8      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=20) ¹
ValidHeaderFieldName-8   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=20) ¹
IsTokenRune-8            0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=20) ¹
geomean                             ²               +0.00%                ²
¹ all samples are equal
² summaries must be >0 to compute geomean

                       │     std      │               jub0bs                │
                       │  allocs/op   │ allocs/op   vs base                 │
IsCookieNameValid-8      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=20) ¹
ValidHeaderFieldName-8   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=20) ¹
IsTokenRune-8            0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=20) ¹
geomean                             ²               +0.00%                ²
¹ all samples are equal
² summaries must be >0 to compute geomean
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
