# Edge AI Landscape and Future Use Cases — Research Report

**Date:** August 2026
**Brief:** Map what is being done in edge AI today, evaluate the thesis that "AI will be incorporated everywhere, much like computers are now everywhere," and predict future use cases.
**Method:** Multi-agent research workflow — six parallel web-research angles (silicon, on-device models/runtimes, consumer products, industrial/physical-world deployments, cloud-vs-edge drivers, forecast landscape), a completeness-critic pass with two follow-up investigations (China's edge AI stack, on-device learning / fleet lifecycle / security), and an adversarially-checked synthesis. 10 agents, 185 findings, ~290 source lookups. Conflicting sources are presented as conflicts rather than silently averaged.

---

*Synthesis of a multi-angle research corpus (silicon, models/runtimes, consumer, industrial, drivers, forecasts, China, fleet lifecycle/security), August 2026. Critic corrections have been applied; unresolved conflicts in the corpus are presented as conflicts.*

---

## 1. Executive Summary

**Edge AI is real, large, and mostly invisible — but it is not what the marketing says it is.** The scaled successes of 2026 are almost all *embedded inference the user never sees*: driver-monitoring cameras, TV picture pipelines, herbicide-spraying booms, checkout-fraud vision, warehouse robot fleets, cardiac rhythm detection on watches. The marquee "on-device AI assistant" category is, by contrast, largely cloud-backed — even Apple, the most on-device-committed vendor, routes hard Siri queries to a custom 1.2T-parameter Gemini in Google Cloud for ~$1B/yr ([Bloomberg](https://www.bloomberg.com/news/articles/2025-11-05/apple-plans-to-use-1-2-trillion-parameter-google-gemini-model-to-power-new-siri)), and Amazon *deleted* Echo's local processing option to feed cloud Alexa+ ([The Register](https://www.theregister.com/2025/03/17/amazon_kills_on_device_alexa/)).

**Headline conclusions:**

1. **The hardware is ahead of the software, and the software is ahead of the demand.** NPU-equipped silicon is approaching majority share of new PCs (~55-59% of 2026 shipments, definition-dependent — see §2) and GenAI-capable phones (~45% of 2026 shipments, [Counterpoint](https://counterpointresearch.com/en/insights/genai-smartphone-share-to-rise-to-45-percent-of-global-shipments-in-2026)), yet mainstream local-LLM tooling still bypasses NPUs, Copilot+ PCs sold ~1.3M units in all of 2024 (~0.5% of the market, [PCWorld/IDC](https://www.pcworld.com/article/2816617/microsofts-copilot-gamble-is-a-bustbut-ai-pcs-still-feel-inevitable.html)), and Microsoft abandoned the NPU-only strategy at Build 2026.

2. **The binding physical constraint is memory bandwidth, not TOPS.** LLM decode speed ≈ bandwidth ÷ active model bytes. Phones (~60-90 GB/s) are capped at 1-4B-active-parameter models; only 256-800+ GB/s unified-memory machines run 70B-class models usably ([memory-bandwidth ladder](https://mlechner.substack.com/p/the-memory-bandwidth-ladder-what)). The industry's response — sparse activation, 2-bit QAT, speculative decoding — is moving the *capability* line faster than the bandwidth line moves.

3. **The "AI everywhere like computers" thesis is directionally right but on a longer clock and a different mechanism than the analogy implies.** AI will diffuse the way *microcontrollers* diffused (invisible, function-specific, uncounted), not the way *PCs* diffused (a visible general-purpose product). Fewer than 1% of the ~30B MCUs shipped annually have AI inference capability today ([Arm](https://newsroom.arm.com/blog/tinyml)); even bullish forecasts (ABI: 4.1B TinyML chipsets/yr by 2031) imply "AI in everything" lands well past 2031.

4. **The equilibrium is hybrid routing, not edge victory.** Every major deployed architecture (Apple device→PCC→cloud, Windows ML, Qualcomm/Google automotive, Samsung Galaxy AI) is on-device-first for narrow tasks with cloud escalation for reasoning. No vendor publishes the on-device/cloud token split — the single most important unmeasured quantity in this debate.

5. **Roughly half the global edge-AI unit map is Chinese and largely decoupled** — Horizon shipped 4M+ ADAS chips in 2025 ([Gasgoo](https://autonews.gasgoo.com/articles/news/horizon-robotics-boasts-577-yoy-surge-in-2025-annual-revenue-2034605818467483649)), Rockchip ~198M edge SoCs, Hikvision/Dahua dominate the world's largest edge-vision install base — running non-NVIDIA silicon, non-CUDA toolchains, and Chinese open-weight models. Western-centric forecasts systematically miss this.

6. **The unsolved problem is not deployment but lifetime.** Cars average 12.6 years on the road ([S&P Global Mobility](https://www.spglobal.com/mobility/en/research-analysis/average-age-of-vehicles-in-the-us-increases-to-126-years.html)); the longest model-update commitments are 7 years; Apple's adapters are version-locked to single OS releases ([Apple developer docs](https://developer.apple.com/apple-intelligence/foundation-models-adapter/)); and frozen on-device models accumulate "forever-day" vulnerabilities. Managing billions of long-lived local models is an unbuilt discipline — the cloud never has this problem.

---

## 2. The Present-Day Map: What Actually Ships, By Layer

### 2.1 Silicon — a four-tier stack with a ~100,000x compute spread

| Tier | Representative parts | What it actually runs | Status |
|---|---|---|---|
| TinyML (µW-mW, ≤4 TOPS) | Syntiant NDP115 (280 µW speech), GreenWaves GAP9 (50 GOPS/50 mW, in earbuds), STM32N6 (600 GOPS), Arm Ethos-U85 in Alif MCUs (SLMs at 36 mW) | Keyword spotting, anomaly detection, low-res vision, first sub-100M SLMs | **Shipping at scale**; Ambiq IPO'd July 2025 ([sources](https://www.cnx-software.com/2025/08/13/alif-ensemble-e4-e6-and-e8-cortex-m55-a32-mcus-and-mpus-feature-ethos-u85-npu-for-small-language-models-slm/)) |
| Consumer NPUs (40-100+ TOPS) | Snapdragon X2 Elite 80 TOPS ([Qualcomm brief](https://www.qualcomm.com/content/dam/qcomm-martech/dm-assets/documents/Snapdragon-X2-Elite-Product-Brief.pdf)), Intel Panther Lake NPU5 50/180 platform TOPS, AMD XDNA2 50-60, Apple M5 (~133 TOPS, third-party est.) | 1-8B LLMs, Studio Effects, photo/scam/translate features | **Shipping in majority of new premium devices; largely idle for third-party LLMs** |
| Discrete edge accelerators ($70-3,499) | Hailo-8/8L ($4-5/TOPS retail via Pi AI HAT+), Hailo-10H (GenAI, 2.5 W), Axelera Metis (~$0.70/TOPS), Jetson Orin Nano Super ($249, 67 sparse TOPS), Jetson AGX Thor (2,070 sparse FP4 TFLOPS, 128 GB @ 273 GB/s, $3,499) | Multi-stream vision; 2-8B LLM/VLM; up to ~120B quantized on Thor | **Shipping, but commercially hard**: Hailo's valuation halved to <$500M ([Calcalist](https://www.calcalistech.com/ctechnews/article/rj000qzaowx)); Google Coral effectively dead ([XDA](https://www.xda-developers.com/intels-cheapest-cpus-do-what-googles-coral-accelerators/)) |
| Automotive/robotics (100-2,500 TOPS) | NVIDIA Drive Thor (in production: BYD, Li Auto, XPeng, ZEEKR), Tesla HW4 (~500 TOPS; AI5 targeting 2,000-2,500 TOPS in 2027 — announcement, not product), Horizon Journey 6 (10-560 TOPS) | End-to-end driving stacks, VLA robot policies | **Scaled in vehicles; pre-scale in robots** |

Key caveat on all TOPS figures: they are peak, mixed-precision (INT8 vs FP4 not comparable), and often thermally unsustainable in phones — and the Intel Panther Lake case (GPU now carries 2.4x the NPU's TOPS) shows "NPU" branding increasingly counts GPU compute ([PCWorld](https://www.pcworld.com/article/2928765/panther-lake-unveiled-a-deep-dive-into-intels-next-gen-laptop-cpu.html)).

**Penetration figures — present the conflict, don't average.** For AI PCs in 2025: Gartner says ~31% ([press release](https://www.gartner.com/en/newsroom/press-releases/2025-08-28-gartner-says-artificial-intelligence-pcs-will-represent-31-percent-of-worldwide-pc-market-by-the-end-of-2025)); Counterpoint's "AI Advanced PC" metric says ~39% ([Counterpoint](https://counterpointresearch.com/en/reports/ai-advanced-pcs-to-surpass-half-of-global-shipments-in-2026)). These are *different definitions quoted as the same quantity* — some firms count any NPU, Copilot+ requires ≥40 TOPS — and the definitional stakes are large: actual Copilot+-class (40+ TOPS) units were only ~0.5-2% of shipments in 2024-early-2025. For GenAI smartphones, 2025 baselines diverge (IDC ~30% vs Counterpoint 36%) and 2028 projections diverge more: **Counterpoint ~54% vs IDC ~70%/912M units** — IDC's definition is materially more aggressive ([Counterpoint](https://counterpointresearch.com/en/insights/genai-smartphone-share-to-rise-to-45-percent-of-global-shipments-in-2026); [IDC via Businesswire](https://www.businesswire.com/news/home/20240730660187/en)). All of these are *supply-side* measures: silicon shipped, not AI used.

### 2.2 Models and runtimes — the fits-on-device line moved from ~1B to ~30B-nominal

What a 2026 flagship genuinely runs locally:

- **Apple AFM (3rd gen, WWDC 2026):** ~3B dense on-device model at ~30 tok/s on iPhone-Pro-class hardware, plus the new **AFM 3 Core Advanced — 20B nominal parameters, sparse, activating only 1-4B per request** ([Apple ML Research](https://machinelearning.apple.com/research/introducing-third-generation-of-apple-foundation-models)). On quantization, keep generations distinct (the corpus conflated them): the *2024* AFM used mixed 2/4-bit palletization averaging ~3.7 bits/weight; the *2025-26* model uses 2-bit quantization-aware training ([arXiv 2507.13575](https://arxiv.org/abs/2507.13575)).
- **Google Gemma 4 E2B/E4B (March/April 2026, Apache 2.0):** multimodal (text+image+audio) in 2-3 GB RAM, ~52 tok/s on Galaxy S26 Ultra GPU via LiteRT-LM, with bundled multi-token-prediction drafters giving 1.8-3x decode speedups ([Google](https://blog.google/innovation-and-ai/technology/developers-tools/gemma-4/); [InfoQ](https://www.infoq.com/news/2026/05/gemma4-multi-token-prediction/)).
- **Qwen3.5-Small (0.8-9B), Phi-4-mini (3.8B), quantized Llama 3.2** round out the field; commodity workloads now include Whisper-class ASR (Apple's SpeechAnalyzer beats Whisper implementations on speed and WER — [MacStories](https://www.macstories.net/stories/hands-on-how-apples-new-speech-apis-outpace-whisper-for-lightning-fast-transcription/)), YOLO-class detection ([YOLO26](https://www.ultralytics.com/blog/ultralytics-yolo26-the-new-standard-for-edge-first-vision-ai)), and sub-1B VLMs.

**A correction on decode-speed claims:** a corpus claim of ~220 tok/s Gemini Nano decode on Snapdragon 8 Elite Gen 5 is almost certainly a *prefill* figure conflated with decode — it contradicts the bandwidth ceiling (phones at ~60-80 GB/s cap 4-bit 2-3B models near 30-55 tok/s) and the measured record of NPU decode weakness (llama.cpp's Hexagon backend launched at 1-3 tok/s in April 2026 — [GitHub](https://github.com/ggml-org/llama.cpp/discussions/8273)). Current Gemini Nano parameter counts are unpublished; the 1.8B/3.25B figures date to Dec 2023. Realistic 2026 phone decode: **~10-52 tok/s for 2-4B models**, device- and runtime-dependent.

**The weakest link is NPU targeting.** Four laptop NPU architectures require four proprietary toolchains; Ollama/LM Studio/stock llama.cpp still default to CPU/GPU; Phi Silica's own design concedes the point (650 tok/s NPU prefill, ~27 tok/s decode partly on CPU — [Microsoft](https://learn.microsoft.com/en-us/windows/ai/apis/phi-silica), with one independent measurement as low as 7.18 tok/s). The abstraction layers that could fix this (Windows ML GA Sept 2025, LiteRT with cross-vendor NPU backends, ExecuTorch 1.0 — already powering on-device AI in Instagram/WhatsApp/Messenger at billions-of-users scale, [PyTorch](https://pytorch.org/blog/introducing-executorch-1-0/)) matured sharply in 2025-26 but are not yet drop-in.

**Quality line:** vendor-acknowledged sufficient for summarization, extraction, rewriting, RAG over local data, image description, simple tool-calling. Long-horizon reasoning, coding, and knowledge-heavy work go to cloud — by every vendor's own routing design.

### 2.3 Consumer — invisible inference works; "AI as the product" mostly doesn't

**Quietly scaled (real usage, no marketing):**
- **EU driver monitoring (GSR2/ADDW):** IR-camera gaze/head-pose inference, processed in a closed loop and deleted by law — on-device by regulation, zero opt-in ([Seeing Machines](https://seeingmachines.com/understanding-advanced-driver-distraction-warning-addw-systems/)). **Date correction:** the July 7, 2026 milestone applies to *new type approvals*; the mandate extends to **all new EU registrations only in July 2028**. The "largest silent deployment of edge AI" framing is right, but its full installed-base effect is a 2028+ story, not a 2026 one.
- **TV picture pipelines** (Samsung ~15 TOPS in-TV NPUs, LG Alpha processors) — always-on, no user action.
- **Live translation** (AirPods fully on-device once packs download — [Apple](https://support.apple.com/en-us/123185); Pixel Voice Translate; Galaxy Interpreter now fully local), **scam-call detection**, photo semantic search.
- **Local security cameras as a paid value proposition:** Eufy/Aqara do face/person/package recognition locally with no subscription — privacy + no fees demonstrably drives purchases.

**Marketing checkboxes with weak usage:** NPU-gated Windows features (Recall re-gated over security a year after launch), most named Apple Intelligence features (adoption evidence conflicts: Morgan Stanley's "80% have used" measures ever-touched, not retention), "AI TV companions."

**High-engagement AI that is cloud in an edge costume:** Meta smart glasses are the runaway hardware hit (~7M units in 2025 alone, $2.15B revenue — [CNBC](https://www.cnbc.com/2026/02/11/ray-ban-maker-essilorluxottica-triples-sales-of-meta-ai-glasses.html)) but essentially all AI inference is cloud-side; the on-device wins are wake-word and the Neural Band's EMG gesture decoding. Alexa+ is aggressively cloud-first across 600M+ devices. The failed AI wearables (Humane — servers shut off Feb 2025, bricking every device; Rabbit R1 — ~5,000 DAU from ~100K sold) were all thin clients for cloud AI ([post-mortems](https://www.digitalapplied.com/blog/ai-product-failures-2026-sora-humane-rabbit-lessons)).

**The consumer pattern:** users adopt edge AI when it is an invisible *property* of a feature (latency, offline, privacy, no subscription) — never when "on-device" is itself the pitch.

### 2.4 Industrial / physical world — a sharp scaled-vs-pilot divide

**Crossed to scale (hard numbers):** Amazon's 1,000,000th warehouse robot (July 2025, ~75% of deliveries robot-assisted — [TechCrunch](https://techcrunch.com/2025/07/01/amazon-deploys-its-1-millionth-robot-releases-generative-ai-model/)); John Deere See & Spray on 5M+ acres in 2025, cutting herbicide ~50% ([Deere](https://www.deere.com/en/news/all-news/see-spray-technology-across-5-million-acres/)); DJI's 600K+ ag drones; Tesla FSD past 10B cumulative supervised miles (May 2026) with HW3's forced retrofit as the canonical staleness failure; Waymo at ~500K paid rides/week ([TechCrunch](https://techcrunch.com/2026/03/27/waymo-skyrocketing-ridership-in-one-chart/)); Everseen checkout vision in 1,700+ Kroger stores; 1,451 FDA-authorized AI medical devices (76% radiology — clearance ≠ deployment, and most run on hospital servers, not true edge); Apple Watch cardiac inference across an installed base of hundreds of millions.

**Still pilots or retreating:** general factory-floor AI (~73-80% of manufacturers in "pilot purgatory" across multiple 2025 surveys); humanoid robots (Unitree *led the world* at ~5,500 units in 2025); grid-edge AI (Utilidata's ~100K Jetson-based meters are DOE-award-funded, not organic); and — the strongest cautionary tale — **5G MEC failed outright**: Microsoft retired Azure Private 5G Core Sept 2025, AWS exited private 5G ([Microsoft lifecycle](https://learn.microsoft.com/en-us/lifecycle/announcements/azure-products-retirement-september-2025)).

**Defense:** Ukraine received ~2.4-3M FPV drones in 2025, but AI terminal guidance remains a minority capability. **A corpus claim that the US Army awarded Anduril a ~$20B Lattice battlefield-integration program should not be treated as established** — the ~$20-22B figure matches the IVAS program ceiling (transitioned from Microsoft to Anduril in 2025) and the single-award characterization is not corroborated by primary contract documents. Doctrine, however, verifiably requires edge inference for denied/disconnected environments.

**China (systematically undercounted in Western corpora):** Horizon at 45.8% ADAS share among Chinese OEMs with VW/CARIAD as anchor investor; Black Sesame +73.4% revenue targeting >10M chips in 2026 (likely across all product lines); Huawei's Kirin NPUs on SMIC 7nm plus an Ascend roadmap through 2028 on self-developed HBM ([Tom's Hardware](https://www.tomshardware.com/tech-industry/artificial-intelligence/huawei-ascend-npu-roadmap-examined-company-targets-4-zettaflops-fp4-performance-by-2028-amid-manufacturing-constraints)); DeepSeek distills and Qwen small models as the default open edge models globally. Note: most edge silicon sits *below* US export-control thresholds, so bifurcation is indirect (SMIC ceiling, entity lists, Beijing's domestic-sourcing mandates) and layered rather than clean.

---

## 3. The Physics and Economics of Where Inference Lives

### 3.1 The physics: bandwidth is destiny (with error bars)

LLM prefill is compute-bound (NPUs shine); decode is memory-bandwidth-bound (NPUs idle). The governing approximation — **max tok/s ≈ bandwidth ÷ active model bytes** — explains the whole device hierarchy: phones (60-90 GB/s) → 1-4B active params; laptops (120-228 GB/s, with Snapdragon X2 Elite Extreme's on-package LPDDR5X the notable move) → 3-8B; Strix Halo/Jetson Thor (256-273 GB/s, 128 GB) → 70-120B quantized; M3 Ultra (819 GB/s) → 70B at usable speed.

**Honest caveat the corpus itself flags:** the formula fails on its own flagship example. 819 GB/s ÷ ~35-40 GB of 4-bit 70B weights gives ~20-23 tok/s theoretical max, yet 25-30 tok/s is reported for M3 Ultra on Llama 3.1 70B. Either the measurement is optimistic (smaller quant, short context) or the formula omits effects (grouped-query attention reducing effective bytes, speculative decoding). Treat the rule as a *ceiling estimator good to ~±30%*, not a law — and treat any decode claim that beats it by 4-7x (like the 220 tok/s phone claim) as prefill conflation.

The datacenter comparison is stark: mobile ~50-90 GB/s vs HBM at 2-3+ TB/s — a 30-50x gap that no roadmap closes ([Chandra, "On-Device LLMs: State of the Union"](https://v-chandra.github.io/on-device-llms/)). NPU TOPS grow ~1.5x/generation, but mobile bandwidth grows far slower. **The industry's real lever is shrinking bytes-per-active-weight**: 2-bit QAT (Apple), INT2 hardware support (Qualcomm), sparse activation (AFM 3 Core Advanced: 7x nominal params at constant active compute), MTP speculative decoding (Gemma 4), BitNet 1.58-bit. This — more than hardware — moved the fits-on-device line from ~1B (2024) to ~30B-nominal (2026).

### 3.2 The economics: colliding curves

Two forces pull in opposite directions:

- **Cloud tokens are collapsing in price** — 9x-900x/yr per capability milestone ([Epoch AI](https://epoch.ai/data-insights/llm-inference-price-trends)); GPT-3.5-class inference fell >280x in 18 months ([Stanford AI Index](https://hai.stanford.edu/ai-index/2025-ai-index-report)). This *weakens* the cost motive for edge.
- **Cloud capacity is rationed** — ~$600-725B hyperscaler capex in 2026 (totals conflict by source), Nadella saying demand outstrips the installed GPU base, Anthropic rationing peak hours, IEA projecting datacenter power 415→~945 TWh by 2030 with ~20% of projects at delay risk ([IEA](https://www.iea.org/reports/energy-and-ai/executive-summary)). This *strengthens* the offload motive: a user's NPU is sunk, user-paid, already-powered silicon with zero marginal cost to the provider. Microsoft's Foundry Local pitch — "unmetered intelligence" — says the quiet part out loud ([Microsoft](https://learn.microsoft.com/en-us/windows/ai/overview)).

Per-query **energy is roughly a wash** (measured ~270 mJ/token on a Hailo-10H vs ~297 mJ/token on an RTX 4050 for the same 1.5B model — [arXiv](https://arxiv.org/html/2603.23640v1) — while a 0.24 Wh cloud Gemini prompt buys a *frontier* answer, [Google](https://cloud.google.com/blog/products/infrastructure/measuring-the-environmental-impact-of-ai-inference)). The energy argument for edge is aggregate grid scarcity, not joules per prompt.

### 3.3 The workload profile where edge is structurally favored

Edge wins when one or more of these hold (and loses otherwise):

1. **Deterministic latency floor <30-50 ms** — robot control loops, AR motion-to-photon (≤20 ms), industrial safety ([The Robot Report](https://www.therobotreport.com/closing-latency-gap-why-physical-ai-requires-edge-first-architectures/)). Note: voice is *not* on this list — cloud speech-to-speech already hits ~320 ms median ([AssemblyAI](https://www.assemblyai.com/blog/low-latency-voice-ai)), so voice assistants are contestable, not edge-locked.
2. **Must run disconnected** — vehicles (tunnels, garages), defense D-DIL, rural/maritime.
3. **Continuous high-duty-cycle sensor streams** where egress dominates — video analytics above ~tens of cameras (vendor-blog magnitudes; directionally credible, unaudited).
4. **Data legally unable to leave the device/premises** — EU health/finance, children's audio, ADDW's mandated closed loop.
5. **Task narrow enough for a ≤10B model** at 80-90% of frontier quality — NVIDIA's own research argues SLMs are sufficient and 10-30x cheaper for most agentic invocations ([arXiv 2506.02153](https://arxiv.org/pdf/2506.02153)).

Cloud structurally retains: frontier reasoning, long-horizon agents, cross-device state, instant model updates, and the developer default (APIs are the path of least resistance). **What moved the boundary toward edge in 2025-26 was capacity economics and physical AI — not privacy ideology and not regulation** (the EU AI Act's high-risk obligations were delayed to Dec 2027/Aug 2028 — [Gibson Dunn](https://www.gibsondunn.com/eu-ai-act-omnibus-agreement-postponed-high-risk-deadlines-and-other-key-changes/)).

---

## 4. Evaluating the Thesis: "AI Will Be Everywhere, Like Computers Are Now"

### Where the analogy holds

1. **The diffusion-by-default mechanism is identical.** Computers became ubiquitous mostly by *embedding* — the average car carries dozens of MCUs nobody calls "computers." NPUs are following exactly this path: they ship in >50% of new phones/PCs regardless of demand, the way MCUs ship in toasters. Supply-side penetration precedes use, then use quietly follows for narrow functions (photo search, translation, upscaling). The 2024 GenAI-smartphone forecasts (IDC 234M, Gartner 240M) essentially came true, and Counterpoint's penetration curve is running slightly *ahead* of its 2023 forecast — the diffusion machinery works.
2. **Capability-per-dollar collapse mirrors Moore's law, faster.** Densing Law — capability density doubling ~3.3 months ([arXiv 2412.04315](https://arxiv.org/abs/2412.04315)) — plus 280x inference-cost collapse means today's frontier is next year's on-device model. Phi-3-mini at 3.8B matched what once took 540B parameters.
3. **The "disappearing technology" endpoint (Weiser, 1991) is already the observed success pattern.** Every scaled consumer edge-AI deployment in §2.3 is invisible. The most-used AI in homes is TV picture processing that no one thinks of as AI — precisely how computing disappeared into thermostats.

### Where the analogy breaks

1. **Computers didn't need a cloud; AI (at frontier quality) still does.** A 1985 PC was self-sufficient. In 2026, every vendor's best product escalates to a datacenter — Apple to a 1.2T Gemini, Google to Private AI Compute, Amazon deleting local mode entirely. Ubiquitous *AI-capable silicon* is not ubiquitous *AI self-sufficiency*, and the bandwidth wall (§3.1) means it won't be soon. The better analogy for the assistant tier is the *browser*: everywhere, but a window onto centralized capability.
2. **Computers were stable artifacts; models are perishable.** A microprocessor from 2005 still runs its washing machine correctly. A frozen 2026 model in a 2038 car is stale, adversarially compromised ("forever-day" sensor attacks — DolphinAttack, speed-sign tape attacks on shipping Mobileye stacks), and version-orphaned (Apple adapters locked to single OS releases; Tesla physically retrofitting HW3). **The microprocessor analogy has no equivalent of model rot**, and no billion-device model-lifecycle platform exists — the cloud-vendor attempts (SageMaker Edge Manager, Azure Percept) were killed.
3. **The historical compute migrations teach oscillation, not a one-way march.** Mainframe→PC decentralized; web re-centralized; smartphone decentralized UX while re-centralizing data; and the *previous edge wave — telco MEC — outright failed*, with infrastructure built ahead of demand and no killer app ([Dell'Oro](https://www.delloro.com/news/mec-market-fails-to-materialize-expectations-lowered-more-than-20-percent-for-2023/)). Benedict Evans' framing fits: new platforms win by reaching new users/use cases, not by displacing the old center — and where inference runs is an open equilibrium, not a destiny ([ben-evans.com](https://www.ben-evans.com/presentations)). Gartner's famous "75% of enterprise data processed outside the datacenter by 2025" was never validated and likely missed.
4. **Economics may favor the meter, not the device.** Altman's "intelligence too cheap to meter" is a *centralized-utility* framing; Gartner's measured money flows agree (AI-optimized IaaS +96% in 2026, inference spend $23.3B exceeding training — [Gartner](https://www.gartner.com/en/newsroom/press-releases/2026-08-10-gartner-forecasts-worldwide-artificial-intelligence-optimized-iaas-spending-to-grow-96-percent-in-2026)). Venture conviction agrees too: pure-play edge-AI chip startups have raised only ~$1.9B cumulatively — a rounding error against one quarter of hyperscaler capex.

### Verdict

**The thesis is correct in its MCU form and misleading in its PC form.** AI inference capability will be as ubiquitous as microcontrollers — embedded, function-specific, invisible, uncounted — on a 10-20-year diffusion clock. But "AI" as users mean it (open-ended intelligence) will remain a *hybrid* service in which the device handles the narrow, instant, private tier and the cloud handles the rest, with the boundary ratcheting outward as small models improve. The right mental model is not "a computer in everything" but "**a reflex arc in everything, a brain on call**."

---

## 5. Predictions

Every prediction names its evidence base and carries a confidence: **High** (extrapolation of a measured, funded, or mandated trend), **Medium** (trend real, timing/magnitude uncertain), **Speculative** (mechanism plausible, no scaled precedent).

### Horizon 1: 1-2 years (through ~2028)

1. **Driver-state inference becomes the largest single consumer edge-AI deployment on Earth.** ADDW extends from new type approvals (July 2026) to *all* new EU registrations in **July 2028**, putting on-device gaze/head-pose models in every new EU car and van (~10-13M vehicles/yr) with zero opt-in. Extrapolates the GSR2 mandate itself ([Seeing Machines](https://seeingmachines.com/understanding-advanced-driver-distraction-warning-addw-systems/)). **High.**
2. **Sparse/MoE becomes the standard on-device architecture: flagship phones ship 20-30B-nominal models activating 1-4B params.** Extrapolates AFM 3 Core Advanced (20B sparse), Gemma 4 26B A4B, and Qwen3 30B-A3B — all shipped 2026. By 2028 this pattern reaches Android flagships via Gemma-class models co-designed with Qualcomm/MediaTek. **High.**
3. **The GPU, not the NPU, becomes the de facto local-LLM engine on PCs; NPUs settle into always-on low-power roles** (wake, vision effects, prefill offload). Extrapolates Microsoft's Build 2026 retreat from NPU-only, Phi Silica on GPU, Intel shipping 2.4x more GPU TOPS than NPU TOPS, and llama.cpp's 1-3 tok/s Hexagon start. **High.**
4. **OS-level hybrid routing becomes the default developer API, and "which model" disappears from app code.** Extrapolates Apple opening the Foundation Models framework to any LLM backend (WWDC 2026, going open source), Windows ML GA, LiteRT consolidation. Apps declare a task; the OS decides device/enclave/cloud. **High.**
5. **On-device scam and fraud interception goes default-on across mid-range Android** (call scam detection, message phishing triage, deepfake-voice flags on calls), driven by Chinese OEMs pushing GenAI NPUs downmarket. Extrapolates Pixel/Galaxy scam detection plus Counterpoint's mid-range GenAI expansion. **High.**
6. **Smart glasses pass 20M units/yr while remaining cloud-inference devices; the first meaningful on-device workloads are EMG gesture decoding, wake/VAD, and frame-selection to cut upload volume.** Extrapolates 7M units in 2025, tripling sales, Meta's 20-30M capacity target — against 30-min Live AI battery life that *forces* local pre-processing. **High** on unit growth, **Medium** on the on-device workload shift.
7. **Waymo-style robotaxi scaling continues (≥1M paid rides/week; Tesla expands supervised autonomy internationally), and China's ADAS attach rate makes L2+ near-universal in new Chinese vehicles.** Extrapolates 10x ridership growth in 22 months, BYD's God's Eye standardization, Horizon's 45.8% share. **High.**
8. **A majority of new Ukrainian/Russian FPV production carries AI terminal guidance by 2028.** Extrapolates mass production starting mid-2025 and jamming pressure making it economically forced. **Medium** (war-reporting quality data).
9. **First "local-only" enterprise laptop policies**: regulated firms (legal, health, defense contractors) mandate on-device transcription/summarization for confidential meetings, using the OS built-ins (SpeechAnalyzer-class ASR + 3-8B summarizers). Extrapolates shipped OS capabilities + EDPB guidance framing local inference as GDPR mitigation. **Medium.**
10. **Humanoids stay pre-scale** (<50K units/yr globally through 2028), with onboard Thor-class compute standard on whatever does ship. Extrapolates Unitree's world-leading ~5,500 units in 2025. **High.**

### Horizon 2: 3-5 years (~2029-2031)

11. **The offline car brain becomes standard equipment:** every mid-range-and-up new vehicle ships a 3-8B onboard assistant (manual Q&A, cabin control, hazard narration) that works in tunnels and garages, with cloud escalation when connected. Extrapolates Snapdragon Cockpit Elite running 3B at ~10 tok/s today, Qualcomm-Google automotive Gemini Nano reference platform, 75M+ vehicles already on Snapdragon cockpits. **High.**
12. **Enterprise edge-AI deployment goes mainstream** — Gartner's projection of two-thirds of enterprises deploying edge AI by 2029 (from 10% in 2025) is directionally validated, concentrated in vision inspection, predictive maintenance, and on-prem SLM appliances for regulated data. Extrapolates the pilot-purgatory backlog + AI Act high-risk obligations landing 2027-28. **Medium** (the same firm's 2018 edge prediction missed).
13. **$10-class AI sensors reach genuine mass deployment:** Ethos-U85/STM32N6-class parts (600 GOPS-4 TOPS, sub-watt, ~$10 BOM) put anomaly detection into pumps, HVAC, elevators, and white goods as a default feature rather than a retrofit — the first true "MCU-ization" wave of AI. Extrapolates first shipping silicon (Alif, ST) + ABI's 4.1B TinyML chipsets/yr by 2031 forecast. **Medium.**
14. **On-watch/on-ear health inference expands from cardiac to multi-signal:** FDA-cleared on-device models for AFib burden (already MDDT-qualified), sleep apnea, and the first respiratory/voice-biomarker screens, running locally for battery and privacy. Extrapolates Apple Watch's regulatory beachhead + sub-mW inference silicon (Syntiant, Ambiq) shipping in hearables. **Medium.**
15. **Hearing augmentation becomes the killer TinyML app:** earbuds do speech-in-noise enhancement, speaker isolation, and live translation continuously at mW power — a use case where the latency floor (~10 ms audio path) makes cloud physically impossible. Extrapolates GAP9 already shipping in earbuds + AirPods Live Translation. **High** for the category, **Medium** for "killer app" status.
16. **Second-order: developers assume a local model the way 2010s apps assumed GPS.** App binaries ship prompts and adapters, not weights; "works offline" becomes a store badge; a generation of apps appears that could not exist with per-token billing — always-on screen understanding, continuous journaling/memory, local RAG over a lifetime of personal files. Extrapolates Apple's free on-device Foundation Models API + Foundry Local's "unmetered intelligence" positioning. **Medium-High.**
17. **Second-order: the on-device semantic index becomes the OS battleground.** Once every file, photo, message, and screen-moment is locally embedded and searchable (Recall done right, Magic Cue generalized), the OS vendor that owns the index owns agent distribution — and antitrust fights follow. Extrapolates Recall/Click-to-Do, Pixel Magic Cue, Apple's on-screen awareness in iOS 26.4. **Medium.**
18. **Robots cross 100K-500K units/yr (warehouse mobile manipulators + narrow-task humanoids), all with onboard VLA inference; general-purpose household humanoids do not arrive.** Extrapolates Amazon's 1M installed base, Figure's onboard Helix, Jetson Thor's 128 GB @ 273 GB/s envelope fitting ~100B-class VLAs. **Medium.**
19. **China's edge stack completes its domestic flip in autos** (near-100% domestic ADAS silicon in new Chinese models around 2027-29 per MIIT guidance) and Chinese open-weight small models (Qwen/DeepSeek lineages) remain the most-deployed edge models globally by unit count. Extrapolates Horizon/Black Sesame/Huawei momentum + Beijing directives + Hugging Face download dominance. **High** for direction, **Medium** for the 100% timeline.
20. **A model-lifecycle regulation moment:** the EU CRA's ≥5-year update floor (from Dec 2027) plus AI Act post-market monitoring produces the first mandated *model* update/recall regime, and the first public "model recall" of a defective on-device model in a safety context. Extrapolates CRA/R156/Article 72 scaffolding + the HW3 retrofit precedent. **Medium.**

### Horizon 3: 5-10+ years (~2031-2036)

21. **The reflex-arc economy: a double-digit percentage of the ~30B+ annual MCU shipments carries inference.** AI becomes a *property of sensors* — every camera, microphone, IMU, and current sensor ships with a model the way sensors ship with ADCs. Nobody will call it AI, exactly as nobody calls a washing machine a computer. Extrapolates Arm's <1%-today baseline + ABI trajectory + $10-class NPU BOMs. **High** on direction, **Medium** on the decade timing.
22. **Yesterday's frontier in your pocket:** if even a diluted Densing Law holds, ~2030-32 phones run models with 2025-26-frontier quality on narrow-domain tasks (not knowledge breadth), making offline expert assistants — medical triage, legal drafting, tutoring — viable in low-connectivity regions. This is the largest second-order prize: **billions of people whose first "computer that talks back" never touches a datacenter.** Extrapolates Densing Law + 280x cost collapse + Gemma-class Apache-2.0 licensing. **Medium.**
23. **Capture-time authenticity inference:** camera pipelines run local deepfake/provenance models at capture and verification models at display, because by then most synthetic media is generated locally too. Extrapolates AI Act transparency rules (in force Aug 2026) + on-device image generation shipping now. **Medium.**
24. **Vehicles become the dominant edge-compute surface:** 2,000+ TOPS (AI5/Thor-successor class) standard in most new cars globally, running driving, cabin agents, and — parked and plugged — offering compute back to homes/grid. The driving part extrapolates Tesla AI5 and Drive Thor roadmaps (**High**); the vehicle-to-home compute part is **Speculative** (no shipped precedent).
25. **True on-device personalization finally ships:** local LoRA/adapter fine-tuning on personal data (your writing voice, your home's sounds, your gait), the piece that is conspicuously *absent* today (Apple adapters are trained off-device; AICore exposes no fine-tuning). Requires solving the version-lock problem — likely via stable adapter interfaces or on-device distillation. Extrapolates Gboard's federated precedent + native LoRA slots already in AFM's architecture. **Speculative** on the 5-year end, **Medium** by 10.
26. **Infrastructure that listens to itself:** decade-battery acoustic/vibration anomaly models embedded in bridges, pipelines, transformers, and rails at sub-mW — predictive maintenance leaving the factory for civil infrastructure. Extrapolates Augury's 1.1B+ machine-hours + Syntiant-class 280 µW inference + Utilidata's meter beachhead. **Medium.**
27. **The disconnected majority thesis:** by ~2035, the *majority of AI inferences on Earth* (by event count, not token count or revenue) happen on devices — because reflex-arc inferences (frames analyzed, sounds classified, gazes tracked) numerically swamp chat tokens. Revenue and frontier capability remain centralized. Extrapolates the surveillance install base (1B+ cameras), ADDW fleets, TinyML unit forecasts. **Medium** — and essentially unmeasurable today, which is itself a finding.
28. **Model heirlooms and digital decay become a consumer-protection issue:** class actions or regulation over devices whose frozen models fail unsafely out of support, echoing the 3G-shutdown and Humane-bricking precedents at far larger scale. Extrapolates 12.6-year car ages vs 7-year update ceilings. **Medium.**

---

## 6. What Could Derail or Reshape These Predictions

1. **Cloud centralization winning on price.** If token prices keep falling 10-900x/yr while capacity catches up (post-2029, per BofA's supply projection), the cost motive for edge evaporates and hybrid routing tilts back toward the API. The measured money already flows centrally (Gartner IaaS inference $23.3B in 2026). *Sensitivity: predictions 16, 22, 27 weaken; the reflex-arc tier (1, 13, 21) survives, since it's latency/physics-driven.*
2. **A frontier capability breakthrough that only scales in datacenters.** If long-horizon agentic reasoning becomes the product and it requires 1T+ active parameters or massive test-time compute, on-device models are demoted permanently to routers and sensor pre-processors. Test-time-scaling-on-edge research exists ([FastTTS](https://arxiv.org/abs/2509.00195)) but carries heavy overhead. *This is the single largest reshaping risk.*
3. **A small-model plateau.** Densing Law is benchmark-defined and may saturate; if sub-10B models stop closing the gap on open-ended tasks, the on-device quality ceiling freezes at "summarize and extract," and predictions 22 and 25 slip a horizon.
4. **The memory-bandwidth wall not moving.** Everything in Horizon 2-3 for phones assumes LPDDR6/on-package memory diffusing downmarket. If phone bandwidth stagnates at ~100 GB/s, the phone tier is permanently capped near 4B active params and the interesting edge migrates to laptops, cars, and robots only.
5. **Energy and grid — in both directions.** Grid scarcity is currently *pro-edge* (offload to already-powered devices); but if datacenter buildout clears its bottlenecks (or SMRs/gas bridge it), that tailwind dies. Conversely, a hard energy crunch would accelerate offload beyond these predictions.
6. **Regulation cutting either way.** The AI Act's high-risk delay (to Dec 2027/Aug 2028) already weakened the near-term compliance-driven edge push; a further softening removes prediction 20's basis. But a major cloud-AI privacy scandal (a FoloToy-at-scale, or an enclave breach) could make "data never leaves the device" a legislated default in children's products, health, and vehicles — accelerating everything in §5.
7. **The lifecycle/security debt coming due early.** A high-profile exploitation of frozen edge models (adversarial attack on a shipping ADAS fleet, mass extraction of app-embedded weights — 41% shipped unprotected per USENIX 2021; EM extraction demonstrated on a commercial edge TPU, [TPUXtract](https://news.ncsu.edu/2024/12/spying-on-ai-models/)) could trigger liability regimes that make long-lived on-device models economically unattractive versus centrally patched cloud inference.
8. **Geopolitical bifurcation deepening.** Full US-China decoupling at the edge layer (DJI Covered-List trigger, mobile SoC controls extending downward, TSMC access loss for Chinese designers) would split the map into two non-interoperable stacks — with Chinese open-weight models continuing to flow outward via open source even where chips cannot, eroding Western on-device model moats.
9. **The MEC precedent repeating.** Silicon shipping ahead of demand is exactly how telco edge failed. If no killer third-party NPU app emerges by ~2028 (none has yet; Dell'Oro's MEC post-mortem and Dell's own "hasn't been what we thought" concession rhyme), the AI-PC/GenAI-phone penetration curves keep rising while mattering less — checkbox ubiquity without use.

---

## 7. What to Watch: Leading Indicators

**The one number that matters most:** any credible measurement of the **on-device vs cloud inference split** (token share or request share). No vendor or analyst publishes it today; third-party claims of 60-80% locally-serviceable queries are estimate-grade. The first vendor telemetry disclosure will reframe this entire debate.

**Software/usage (the demand-side truth):**
- Do Ollama/LM Studio/llama.cpp ship *default-on* NPU backends with >20 tok/s decode? (Currently 1-3 tok/s on Hexagon — the gap to close.)
- Third-party adoption rate of Apple's Foundation Models framework and its promised open-sourcing; count of shipping apps using Phi Silica/ML Kit GenAI APIs.
- Whether a third-party NPU killer app exists by end-2027 — the anti-MEC test.

**Hardware physics:**
- LPDDR6 and on-package memory in mid-range phones/laptops (the X2 Elite Extreme pattern diffusing downmarket) — the only thing that moves the phone decode ceiling.
- Whether Densing Law's ~3.3-month doubling holds through 2027 ([tracking](https://arxiv.org/abs/2412.04315)).

**Deployment milestones (dated, falsifiable):**
- July 2028: ADDW mandatory in all new EU registrations — watch attach execution, not the mandate.
- Waymo at 1M rides/week by end-2026 (its own stated metric); Tesla AI5 volume production in 2027.
- Gartner's enterprise edge-AI curve: is deployment visibly above ~25-30% by 2027 on the way to two-thirds by 2029?
- Humanoid annual shipments crossing 50K (would pull Horizon-2 robotics predictions forward).
- Eufy EdgeAgent (H2 2026) and successors: does "local AI agent" ship as a product category or stay an announcement?

**Economics/market structure:**
- Hyperscaler capex trajectory in 2027 (sustained >$700B vs retrenchment) and whether cloud rationing ends — the offload pressure gauge.
- Survival of merchant edge silicon: Hailo's IPO outcome, Axelera Europa actually shipping H1 2026 volumes — or further consolidation into Qualcomm/NVIDIA.
- Counterpoint vs IDC 2028 GenAI-phone reconciliation (54% vs 70%): which definition reality picks.

**Policy/security tripwires:**
- CRA obligations landing Dec 2027; AI Act high-risk enforcement Dec 2027/Aug 2028; any first "model recall."
- DJI's FCC Covered List status (post-Dec-2025 audit deadline) as the bellwether for demand-side edge bifurcation.
- First in-the-wild (not lab) prompt-injection or adversarial attack against a shipped on-device agent with actuator access — the Black Hat 2025 smart-home hijack moving from demo to incident.

---

*Method note: figures above inherit the corpus's sourcing; where research angles conflicted (AI-PC penetration definitions, GenAI-phone 2028 share, predictive-maintenance market size, M5 Max decode rates, the Anduril award, edge-market sizing spanning $25-48B for 2026), the conflict is stated rather than resolved. Claims flagged in the corpus as training-knowledge, vendor-fed, or single-source (e.g., Gemini Nano throughput measurements, "68% of queries local," Eufy's 63%-faster claim) are excluded from load-bearing conclusions.*