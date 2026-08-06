# When Test-First Adds Value

Load this reference only when the appropriate test order is unclear.

## Prefer test-first when

- reproducing a defect proves the regression test targets the real failure;
- examples help define a business rule, parser, validator, or state transition;
- edge cases are easier to discover before committing to an implementation;
- a stable public seam already makes focused tests inexpensive.

Seeing the test fail for the expected reason is useful evidence that it can detect the missing behavior.

## Use another sequence when

- characterization is needed before changing poorly understood legacy behavior;
- the change is mechanical and a reliable existing suite already protects behavior;
- generated output or declarative configuration has a stronger schema/build validator;
- an exploratory spike is needed to discover feasibility;
- the failure cannot be reproduced deterministically and a different verification method is more credible.

Do not delete correct existing work merely because tests were written later. Instead, assess whether the resulting tests can fail for relevant defects, cover the required behavior, and support safe future changes.

## Decision rule

Choose the cheapest verification sequence that provides credible evidence for the actual risk. Document exceptions when a reader could reasonably expect an automated behavior test.
