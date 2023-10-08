---
title: Math
description: Hugo provides mathematical operators in templates.
keywords: [math, operators]
categories: [functions]
menu:
  docs:
    parent: functions
toc:
signature: []
relatedfuncs: []
---

| Function        | Description                                                                 | Example                                           |
|-----------------|-----------------------------------------------------------------------------|---------------------------------------------------|
| `add`           | Adds two or more numbers.                                                   | `{{ add 12 3 2 }}` &rarr; `17`                    |
|                 | *If one of the numbers is a float, the result is a float.*                  | `{{ add 1.1 2 }}` &rarr; `3.1`                    |
| `sub`           | Subtracts one or more numbers from the first number.                        | `{{ sub 12 3 2 }}` &rarr; `7`                     |
|                 | *If one of the numbers is a float, the result is a float.*                  | `{{ sub 3 2.5 }}` &rarr; `0.5`                    |
| `mul`           | Multiplies two or more numbers.                                             | `{{ mul 12 3 2 }}` &rarr; `72`                    |
|                 | *If one of the numbers is a float, the result is a float.*                  | `{{ mul 2 3.1 }}` &rarr; `6.2`                    |
| `div`           | Divides the first number by one or more numbers.                            | `{{ div 12 3 2 }}` &rarr; `2`                     |
|                 | *If one of the numbers is a float, the result is a float.*                  | `{{ div 6 4.0 }}` &rarr; `1.5`                    |
| `mod`           | Modulus of two integers.                                                    | `{{ mod 15 3 }}` &rarr; `0`                       |
| `modBool`       | Boolean of modulus of two integers. Evaluates to `true` if result equals 0. | `{{ modBool 15 3 }}` &rarr; `true`                |
| `math.Abs`      | Returns the absolute value of the given number.                             | `{{ math.Abs -2.1 }}` &rarr; `2.1`                |
| `math.Ceil`     | Returns the least integer value greater than or equal to the given number.  | `{{ math.Ceil 2.1 }}` &rarr; `3`                  |
| `math.Floor`    | Returns the greatest integer value less than or equal to the given number.  | `{{ math.Floor 1.9 }}` &rarr; `1`                 |
| `math.Log`      | Returns the natural logarithm of the given number.                          | `{{ math.Log 42 }}` &rarr; `3.737`                |
| `math.Max`      | Returns the greater of all numbers. Accepts scalars, slices, or both.       | `{{ math.Max 1 (slice 2 3) 4 }}` &rarr; `4`       |
| `math.Min`      | Returns the smaller of all numbers. Accepts scalars, slices, or both.       | `{{ math.Min 1 (slice 2 3) 4 }}` &rarr; `1`       |
| `math.Product`  | Returns the product of all numbers. Accepts scalars, slices, or both.       | `{{ math.Product 1 (slice 2 3) 4 }}` &rarr; `24`  |
| `math.Pow`      | Returns the first number raised to the power of the second number.          | `{{ math.Pow 2 3 }}` &rarr; `8`                   |
| `math.Round`    | Returns the nearest integer, rounding half away from zero.                  | `{{ math.Round 1.5 }}` &rarr; `2`                 |
| `math.Sqrt`     | Returns the square root of the given number.                                | `{{ math.Sqrt 81 }}` &rarr; `9`                   |
| `math.Sum`      | Returns the sum of all numbers. Accepts scalars, slices, or both.           | `{{ math.Sum 1 (slice 2 3) 4 }}` &rarr; `10`      |
