# Groth16 Minimal Flow Explanation

This document explains the **Groth16 minimal flow** in the context of zero-knowledge proofs (zk-SNARKs). It provides a theoretical understanding of the problem being solved and how it is proven.

---

## What Are We Proving?

We aim to prove that a prover knows a **private input \( x \)** such that it satisfies the equation:

\[
y = x^3 + x + 5
\]

### Inputs:
1. **Private Input (\( x \))**: Known only to the prover.
2. **Public Signal (\( y \))**: Known to both the prover and the verifier.

### Goal:
- The prover must convince the verifier that:
  - They know a valid \( x \) that satisfies the equation.
  - The computation of \( y \) was performed correctly.
  - No information about \( x \) is revealed during the proof.

---

## Theoretical Steps

1. **Circuit Representation**:
   - The equation \( y = x^3 + x + 5 \) is represented as a computational circuit.
   - This circuit defines the constraints for proving the relationship between \( x \) and \( y \).

2. **Prover's Role**:
   - The prover computes \( y \) using their private input \( x \).
   - They generate a cryptographic proof showing that \( x \) satisfies the circuit constraints.

3. **Verifier's Role**:
   - The verifier checks the proof to confirm:
     1. \( y \) is computed correctly according to the circuit.
     2. The proof is valid.
   - Importantly, the verifier learns nothing about \( x \).

---

## Key Properties of zk-SNARKs in Groth16

1. **Completeness**:
   - If the prover knows a valid \( x \) and computes \( y \) correctly, they can generate a proof that always convinces the verifier.

2. **Soundness**:
   - If the prover does not know a valid \( x \), they cannot generate a valid proof to convince the verifier.

3. **Zero-Knowledge**:
   - The proof reveals no information about \( x \), ensuring privacy.

---

## Why Does the Verifier Know \( y \)?

In zk-SNARKs, the **public signals** (\( y \)) are shared with the verifier because:
1. \( y \) is the value the verifier wants to validate.
2. The verifier uses \( y \) to ensure the proof aligns with the computation in the circuit.

However, the verifier does **not** learn the private input \( x \), preserving the prover's privacy.

---

## Real-World Analogy

Imagine a prover claims to know the secret ingredient of a famous recipe but doesn’t want to reveal it. Instead, they prepare the dish and show the result to the verifier. The verifier can confirm the dish is correct without knowing the secret ingredient.

- **Private Input (\( x \))**: Secret ingredient.
- **Public Signal (\( y \))**: The final dish (result).
- **Proof**: A cryptographic assurance that the dish was made correctly.

---

## Summary

The Groth16 minimal flow provides a way to:
- Prove knowledge of a private input (\( x \)).
- Validate the computation of a public signal (\( y \)).
- Preserve the privacy of the private input.

This enables efficient and secure proofs, making zk-SNARKs suitable for applications in blockchain, privacy-preserving protocols, and beyond.