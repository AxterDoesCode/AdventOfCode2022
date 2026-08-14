# The Hardware Sabbatical — LLM Inference Chip Roadmap

**Owner:** Alex · **Horizon:** ~12 months · **Written:** August 2026
**Goal:** From "unfamiliar with LLMs and accelerators" to (1) a simulated LLM inference chip with defensible cycle, area, and timing numbers, (2) a novel research contribution built on it, and (3) a PhD reapplication anchored by both.
**Companion documents:** `gaussian-splatting-3d-calls.md` (n/a), `edge-ai-landscape.md` (market context), `taalas-llm-asic.md` (architecture literature map §4, design playbook §5, reading list §7).

---

## Operating rules (read these when tempted to skip ahead)

1. **Exit criteria gate every phase.** No starting phase N+1 until phase N's exit criteria pass. The ladder fails when the interesting part (the prefetcher) starts before the boring part (the GEMV engine) sustains its bandwidth target.
2. **MBU is the scoreboard.** Model Bandwidth Utilization — achieved ÷ theoretical bandwidth — is the metric every block reports from day one. A correct engine at 30% MBU has missed the point of decode hardware.
3. **The golden model is bit-exact or it is useless.** The Python reference implements the *quantized* arithmetic exactly: rounding mode, accumulator width, saturation. "Matches reference" must mean something.
4. **Verify at the layer level, not the story level.** Checkpointed activations from the golden model, compared per layer. A bit error in layer 2 discovered as weird prose at token 40 is a miserable debug loop.
5. **Reading interleaves with building.** Only Phase 0's reading is prerequisite. Everything else is read while simulations run — papers land differently with a design problem open.
6. **Ship a public writeup at the end of every phase.** A blog post or repo README section. By December this is the portfolio the PhD application points at; by month 12 it is the paper's related-work section, already drafted.

**Cadence note:** week counts below assume full-time. At evenings-and-weekends pace, multiply by ~2.5 — and reading starts competing with building instead of interleaving, which is where the multiplier comes from.

---

## Phase 0 — Foundations (Weeks 1–2)

**Build:** the verification harness before any DUT. Verilator + cocotb skeleton; a bit-exact INT8 quantized golden model of a decoder layer in numpy; llama2.c and its TinyStories 15M checkpoint as the executable spec and test-vector source.

**Learn:** Karpathy's "Let's build GPT" + nanoGPT (build one, don't read about one); kipply's *Transformer Inference Arithmetic*; Llama 3 paper architecture section (GQA, RoPE, SwiGLU — the modern reference, not the 2017 paper); *LLM Inference Unveiled* (arXiv 2402.16363) for roofline framing.

**The drill that makes it stick** — derive by hand for Llama-3.1-8B until reflexive:
- Weight bytes at FP16 / INT4
- KV bytes per token: 2 × layers × kv_heads × head_dim × bytes = 64 KB (8-bit)
- Decode tok/s ceiling on an H100 from bandwidth alone
- Prefill FLOPs for a 4K prompt
- Batch size where decode crosses from bandwidth-bound to compute-bound

**Exit criteria:** harness runs a trivial DUT against the golden model with randomized stimulus; the five numbers above derivable without notes.

---

## Phase 1 — Streaming GEMV engine (Weeks 3–5)

The atom of decode. At batch 1 there is no weight reuse, so this is not a systolic array: it is a weight-streaming MAC tree that must sustain full modeled-memory bandwidth. Sequential weight reads, INT8 multiplies, adder tree, configurable-width accumulators, double-buffered activations, a clean streaming memory interface (AXI-flavored).

**Interleaved reading:** Eyeriss (ISCA 2016) — dataflow taxonomy, and why it mostly *doesn't* apply to GEMV; TPU v1 (ISCA 2017).

**Exit criteria:** randomized verification against the quantized golden model passes; **MBU ≥ 85%** against the modeled memory across matrix shapes; assertions on every interface.

---

## Phase 2 — Nonlinear blocks (Weeks 6–7)

RMSNorm, SwiGLU, RoPE (LUT-based complex rotate), and the one that bites: **online softmax** — running max, running sum, accumulator rescaling as the max updates (the FlashAttention trick, in hardware), exp via piecewise-linear LUT.

Write the **fixed-point error budget** as a document: error tolerance per block, measured, so layer-level tolerance is engineered rather than hoped for.

**Interleaved reading:** Brainwave (ISCA 2018) — narrow block floating point in production; OCP MX v1.0 spec.

**Exit criteria:** full layer-level error vs the float reference within the documented budget; softmax stable across adversarial inputs (large logits, ties, long rows).

---

## Phase 3 — Attention unit + KV-cache controller (Weeks 8–12)

The meaty block, and home turf. The KV cache is a banked, paged memory system; PagedAttention's block tables are page tables — build the controller as a TLB-flavored page-table walk over KV blocks. Attention unit: QK^T GEMV against cached keys, online softmax, V accumulation, RoPE applied at ingest. Controller: append-on-decode, multi-turn growth, paging, bank-conflict handling.

**Interleaved reading:** PagedAttention (SOSP 2023), FlashAttention, SpAtten/ELSA (attention sparsity in hardware — also scouting for Phase 6), Groq ISCA 2020/2022.

**Exit criteria:** full attention op verified at layer level against golden model; KV controller correct across multi-turn append/eviction sequences; controller microarchitecture documented (this document seeds Phase 6).

---

## Phase 4 — Integration: the simulated chip (Weeks 13–16)

Sequencer for a full decoder layer; loop all layers of TinyStories 15M; weights streamed from modeled memory; sampling on the testbench side. At ~16 MACs/cycle, ~1M cycles/token — seconds per token in Verilator. **You can watch it write stories.**

Debug methodology: layer-by-layer checkpoint comparison (rule 4), never story-level debugging.

**Exit criteria:**
- End-to-end tokens match the golden model bit-exactly at temperature 0
- Cycle-accurate tok/s **predicted in advance** from bandwidth ÷ model bytes, measured within ~10% of prediction
- A generated story worth putting at the top of the writeup

---

## Phase 5 — Make it a chip (Weeks 17–18)

RTL that runs is not a chip. Push it through Yosys + OpenROAD on an open PDK — SkyWater 130 for realism or ASAP7 for a modern-node story — for post-synthesis area, timing, and power estimates. This is the evidentiary standard accelerator papers use: full RTL + synthesis numbers for the blocks, simulator/model for system claims.

**Deliverable:** the microarchitecture technical report — block diagrams, the error budget, MBU results, synthesis numbers, honest limitations. This is the centerpiece artifact of the PhD application and the skeleton of the eventual paper's methodology section.

**Checkpoint — month ~4–5:** simulated chip complete. Decision point: pick the Phase 6 direction.

---

## Parallel track — PhD reapplication (September–December)

Runs alongside Phases 3–6. **Most Cambridge funding rounds close early December** — operationally ~3 months away, not twelve.

- **September:** send the deferral/decline email (drafted). If deferral granted, everything below relaxes.
- **October:** research statement drafted, anchored on the project: the through-line is *prefetching and speculation applied to LLM inference memory systems* — MPhil skill, new workload, live research seam.
- **November:** statement + writeups to referees; contact the professor with the Phase 0–4 progress (by then there is real progress to show).
- **Early December:** submit. The Phase 5 technical report is the writing sample if timing allows; otherwise the phase writeups are.

---

## Phase 6 — The novel contribution (Months 5–9)

Where the ladder was earned. Pick **one** at the Phase 5 checkpoint:

1. **KV-page prefetcher for sparse attention** *(default pick — closest to the MPhil, clearest gap).* Sparse/retrieval attention (Quest, InfLLM, LongSight) selects which KV pages each token touches — today a demand fetch. Build a predictor of near-future page touches informed by the prefetching literature (stride, correlation, temporal streaming) in front of the Phase 3 controller, with KV tiered across modeled fast/slow memory. Nobody has done a prefetcher-literature-informed treatment.
2. **Speculation-aware KV controller.** Speculative decoding as control speculation: draft = predictor, verify = commit, rejection = squash. A controller that checkpoints and rolls back KV state cheaply on draft rejection; hardware co-design here has barely begun (SpecPIM is the lone early entry).
3. **MoE expert prefetching.** Predict layer N+1 routing during layer N to hide expert fetch from slow memory. Requires adding a small MoE to the platform — more plumbing, hence third.

**Method:** baseline = Phase 5 chip + tiered memory model with demand fetch. Metrics: prediction accuracy/coverage, stall cycles hidden, tok/s uplift, area cost of the predictor (synthesized). Sweep model sizes and context lengths. Honest negative results included — a prefetcher that doesn't help at short context is a finding, not a failure.

**Exit criteria:** a submittable draft + open-source release of platform and predictor.

---

## Phase 7 — Write, submit, iterate (Months 9–12)

- Venue targeting (check exact CFP dates — they shift): ISCA submissions ~November; MICRO ~April; HPCA ~July–August; ASPLOS rolling cycles (~spring/summer/autumn). A workshop at ISCA/MICRO is a legitimate first outing if the full-paper timing is wrong; arXiv the report regardless.
- Iterate on reviews; extend the evaluation where reviewers push.
- If the PhD offer/deferral landed: arrive with a paper in flight and a running platform.

**Stretch goals (only if ahead of schedule, in order):**
- FPGA port of the core (Kria KV260-class board, ~$250) — a physical demo, budget the bring-up pain honestly
- Tiny Tapeout ternary/BitNet micro-engine (~€70–300) — a die photo of your own inference silicon
- Scale the platform to a 1B-class model with INT4/MXFP4 datapath

---

## Risk register

| Risk | Signal | Mitigation |
|---|---|---|
| Scope creep ("build the accelerator") | Generalizing before Phase 4 passes | Operating rule 1; the roadmap is the contract |
| Softmax fixed-point rabbit hole | Phase 2 exceeds 2 weeks | Adopt a known-good online-softmax formulation; error budget caps the pursuit of precision |
| Verification debt | "I'll test it after integration" | Rule 3/4; no block merges without randomized verification |
| Simulation speed wall | Full-model runs too slow to iterate | Reduced-layer configs for iteration; full model for release runs only |
| Reading as procrastination | A week with no commits | Rule 5: reading happens while sims run, not instead of them |
| December deadline collision | Phases 3–4 slip into November | The application needs Phases 0–3 + writeups, not the finished chip; protect the writeups |
| Novelty scoop | A KV-prefetching paper appears | The platform is the moat — redirect to direction 2 or 3 on the same chip |

---

## Cost

Everything through Phase 6 is free: Verilator, cocotb, Yosys, OpenROAD, open PDKs, llama2.c. Money enters only at the stretch goals (~$250 FPGA, ~€70–300 tapeout) — and optionally a GPU cloud hour here and there for golden-model baselines.
