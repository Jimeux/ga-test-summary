# Test Summary

|     Status      | Count |
|-----------------|-------|
| ✅ Passed       | 7   |
| ❌ Failed       | 9   |
| ⏩ Skipped      | 1   |
| 💥 Parse Errors | 0   |

## Run Failed Tests Locally

	```bash
	go test ./... -v -run 'TestTests1_First|TestTests1_FourthTable|TestTests1_Third|TestTests2_First|TestTests2_Third|TestTests3_First|TestTests3_Third'
	```

	## Failure Details
---

#### `tests/tests_1_test.go`

<details>
<summary>TestTests1_First</summary>

```diff
=== RUN   TestTests1_First
2022/05/23 20:59:35 Example log
tests_1_test.go:10: failed first
--- FAIL: TestTests1_First (0.00s)
```

</details>

<details>
<summary>TestTests1_FourthTable</summary>

```diff
=== RUN   TestTests1_FourthTable
tests_1_test.go:24: fail filename regexp-catcher
--- FAIL: TestTests1_FourthTable (0.00s)
```

</details>

<details>
<summary>TestTests1_FourthTable/subtest_1</summary>

```diff
=== RUN   TestTests1_FourthTable/subtest_1
tests_1_test.go:37: failed sub-test
--- FAIL: TestTests1_FourthTable/subtest_1 (0.00s)
```

</details>

<details>
<summary>TestTests1_FourthTable/subtest_2</summary>

```diff
=== RUN   TestTests1_FourthTable/subtest_2
tests_1_test.go:37: failed sub-test
--- FAIL: TestTests1_FourthTable/subtest_2 (0.00s)
```

</details>

<details>
<summary>TestTests1_Third</summary>

```diff
=== RUN   TestTests1_Third
2022/05/23 20:59:35 Example log third 1
2022/05/23 20:59:35 Example log third 2
tests_1_test.go:20: failed third
--- FAIL: TestTests1_Third (0.00s)
```

</details>

---

#### `tests/tests_2_test.go`

<details>
<summary>TestTests2_First</summary>

```diff
=== RUN   TestTests2_First
tests_2_test.go:6: failed first
--- FAIL: TestTests2_First (0.00s)
```

</details>

<details>
<summary>TestTests2_Third</summary>

```diff
=== RUN   TestTests2_Third
tests_2_test.go:13: failed third
--- FAIL: TestTests2_Third (0.00s)
```

</details>

---

#### `tests/tests_3_test.go`

<details>
<summary>TestTests3_First</summary>

```diff
=== RUN   TestTests3_First
tests_3_test.go:6: failed first
--- FAIL: TestTests3_First (0.00s)
```

</details>

<details>
<summary>TestTests3_Third</summary>

```diff
=== RUN   TestTests3_Third
tests_3_test.go:13: failed third
--- FAIL: TestTests3_Third (0.00s)
```

</details>

