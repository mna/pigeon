# pigeon performance

## VM implementation notes

* Initializing the stacks capacities to 128 elements seems to help a little bit, but there is no noticeable improvement by using 512 or 1024.
* Removing the bounds checks in the stacks don't translate to any noticeable improvement.

### Commit info

* 6be069a : change parser.peek to return a rune instead of the svpt struct.
* cffe69a : use decoded instruction (struct instead of encoded uint64).

// Commit df3f721 (recursive)
BenchmarkParseUnicodeClass         10000            169352 ns/op           15013 B/op        234 allocs/op
BenchmarkParseKeyword              10000            139439 ns/op           14070 B/op        202 allocs/op

// Commit 0850249 (vm)
BenchmarkParseUnicodeClass          2000           1657522 ns/op          361434 B/op       3016 allocs/op
BenchmarkParseKeyword               1000           1972652 ns/op          207480 B/op       1870 allocs/op

// Commit 1317e07 (vm+stacks 128)
BenchmarkParseUnicodeClass          2000           2845717 ns/op          367405 B/op       2990 allocs/op
BenchmarkParseKeyword               2000            649629 ns/op          213968 B/op       1845 allocs/op

// Commit 0b32ca6 (stack rewritten with sp)
BenchmarkParseUnicodeClass          2000            826942 ns/op          338877 B/op       2399 allocs/op
BenchmarkParseKeyword               3000            549708 ns/op          195856 B/op       1471 allocs/op

// Commit 6be069a
BenchmarkParseUnicodeClass          2000            867650 ns/op          338877 B/op       2399 allocs/op
BenchmarkParseKeyword               3000            579195 ns/op          195856 B/op       1471 allocs/op

// Commit cffe69a
BenchmarkParseUnicodeClass          2000            670882 ns/op          105810 B/op       1949 allocs/op
BenchmarkParseKeyword               3000            432649 ns/op           68926 B/op       1204 allocs/op



// Commit df3f721 (recursive)
BenchmarkParsePigeonNoMemo            30          38052374 ns/op         3212049 B/op      71466 allocs/op
BenchmarkParsePigeonMemo              20          82778941 ns/op        30789484 B/op      66046 allocs/op

// Commit 0850249 (vm - memo is nop)
BenchmarkParsePigeonNoMemo            10         143073240 ns/op        87630939 B/op     693740 allocs/op
BenchmarkParsePigeonMemo              10         142642384 ns/op        87630939 B/op     693740 allocs/op

// Commit 1317e07 (vm+stacks 128, memo is nop)
BenchmarkParsePigeonNoMemo            10         141237021 ns/op        87629233 B/op     693708 allocs/op
BenchmarkParsePigeonMemo              10         141600274 ns/op        87629233 B/op     693708 allocs/op

// Commit 0b32ca6 (stack rewritten with sp)
BenchmarkParsePigeonNoMemo            10         116394376 ns/op        80393633 B/op     542970 allocs/op
BenchmarkParsePigeonMemo              10         116237323 ns/op        80393633 B/op     542970 allocs/op

// Commit 6be069a
BenchmarkParsePigeonNoMemo            10         116662711 ns/op        80393633 B/op     542970 allocs/op
BenchmarkParsePigeonMemo              10         117390138 ns/op        80393633 B/op     542970 allocs/op

// Commit cffe69a
BenchmarkParsePigeonNoMemo            20          89610719 ns/op        25281650 B/op     458273 allocs/op
BenchmarkParsePigeonMemo              20          88681471 ns/op        25281651 B/op     458273 allocs/op


// Commit df3f721 (recursive)
BenchmarkPigeonJSONNoMemo             50          25212587 ns/op         3328296 B/op      86105 allocs/op
BenchmarkPigeonJSONMemo               20          86689562 ns/op        25050390 B/op     131153 allocs/op

// Commit 0850249 (vm - memo is nop)
BenchmarkPigeonJSONNoMemo             20          93411006 ns/op        56357596 B/op     492396 allocs/op
BenchmarkPigeonJSONMemo               20          93271080 ns/op        56357640 B/op     492396 allocs/op

// Commit 1317e07 (vm+stacks 128, memo is nop)
BenchmarkPigeonJSONNoMemo             20          93790644 ns/op        56363398 B/op     492369 allocs/op
BenchmarkPigeonJSONMemo               20          93547428 ns/op        56363254 B/op     492369 allocs/op

// Commit 0b32ca6 (stack rewritten with sp)
BenchmarkPigeonJSONNoMemo             20          75739611 ns/op        52532741 B/op     412572 allocs/op
BenchmarkPigeonJSONMemo               20          75758382 ns/op        52533317 B/op     412574 allocs/op

// Commit 6be069a
BenchmarkPigeonJSONNoMemo             20          76681525 ns/op        52532914 B/op     412573 allocs/op
BenchmarkPigeonJSONMemo               20          77146195 ns/op        52533188 B/op     412574 allocs/op

// Commit cffe69a
BenchmarkPigeonJSONNoMemo             20          60564462 ns/op        18194950 B/op     371580 allocs/op
BenchmarkPigeonJSONMemo               20          61545542 ns/op        18195112 B/op     371580 allocs/op


// Commit df3f721 (recursive)
BenchmarkPigeonCalculatorNoMemo    10000            169574 ns/op           17958 B/op        390 allocs/op
BenchmarkPigeonCalculatorMemo       2000            672838 ns/op          132173 B/op        515 allocs/op

// Commit 0b32ca6 (stack rewritten with sp)
BenchmarkPigeonCalculatorNoMemo     3000            494580 ns/op          167408 B/op       1329 allocs/op
BenchmarkPigeonCalculatorMemo       3000            498565 ns/op          167408 B/op       1329 allocs/op

// Commit 6be069a
BenchmarkPigeonCalculatorNoMemo     3000            423499 ns/op          155104 B/op       1329 allocs/op
BenchmarkPigeonCalculatorMemo       3000            426319 ns/op          155104 B/op       1329 allocs/op

// Commit cffe69a
BenchmarkPigeonCalculatorNoMemo     5000            389739 ns/op           64723 B/op       1224 allocs/op
BenchmarkPigeonCalculatorMemo       5000            389930 ns/op           64723 B/op       1224 allocs/op


// Go1.4 stdlib
BenchmarkStdlibJSON                 2000            861586 ns/op           74094 B/op       1055 allocs/op
