# Test Summary

|     Status      | Count |
|-----------------|-------|
| ✅ Passed       | 7   |
| ❌ Failed       | 9   |
| ⏩ Skipped      | 1   |
| 💥 Parse Errors | 0   |

## Run Failed Tests Locally

```bash
go test ./... -run 'TestTests1_First|TestTests1_FourthTable|TestTests1_Third|TestTests2_First|TestTests2_Third|TestTests3_First|TestTests3_Third'
```

## Failure Details
---

#### `tests/tests_1_test.go`

<details>
<summary>TestTests1_First</summary>

```diff
2022/05/24 10:59:03 Example log
tests_1_test.go:10: failed first
```

</details>

<details>
<summary>TestTests1_FourthTable</summary>

```diff
tests_1_test.go:24: fail filename regexp-catcher
```

</details>

<details>
<summary>TestTests1_FourthTable/subtest_1</summary>

```diff
tests_1_test.go:37: failed sub-test
```

</details>

<details>
<summary>TestTests1_FourthTable/subtest_2</summary>

```diff
tests_1_test.go:37: failed sub-test
```

</details>

<details>
<summary>TestTests1_Third</summary>

```diff
2022/05/24 10:59:03 Example log third 1
2022/05/24 10:59:03 Example log third 2
tests_1_test.go:20: failed third
```

</details>

---

#### `tests/tests_2_test.go`

<details>
<summary>TestTests2_First</summary>

```diff
tests_2_test.go:6: failed first
```

</details>

<details>
<summary>TestTests2_Third</summary>

```diff
tests_2_test.go:13: failed third
```

</details>

---

#### `tests/tests_3_test.go`

<details>
<summary>TestTests3_First</summary>

```diff
tests_3_test.go:6: failed first
```

</details>

<details>
<summary>TestTests3_Third</summary>

```diff
tests_3_test.go:13: failed third
```

</details>

