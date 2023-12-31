# scrl
`scrl` is a scripting language implemented and embedded in Go.

It's designed to complement the host language by adding a convenient meta level on top that's completely under the users control.

## setup

The REPL may be started from a shell like this:

```
$ go run main/scrl.go
scrl v1
  
```

## types
`type-of` may be used to get the type of an expression.

```
  type-of 1
  
[Int]
```
```
  type-of Int

[Meta]
```
```
  type-of Meta

[Meta]
```

### booleans
Booleans have one of two values, `T` or `F`.

```
  or T F

[T]
```

All values have boolean representations; many are unconditionally `T`, `0` and empty sequences are `F`.

```
  and "foo" 42

[42]
```

Logical operatos are short-circuiting.

```
  and T say "hello"

hello
[]
```
```
  or T say "hello"

[T]
```

### integers

```
  + 1 2
  
[3]
```

### symbols
New symbols may be created using `'...` or by interning strings.

```
  = sym "foo" 'foo

[T]
```

### strings
New strings may be created using `"..."`.

```
"foo"
  
["foo"]
```

### pairs
New pairs may be created using `:`.

```
 1:2
  
[1:2]
```

### deques
Deques are double ended queues of values.
New deques may be created using `[...]`.

```
  [3 2 1]
  
[[3 2 1]]
```

### sets
Sets are ordered collections of values.
New sets may be created using `{...}`.

```
 {3 2 1}
  
[{1 2 3}]
```

## flow control
`if` may be used to branch on a condition.

```
  if F "foo"

[]
``

`else` may be used to evaluate code when the condition is `F`.

```
  if F "foo" else "bar"

["bar"]
````

## debugging
`trace` may be used to toggle VM tracing.

```
  trace

[T]
  1 2 3

3 Push 1
5 Push 2
7 Push 3
9 Stop
[T 1 2 3]
```

## benchmarking
`bench` may be used to measure elapsed time for specified number of repetitions.

```
  bench 2 sleep milliseconds 500
  
[1.001917973s]
```

```
 fun fib (n) 
   if < n 2 n else + fib - n 1 fib - n 2

 bench 100 fib 20

[180.872334ms]
```

```
 fun fib(n a b)
   if > n 1 fib - n 1 b + a b else if = n 0 a else b

 bench 10000 fib 70 0 1

[78.846962ms]
```