# The Hardware Sabbatical — LLM Inference Chip Roadmap (Detailed)

**Owner:** Alex · **Horizon:** ~12 months · **Written:** August 2026
**Goal:** (1) a simulated LLM inference chip with defensible cycle, area, and timing numbers; (2) a novel research contribution built on it; (3) a PhD reapplication anchored by both.
**Companions:** `taalas-llm-asic.md` (architecture literature map, design playbook, industry reading list).

**Cadence:** week counts assume full-time; multiply by ~2.5 at evenings-and-weekends pace.

---

## Operating rules

1. **Exit criteria gate every phase.** No phase N+1 until phase N passes.
2. **MBU (achieved ÷ theoretical bandwidth) is the scoreboard** every block reports from day one.
3. **The golden model is bit-exact or it is useless** — rounding mode, accumulator width, saturation all modeled.
4. **Verify at the layer level, never the story level** — checkpointed activations, compared per layer.
5. **Reading interleaves with building.** Only Phase 0's reading is prerequisite; the rest happens while sims run.
6. **Ship a public writeup at the end of every phase** — the December application portfolio and the eventual paper's related work.

---

## Phase 0 — Foundations (Weeks 1–2)

### Build
- **Dev environment:** [Verilator](https://verilator.org) (fast open-source RTL simulator) + [cocotb](https://docs.cocotb.org) (Python testbenches — this is what lets the golden model and the DUT live in one process). Waveforms via GTKWave or [Surfer](https://surfer-project.org).
- **Golden model:** a bit-exact INT8 quantized decoder layer in numpy. Start from [llama2.c](https://github.com/karpathy/llama2.c) — read `run.c` line by line (it *is* the spec: ~700 lines covering RMSNorm, RoPE, GQA attention, SwiGLU, sampling), then port its `runq.c` INT8 path to numpy with explicit rounding/saturation. Test vectors from the TinyStories 15M checkpoint (in the llama2.c README; model background: [TinyStories, arXiv 2305.07759](https://arxiv.org/abs/2305.07759)).
- **Harness skeleton:** cocotb driving a trivial DUT (e.g., an 8-lane MAC) against the numpy reference with randomized stimulus — proves the whole verification loop before any real RTL exists.

### Read (prerequisite — the only phase where reading comes first)
1. [Karpathy, "Let's build GPT: from scratch"](https://www.youtube.com/watch?v=kCc8FmEb1nY) + [nanoGPT](https://github.com/karpathy/nanoGPT) — build one; don't just watch. The single fastest route to transformer mechanics.
2. [The Illustrated Transformer](https://jalammar.github.io/illustrated-transformer/) — the pictures your mental model will reuse forever.
3. [kipply, "Transformer Inference Arithmetic"](https://kipp.ly/transformer-inference-arithmetic/) — **the bridge document for an architect**: params → bytes → FLOPs → latency napkin math. Read twice.
4. [Llama 3 herd of models, arXiv 2407.21783](https://arxiv.org/abs/2407.21783) — §architecture only. The modern reference: GQA, RoPE, SwiGLU. Supporting depth as needed: [GQA (2305.13245)](https://arxiv.org/abs/2305.13245), [RoPE (2104.09864)](https://arxiv.org/abs/2104.09864), [SwiGLU/GLU variants (2002.05202)](https://arxiv.org/abs/2002.05202), [RMSNorm (1910.07467)](https://arxiv.org/abs/1910.07467).
5. [A Visual Guide to Quantization (Grootendorst)](https://newsletter.maartengrootendorst.com/p/a-visual-guide-to-quantization) — scales, zero-points, group quantization; what your golden model must reproduce exactly.
6. [LLM Inference Unveiled, arXiv 2402.16363](https://arxiv.org/abs/2402.16363) — the roofline framing that organizes everything else you'll read.
7. Skim only: [Attention Is All You Need (1706.03762)](https://arxiv.org/abs/1706.03762) — historical context; modern decoders differ in ways that matter to hardware.

### The drill (do on paper until reflexive; check yourself)
For Llama-3.1-8B: weight bytes at FP16 (~16 GB) and INT4 (~4.5 GB with norms/embeddings) · KV bytes/token = 2 × 32 layers × 8 kv_heads × 128 head_dim × 1 B = **64 KB** at INT8 · H100 SXM decode ceiling ≈ 3.35 TB/s ÷ 16 GB ≈ **~200 tok/s** FP16 (×4 at INT4) · prefill FLOPs ≈ 2 × 8e9 × 4096 ≈ **65 TFLOP** for a 4K prompt · batch crossover where decode goes compute-bound (~arithmetic-intensity ridge ÷ 2 FLOPs/byte-class, work it through).

### Exit criteria
- [ ] Harness runs a trivial DUT vs the golden model with randomized stimulus, failures reproduce from a seed
- [ ] The five drill numbers derivable without notes

---

## Phase 1 — Streaming GEMV engine (Weeks 3–5)

### Build
At batch 1 there is no weight reuse — this is **not** a systolic array. Spec:
- N parallel INT8 MAC lanes (start N=16), adder tree, 32-bit accumulators
- Requantization stage: per-group scale, rounding, saturation to INT8 — matching the golden model bit-exactly
- Weight-streaming interface (AXI-Stream-flavored, ready/valid) from a **modeled memory** with configurable bandwidth and latency
- Double-buffered activation SRAM; MBU + stall-cause performance counters from day one
- Testbench: shape sweeps (thin/tall/wide), backpressure injection, SVA assertions on every interface

### Read (while sims run)
1. [Eyeriss, ISCA 2016](https://dl.acm.org/doi/10.1145/3007787.3001177) + [Sze/Chen/Yang/Emer survey, arXiv 1703.09039](https://arxiv.org/abs/1703.09039) — the dataflow taxonomy and the energy hierarchy (DRAM ≈ 100–200× a MAC). Understand it, then note why weight-stationary reuse ≈ 1 for GEMV — the insight that shapes this whole chip.
2. [TPU v1, arXiv 1704.04760](https://arxiv.org/abs/1704.04760) — the canonical industrial datapoint; bandwidth-bound even in 2017.
3. [ZipCPU blog](https://zipcpu.com) — the best practical material on ready/valid handshakes, AXI pitfalls, and formal-friendly interface design. Read the AXI-Stream and skid-buffer posts before freezing your interfaces.

### Exit criteria
- [ ] Randomized verification vs golden model passes across shape sweep
- [ ] MBU ≥ 85% against modeled memory
- [ ] Assertions on every interface; failures reproduce from seed

---

## Phase 2 — Nonlinear blocks (Weeks 6–7)

### Build
- RMSNorm (rsqrt via LUT + Newton step), SwiGLU (SiLU LUT × gate), RoPE (paired rotate, LUT sin/cos)
- **Online softmax**: running max, running sum, accumulator rescale on max update; exp via piecewise-linear LUT
- **The error-budget document**: per-block tolerance vs float reference, measured — layer-level accuracy engineered, not hoped for

### Read
1. [Milakov & Gimelshein, "Online normalizer calculation for softmax", arXiv 1805.02867](https://arxiv.org/abs/1805.02867) — the online softmax algorithm itself, 4 pages.
2. [FlashAttention, arXiv 2205.14135](https://arxiv.org/abs/2205.14135) (+ [FA2, 2307.08691](https://arxiv.org/abs/2307.08691)) — read as *the online-softmax trick at kernel scale*; this is what your Phase 3 attention unit implements in hardware.
3. [I-BERT, arXiv 2101.01321](https://arxiv.org/abs/2101.01321) — integer-only exp/softmax/GELU approximations; directly reusable formulations.
4. [Softermax, arXiv 2103.09301](https://arxiv.org/abs/2103.09301) — hardware softmax design-space study; sanity-checks your LUT sizing.
5. [Brainwave, ISCA 2018](https://www.microsoft.com/en-us/research/publication/a-configurable-cloud-scale-dnn-processing-unit/) — narrow block floating point in production; the precedent for your precision decisions.
6. [OCP Microscaling MX v1.0 spec](https://www.opencompute.org/documents/ocp-microscaling-formats-mx-v1-0-spec-final-pdf) + [arXiv 2310.10537](https://arxiv.org/abs/2310.10537) — where the industry's formats landed; skim now, implement later if you go MXFP4.

### Exit criteria
- [ ] Layer-level error vs float reference within documented budget
- [ ] Softmax stable on adversarial inputs (large logits, ties, long rows)
- [ ] Error-budget doc committed

---

## Phase 3 — Attention unit + KV-cache controller (Weeks 8–12)

### Build
Home turf: the KV cache is a banked, paged memory system.
- **KV controller:** block/page-based storage (PagedAttention block tables *are* page tables — build the lookup as a TLB-flavored walk), append-on-decode, multi-turn growth, eviction, bank-conflict arbitration
- **Attention unit:** RoPE at ingest, QK^T GEMV against cached K, online softmax, V accumulation; GQA-aware (KV heads ≠ Q heads)
- Memory model upgrade: consider [Ramulator 2](https://github.com/CMU-SAFARI/ramulator2) or [DRAMsim3](https://github.com/umd-memsys/DRAMsim3) behind the modeled interface — a real DRAM model makes Phase 6 results credible
- Document the controller microarchitecture — this document seeds Phase 6

### Read
1. [PagedAttention / vLLM, arXiv 2309.06180](https://arxiv.org/abs/2309.06180) (SOSP 2023) — read as a virtual-memory paper wearing an ML costume.
2. [SpAtten, arXiv 2012.09852](https://arxiv.org/abs/2012.09852) (HPCA 2021) and [A3, arXiv 2002.10941](https://arxiv.org/abs/2002.10941) (HPCA 2020); ELSA (ISCA 2021) — the attention-accelerator trilogy; also Phase 6 scouting: note what they *predict* vs demand-fetch.
3. [Groq ISCA 2022 paper](https://groq.com/isca-2022-paper/) — compiler-scheduled determinism; what a cache-free machine looks like.
4. [KIVI, arXiv 2402.02750](https://arxiv.org/abs/2402.02750) — KV quantization; decides your KV storage format.

### Exit criteria
- [ ] Full attention op verified at layer level vs golden model
- [ ] KV controller correct across randomized multi-turn append/eviction sequences
- [ ] Controller microarchitecture documented

---

## Phase 4 — Integration: the simulated chip (Weeks 13–16)

### Build
- Layer sequencer (simple microcoded FSM beats a general controller here), weight-layout plan and memory map document, full-model loop over TinyStories 15M, sampling on the testbench side
- Perf counters everywhere; a run-report script that prints predicted vs measured tok/s and stall breakdown per run
- ~16 MACs/cycle ⇒ ~1M cycles/token ⇒ seconds/token in Verilator: **watch it write stories**, but debug only via layer-checkpoint comparison (rule 4). Keep a reduced-layer config for iteration.

### Read
1. [TPUv4i "Ten Lessons", ISCA 2021](https://dl.acm.org/doi/abs/10.1109/ISCA52012.2021.00010) — the industrial doctrine, best read while making integration compromises; argue with it in the writeup.
2. [gpt-fast](https://github.com/pytorch-labs/gpt-fast) — the cleanest minimal high-performance inference code; a software mirror of your sequencer decisions.
3. Skim [Gemmini](https://github.com/ucb-bar/gemmini) / [Chipyard](https://github.com/ucb-bar/chipyard) — generator structure worth stealing patterns from, even though its GEMM shape isn't yours.

### Exit criteria
- [ ] End-to-end tokens match golden model bit-exactly at temperature 0
- [ ] tok/s predicted in advance from bandwidth ÷ model bytes; measured within ~10%
- [ ] A generated story worth topping the writeup with

---

## Phase 5 — Make it a chip (Weeks 17–18)

### Build
- Synthesis + P&R via [Yosys](https://yosyshq.net) and [OpenROAD-flow-scripts](https://github.com/The-OpenROAD-Project/OpenROAD-flow-scripts); PDK: [SkyWater 130](https://github.com/google/skywater-pdk) (realism) or [ASAP7](https://github.com/The-OpenROAD-Project/asap7) (modern-node story) — run both if cheap
- SRAM macros are the classic snag: [OpenRAM](https://github.com/VLSIDA/OpenRAM) / DFFRAM for SKY130; behavioral + area model for ASAP7 is acceptable if documented
- **Deliverable: the microarchitecture technical report** — block diagrams, error budget, MBU results, area/timing/power, honest limitations

### Read
- OpenROAD flow docs + [Tiny Tapeout](https://tinytapeout.com) docs (the gentlest introduction to the same flow); TPU/Brainwave papers' evaluation sections as templates for how to report numbers.

### Exit criteria / checkpoint
- [ ] Post-synthesis area/timing (and power estimate) published in the report
- [ ] **Month ~4–5 decision point: pick the Phase 6 direction**

---

## Parallel track — PhD reapplication (September–December)

- **September:** send the deferral/decline email (drafted). Deferral granted ⇒ everything relaxes.
- **October:** research statement drafted around *prefetching and speculation applied to LLM inference memory systems*.
- **November:** statement + phase writeups to referees; update the professor with Phase 0–4 progress.
- **Early December:** submit (most Cambridge funding rounds close early December). Phase 5 report is the writing sample if ready; otherwise the phase writeups.
- **Rule:** the application needs Phases 0–3 + writeups, **not the finished chip**. Protect the writeups.

---

## Phase 6 — The novel contribution (Months 5–9)

Pick **one** at the Phase 5 checkpoint. Method for all: baseline = Phase 5 chip + tiered memory model (fast SRAM / slow DRAM via Ramulator) with demand fetch; metrics = prediction accuracy/coverage, stall cycles hidden, tok/s uplift, synthesized predictor area; sweeps over model size and context length; honest negatives reported.

### Direction 1 (default): KV-page prefetcher for sparse attention
Sparse/retrieval attention selects which KV pages each token touches — today, demand-fetched. Build a predictor (stride/correlation/temporal-streaming instincts from the MPhil) in front of the Phase 3 controller.
**Read:** [Quest, arXiv 2406.10774](https://arxiv.org/abs/2406.10774) · [InfLLM, arXiv 2402.04617](https://arxiv.org/abs/2402.04617) · [H2O, arXiv 2306.14048](https://arxiv.org/abs/2306.14048) · LongSight (MICRO 2025, [DOI](https://dl.acm.org/doi/10.1145/3725843.3756062)) — in each, identify the demand-fetch moment your predictor front-runs.

### Direction 2: Speculation-aware KV controller
Draft = predictor, verify = commit, rejection = squash: a controller that checkpoints and rolls back KV state cheaply on draft rejection.
**Read:** [Leviathan et al., arXiv 2211.17192](https://arxiv.org/abs/2211.17192) · [Chen et al., arXiv 2302.01318](https://arxiv.org/abs/2302.01318) · [Medusa, arXiv 2401.10774](https://arxiv.org/abs/2401.10774) · [EAGLE, arXiv 2401.15077](https://arxiv.org/abs/2401.15077) · SpecPIM (ASPLOS 2024) — the lone hardware co-design entry.

### Direction 3: MoE expert prefetching
Predict layer N+1 routing during layer N to hide expert fetch from slow memory. Requires adding a small MoE to the platform (more plumbing — hence third).
**Read:** [Mixtral, arXiv 2401.04088](https://arxiv.org/abs/2401.04088) · [Pre-gated MoE, arXiv 2308.12066](https://arxiv.org/abs/2308.12066) · [LLM in a flash, arXiv 2312.11514](https://arxiv.org/abs/2312.11514) · [SiDA-MoE, arXiv 2310.18859](https://arxiv.org/abs/2310.18859).

**Design-space tooling (any direction):** [Timeloop](https://github.com/NVlabs/timeloop) + [Accelergy](https://github.com/Accelergy-Project/accelergy) for area/energy sweeps beyond what you synthesize.

### Exit criteria
- [ ] Submittable draft
- [ ] Open-source release of platform + predictor

---

## Phase 7 — Write, submit, iterate (Months 9–12)

- Venues (verify CFPs — dates shift): ISCA ~Nov · MICRO ~Apr · HPCA ~Jul–Aug · ASPLOS rolling. Workshop at ISCA/MICRO is a legitimate first outing; arXiv the report regardless.
- Iterate on reviews; extend evaluation where pushed.
- **Stretch goals, in order, only if ahead:** FPGA port (Kria KV260-class, ~$250) · [Tiny Tapeout](https://tinytapeout.com) ternary/BitNet micro-engine (~€70–300; see [BitNet b1.58, arXiv 2402.17764](https://arxiv.org/abs/2402.17764)) · 1B-class model with MXFP4 datapath.

---

## Standing resources (all year)

- **Book:** Sze, Chen, Yang, Emer, *Efficient Processing of Deep Neural Networks* — the accelerator textbook; reference, not cover-to-cover.
- **Industry pulse:** SemiAnalysis, Chipstrat, The Next Platform, EE Times — plus `taalas-llm-asic.md` §7 for the curated set.
- **Benchmark mindset:** MLPerf Inference rules — how serious inference measurement is done.

## Risk register

| Risk | Signal | Mitigation |
|---|---|---|
| Scope creep | Generalizing before Phase 4 passes | Rule 1; the roadmap is the contract |
| Softmax rabbit hole | Phase 2 > 2 weeks | Milakov formulation + I-BERT approximations; error budget caps it |
| Verification debt | "Test after integration" | Rules 3/4; no unverified block merges |
| Sim speed wall | Full-model too slow to iterate | Reduced-layer configs; full model for release runs |
| Reading as procrastination | A week with no commits | Rule 5 |
| December collision | Phases 3–4 slip into November | Application needs Phases 0–3 + writeups; protect writeups |
| Novelty scoop | A KV-prefetching paper appears | Platform is the moat; redirect to Direction 2/3 |

**Cost:** $0 through Phase 6 (Verilator, cocotb, Yosys, OpenROAD, open PDKs, llama2.c); money only at stretch goals.
