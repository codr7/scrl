# scrl

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

## integers

```
  1 2 3

[1 2 3]
```

## booleans
Booleans have one of two values, `T` or `F`.

```
  T F

[T F]
```

## strings
New strings may be created using `"..."`.

```
"foo"
  
["foo"]
```

## pairs
New pairs may be created using `:`.

```
 1:2
  
[1:2]
```

## deques
Deques are double ended queues of values.
New deques may be created using `[...]`.

```
  [3 2 1]
  
[[3 2 1]]
```

## sets
Sets are ordered collections of values.
New sets may be created using `{...}`.

```
 {3 2 1}
  
[{1 2 3}]
```