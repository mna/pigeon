# pigeon performance

// Commit df3f721 (recursive)
BenchmarkParseUnicodeClass         10000            169352 ns/op           15013 B/op        234 allocs/op
BenchmarkParseKeyword              10000            139439 ns/op           14070 B/op        202 allocs/op

// Commit 0850249 (vm)
BenchmarkParseUnicodeClass          2000           1657522 ns/op          361434 B/op       3016 allocs/op
BenchmarkParseKeyword               1000           1972652 ns/op          207480 B/op       1870 allocs/op

// Commit 1317e07 (vm+stacks 128)
BenchmarkParseUnicodeClass          2000            991382 ns/op          361434 B/op       3016 allocs/op
BenchmarkParseKeyword               2000            657234 ns/op          207480 B/op       1870 allocs/op


// Commit df3f721 (recursive)
BenchmarkParsePigeonNoMemo            30          38052374 ns/op         3212049 B/op      71466 allocs/op
BenchmarkParsePigeonMemo              20          82778941 ns/op        30789484 B/op      66046 allocs/op

// Commit 0850249 (vm - memo is nop)
BenchmarkParsePigeonNoMemo            10         143073240 ns/op        87630939 B/op     693740 allocs/op
BenchmarkParsePigeonMemo              10         142642384 ns/op        87630939 B/op     693740 allocs/op

// Commit 1317e07 (vm+stacks 128)
BenchmarkParsePigeonNoMemo            10         141508295 ns/op        87630939 B/op     693740 allocs/op
BenchmarkParsePigeonMemo              10         143317953 ns/op        87630939 B/op     693740 allocs/op


// Commit df3f721 (recursive)
BenchmarkPigeonJSONNoMemo             50          25212587 ns/op         3328296 B/op      86105 allocs/op
BenchmarkPigeonJSONMemo               20          86689562 ns/op        25050390 B/op     131153 allocs/op

// Commit 0850249 (vm - memo is nop)
BenchmarkPigeonJSONNoMemo             20          93411006 ns/op        56357596 B/op     492396 allocs/op
BenchmarkPigeonJSONMemo               20          93271080 ns/op        56357640 B/op     492396 allocs/op

// Commit 1317e07 (vm+stacks 128)
BenchmarkPigeonJSONNoMemo             20          94468129 ns/op        56357323 B/op     492395 allocs/op
BenchmarkPigeonJSONMemo               20          95478650 ns/op        56357121 B/op     492394 allocs/op


// Go1.4 stdlib
BenchmarkStdlibJSON                 2000            861586 ns/op           74094 B/op       1055 allocs/op
