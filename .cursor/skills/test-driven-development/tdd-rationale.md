# When Test-First Adds Value

Load this reference only when the appropriate test order is unclear.

## Prefer test-first when

- reproducing a defect proves the regression test targets the real failure;
- examples help define a business rule, parser, validator, or state transition;
- edge cases are easier to discover before committing to an implementation;
- an existing usecase, handler, repository, or gateway seam makes focused tests inexpensive.

Seeing the test fail for the expected reason is useful evidence that it detects missing behavior.

## Use another sequence when

- characterization is needed before changing poorly understood legacy behavior;
- the change is mechanical and a reliable existing suite protects behavior;
- generated mocks, Wire output, Swagger, or declarative configuration have a stronger generator/build check;
- an exploratory spike is needed to discover feasibility;
- the failure cannot be reproduced deterministically and another verification method is more credible.

Do not delete correct existing work merely because tests were written later. Assess whether tests can fail for relevant defects, cover required behavior, and support future changes.

## Decision rule

Choose the cheapest verification sequence that provides credible evidence for the actual risk. Document exceptions when an automated behavior test would normally be expected.
