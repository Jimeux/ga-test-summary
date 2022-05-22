# Test Summary

|     Status      | Count |
|-----------------|-------|
| ✅ Passed       | 3   |
| ❌ Failed       | 6   |
| ⏩ Skipped      | 0   |
| 💥 Parse Errors | 0   |

## Run Failed Tests Locally

```bash
go test ./... -v -run 'TestTests1_First|TestTests1_Third|TestTests2_First|TestTests2_Third|TestTests3_First|TestTests3_Third'
```

## Failure Details

---

#### `tests/tests_1_test.go`

<details>
<summary>TestTests1_First</summary>

```bash
=== RUN   TestTests1_First
2022/05/22 20:28:22 Example log
tests_1_test.go:10: failed first
--- FAIL: TestTests1_First (0.00s)
```

</details>

<details>
<summary>TestTests1_Third</summary>

```bash
=== RUN   TestTests1_Third
2022/05/22 20:28:22 Example log third 1
2022/05/22 20:28:22 Example log third 2
tests_1_test.go:19: failed third
--- FAIL: TestTests1_Third (0.00s)
```

</details>

---

#### `tests/tests_2_test.go`

<details>
<summary>TestTests2_First</summary>

```bash
=== RUN   TestTests2_First
tests_2_test.go:6: failed first
--- FAIL: TestTests2_First (0.00s)
```

</details>

<details>
<summary>TestTests2_Third</summary>

```bash
=== RUN   TestTests2_Third
tests_2_test.go:13: failed third
--- FAIL: TestTests2_Third (0.00s)
```

</details>

---

#### `tests/tests_3_test.go`

<details>
<summary>TestTests3_First</summary>

```bash
=== RUN   TestTests3_First
tests_3_test.go:6: failed first
--- FAIL: TestTests3_First (0.00s)
```

</details>

<details>
<summary>TestTests3_Third</summary>

```bash
=== RUN   TestTests3_Third
tests_3_test.go:13: failed third
--- FAIL: TestTests3_Third (0.00s)
```

</details>

