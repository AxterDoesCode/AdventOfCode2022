# Real-Time Gaussian Splatting for 3D Video Calls — Research Report

**Date:** August 2026
**Question:** Is real-time Gaussian splatting the right technology for real-time 3D Zoom-style video calls — and is there a better solution?
**Method:** Multi-agent research workflow — six parallel web-research angles (splatting state of the art, industry telepresence systems, live-pipeline feasibility, alternative approaches, academic literature, deployment reality), a completeness-critic pass with two follow-up investigations (group-call scaling, commodity-hardware coverage), and an adversarially-checked synthesis. 10 agents, 122 findings, ~274 source lookups. Conflicting sources are surfaced in the appendix rather than silently resolved.

---

## 1. Verdict

**Part one — is splatting the right technology?** Yes for the *rendering substrate*, no for the *end-to-end pipeline*. Gaussian splatting has effectively won the race for how a 3D human should be *rendered* on a receiving device — even competing research lineages (Meta's mesh-based codec avatars, Microsoft's 2D VASA line) have converged onto Gaussian or hybrid mesh+Gaussian output representations ([Gaussian Pixel Codec Avatars](https://arxiv.org/pdf/2512.15711), [VASA-3D](https://arxiv.org/pdf/2512.14677)). But splatting does **not** solve the parts of a video call that are actually hard: live capture from commodity hardware, real-time reconstruction, transmission at consumer bitrates, and a display that makes 3D worth having. The only demonstrated live-authentic splat call ([VoluMe](https://microsoft.github.io/VoluMe/), ICCV 2025) needs an RTX 4090 Mobile for 28 FPS on **one** feed — hardware present in roughly 1–3% of the actual videoconferencing endpoint base ([JPR Q3 2025](https://www.jonpeddie.com/news/q325-pc-gpu-shipments-increased-by-2-5-from-last-quarter-which-might-suggest-a-creep-forward/), [Steam survey](https://www.thefpsreview.com/2026/07/06/steam-hardware-survey-june-2026-a-laptop-gpu-just-topped-the-chart-for-the-first-time/)).

**Part two — is there a better solution?** There is no single better solution; there are three better-fitting architectures depending on the endpoint, and splats are a component in two of them:

- **Consumer 2D screens (the actual Zoom market):** 2D neural rendering / generative face codecs win today — 10–30 kbps, single webcam, being standardized by MPEG/JVET as [Generative Face Video Coding](https://arxiv.org/pdf/2311.02649). A 3D pipeline feeding a flat monitor has an unproven marginal payoff (see §6 caveat).
- **Enterprise rooms:** multi-camera capture + **cloud** neural volumetric reconstruction + light-field display — Google Beam's architecture ([blog.google](https://blog.google/technology/research/project-starline-google-beam-update/)) — wins, at $24,999/room plus undisclosed license fees.
- **Headset-to-headset:** pre-enrolled avatars driven by low-rate signals (Apple Personas, Meta's Gaussian codec avatars) win — kbps-scale bandwidth, 72–90 FPS on mobile silicon — at the cost of live authenticity and a tiny install base.

The honest one-liner: **splatting is the right renderer, the wrong codec, and premature as a consumer product.**

---

## 2. What a Real-Time 3D Call Actually Requires

These are the evaluation criteria, with the bars they must clear:

| Criterion | Bar | Basis |
|---|---|---|
| **One-way latency** | <150 ms (150–400 ms degraded; >400 ms unacceptable) | [ITU-T G.114](https://www.itu.int/rec/dologin_pub.asp?lang=e&id=T-REC-G.114-200305-I!!PDF-E&type=items) |
| **Frame budget** | ~33 ms/frame at 30 fps for every live per-frame stage | Derived from G.114 + 30 fps standard |
| **Bandwidth** | ~2–4 Mbps symmetric = a good 1080p 2D call; consumer uplinks cannot sustain 30–100 Mbps bidirectionally | [Zoom docs](https://broadbandnow.com/guides/internet-speed-zoom), [Starline paper](https://dl.acm.org/doi/10.1145/3478513.3480490) |
| **Capture** | Ideally one webcam, no enrollment; each added camera/sensor/scan shrinks the market | — |
| **Display** | A 2D screen shows no parallax/stereo; 3D value requires light-field ($2k–$25k), autostereo, or a headset | [Looking Glass](https://lookingglassfactory.com/looking-glass-27), [HP Dimension](https://www.roadtovr.com/hp-dimension-google-beam-price-release-date/) |
| **Install base / min spec** | ~75% of PCs are iGPU-only ([JPR](https://www.jonpeddie.com/news/q325-pc-gpu-shipments-increased-by-2-5-from-last-quarter-which-might-suggest-a-creep-forward/)); >70% of phones sold are sub-premium ([Counterpoint](https://counterpointresearch.com/en/insights/global-smartphone-shipments-grew-2-percent-YoY-in-2025)); Vision Pro ~475–600k units cumulative ([reporting](https://virtual.reality.news/news/apple-vision-pro-sales-plunge-95-as-production-halts/)) |
| **Group scaling (N-way)** | Cost must scale to 4–9+ tiles; streams must relay through commodity SFUs | [NOSSDAV 2026](https://dl.acm.org/doi/10.1145/3798065.3798070) |
| **Authenticity/anti-fraud** | Photoreal tracker-driven avatars are consented deepfake rigs; $25.6M Arup fraud, ~$1.1B US deepfake losses in 2025 | [deepstrike.io](https://deepstrike.io/blog/deepfake-statistics-2025) |

---

## 3. How Splatting Scores Today

The critical move is to separate splatting's **four pipeline roles**, because they score very differently.

### 3a. Rendering (receiver-side): solved

- Desktop: 200–350 FPS ([QUEEN](https://research.nvidia.com/labs/amri/publication/girish2024queen/), [3DGStream](https://arxiv.org/abs/2403.01444)); browser WebGPU: 50–147 FPS ([SuperSplat](https://radiancefields.com/supersplat-ships-compute-based-webgpu-rendering-and-automatic-streamed-lod)); flagship phones: 60–116 FPS ([Mobile-GS](https://xiaobiaodu.github.io/mobile-gs-project/)); VR: [SqueezeMe](https://arxiv.org/abs/2412.15171) runs **3 full-body photoreal avatars at 72 FPS on a standalone Quest 3**, [TaoAvatar](https://arxiv.org/abs/2503.17032) hits 90 FPS stereo on Vision Pro, [HRM2Avatar](https://arxiv.org/abs/2510.13587) 120 FPS on iPhone 15 Pro Max.
- **Caveat:** these are flagship devices. On integrated GPUs large splat scenes fall to ~10 FPS and on a 2023 flagship Android GPU (Adreno 740) to 5–20 FPS ([WebSplatter](https://arxiv.org/abs/2602.03207)); avatar-scale scenes are several times smaller, so 30 FPS on recent iGPU/M-series is plausible but **mid-range Android remains below real time**, and NeRF's 10–15 FPS interactive rates are why the field migrated to splats at all ([FlashAvatar](https://arxiv.org/html/2312.02214)).

### 3b. Live reconstruction (sender-side): the wall

- **Optimization-based** per-frame reconstruction is ~100–360x over budget: 12 s/frame ([3DGStream](https://arxiv.org/abs/2403.01444)), 2.67 s/frame ([Instant Gaussian Stream](https://arxiv.org/abs/2503.16979)), <5 s/frame ([QUEEN](https://research.nvidia.com/labs/amri/publication/girish2024queen/)). Not usable live.
- **Feed-forward** prediction works but is expensive: [VoluMe](https://arxiv.org/abs/2507.21311) predicts per-frame head-region Gaussians from a single 2D webcam at **28 FPS on an RTX 4090 Mobile** — the first live monocular splat call, demoed two-way, but head-only, limited novel-view range, one feed ≈ one high-end GPU. [GPS-Gaussian](https://arxiv.org/abs/2312.02155) gets 2K at 25 FPS but needs calibrated stereo rigs.
- Note the "real-time" trap: most 3DGS streaming papers mean real-time *rendering* of offline-encoded content, not live capture-to-display.

### 3c. Transmission: splats are not a codec

- Raw dynamic 3DGS: ~49 MB/frame (~11.8 Gbps at 30 fps) ([arXiv 2507.14432](https://arxiv.org/html/2507.14432v1)). Best compressed streams: QUEEN ~0.7 MB/frame (~168 Mbps), [V³](https://arxiv.org/abs/2409.13648) ~0.5 MB/frame (~125 Mbps, offline-encoded), [GIFStream](https://xdimlab.github.io/GIFStream/) 30 Mbps (offline encoding). All 10–1000x above a 2–4 Mbps Zoom call, and **no real-time encoder exists at any of these points**.
- The fix every working system uses: **don't ship splats.** Ship 2D video and reconstruct at the receiver (VoluMe), or ship driving parameters for a pre-built avatar — [Mon3tr](https://arxiv.org/abs/2601.07518) streams <0.2 Mbps of motion/appearance features over WebRTC with ~80 ms end-to-end latency and ~60 FPS on a Quest 3 receiver. That is *below* 2D-call bitrates, matching the measured finding that Apple's spatial Personas use less bandwidth than 2D video ([arXiv 2405.10422](https://arxiv.org/html/2405.10422v1)).

### 3d. Avatar construction (enrollment): fast but tiered

- One-shot: [LAM](https://github.com/aigc3d/LAM) builds an animatable Gaussian head from one photo — sources conflict on timing (**1.4 s on A100** vs **"under 1 second"**; seconds-scale either way) — rendering 110+ FPS on a phone via WebGL. [GAGAvatar](https://arxiv.org/abs/2410.07971) does single-image reenactment, but sources contradict on its speed: **246 FPS on RTX 3090** (project-page framing) vs **67 FPS on A100** (paper figure with pre-tracked drivers); the gap is likely rendering-only vs full reenactment, and third-party tests report 30–60 FPS on consumer GPUs with online face tracking cutting throughput 50–70%. Treat headline FPS in this subfield as rendering-only until proven otherwise.
- Quality tiers are real: one-shot avatars hallucinate occluded regions, bake lighting, and fail extreme expressions; rig-grade quality ([RGCA](https://shunsukesaito.github.io/rgca/), [TaoAvatar](https://pixelai-team.github.io/TaoAvatar/)) still needs multi-view capture or guided phone scans. A 2026 industry assessment still rates splat avatars "research-grade" vs production 2D neural rendering ([Forasoft](https://www.forasoft.com/blog/article/interactive-ai-avatar-development)).
- The real live bottleneck is often the **driver**, not the splats: research FLAME trackers add 200–300 ms/frame; deployed systems substitute ARKit/MediaPipe blendshapes or audio ([Meta's audio-driven codec animation: <15 ms GPU time](https://arxiv.org/abs/2510.01176)).

### 3e. Group calls: only the driven-avatar class scales

No published 3+ party splat call exists. Derived arithmetic: VoluMe-class feed-forward reconstruction scales to ~3–7 FPS/feed at 4–9 feeds on one GPU — non-viable — while driven avatars cost (N−1)×0.2 Mbps and rendering headroom (published mobile ceiling: 3 avatars/Quest 3; [5 mesh avatars saturated a Quest 2](https://arxiv.org/abs/2104.04638)). Motion-feature streams relay through commodity SFUs today ([LiveKit data tracks](https://livekit.com/blog/livekit-data-tracks-realtime-streaming)); volumetric streams needed custom research SFUs ([NOSSDAV 2026](https://dl.acm.org/doi/10.1145/3798065.3798070)). Most tellingly, **Google's own Beam group-meeting experiment renders remote participants as 2D life-size video around a virtual table**, not N-way volumetric capture ([blog.google](https://blog.google/innovation-and-ai/models-and-research/google-research/google-beam-group-meetings/)).

**Scorecard for splatting:** latency ✓ (80–150 ms demonstrated in avatar/feed-forward configurations — [Mon3tr](https://mon3tr3d.github.io/), [Tele-Aloha](https://dl.acm.org/doi/10.1145/3641519.3657491)); rendering ✓ on flagship hardware; bandwidth ✓ *only* in driven-avatar form, ✗ as a streamed format; capture ✗ (webcam works only at reduced fidelity on a 4090-class GPU); display ✗ (no consumer 3D display base); install base ✗ (~1–3% meet the demonstrated live-capture spec); group scaling ✗ except driven avatars; authenticity ⚠ (driven avatars are structurally deepfake-adjacent).

---

## 4. The Alternatives, Scored on the Same Criteria

### 4a. 2D neural talking heads / generative face codecs

[Maxine face-vid2vid](https://nvidianews.nvidia.com/news/nvidia-announces-cloud-ai-video-streaming-platform-to-better-connect-millions-working-and-studying-remotely) (~10x bitrate savings, vendor-claimed), [LivePortrait](https://github.com/KwaiVGI/LivePortrait) (~22 kbps "perceptually lossless" per an [independent analysis](https://mlumiste.com/technical/liveportrait-compression/), but needs a 4090-class receiver and is ~20x slower than real time on Apple Silicon), [VASA-1](https://arxiv.org/pdf/2404.10667) (40 FPS from a photo + audio; withheld by Microsoft over deepfake risk). JVET is standardizing this as SEI on VVC with >70% bitrate savings below 10 kbps ([GFVC survey](https://arxiv.org/pdf/2311.02649)) — the strongest institutional signal about the near-term path. **Weaknesses:** narrow pose range, head-and-shoulders only, hallucination/identity-drift risk, no true 3D output for stereo displays.

### 4b. NeRF avatars

[INSTA](https://zielon.github.io/insta/)-class heads train in <10 min but render at 10–15 FPS interactive vs 300+ FPS for equivalent Gaussian avatars ([FlashAvatar](https://arxiv.org/html/2312.02214)) with no compensating advantage. **Effectively superseded; not a live contender.**

### 4c. Pre-enrolled codec avatars (mesh/latent or Gaussian)

Meta [Codec Avatars](https://arxiv.org/pdf/2104.04638) and Apple Personas: kbps-scale bandwidth, full novel-view freedom, mobile-renderable. Costs: per-person enrollment, premium sensor hardware on both ends, and no live authenticity (your avatar wears what it was enrolled in). Apple's [visionOS 26 Personas](https://www.apple.com/newsroom/2025/06/visionos-26-introduces-powerful-new-spatial-experiences-for-apple-vision-pro/) (announced June 9, 2025; shipped Sept 15, 2025) reversed the feature's early ridicule ([Gizmodo](https://gizmodo.com/okay-apple-vision-pros-spatial-personas-dont-look-like-trash-anymore-2000614401)); reporting suggests they are Gaussian-splat-based, though Apple has never technically confirmed this.

### 4d. True-capture volumetric (RGBD mesh / multi-camera neural)

[Holoportation](https://arxiv.org/pdf/2209.01982): 1–2 Gbps raw, ~30–50 Mbps mobile-compressed, visibly worst photorealism. Starline 2021: photoreal, 30–100 Mbps, **105.8 ms mean end-to-end** with 4 workstation GPUs per booth ([TOG paper](https://hhoppe.com/starline.pdf)) — the proof the latency budget closes. **Important correction:** that 105.8 ms was measured on the *2021 prototype with local GPUs*. Google Beam moved reconstruction to Google Cloud, **no public end-to-end latency figure exists for the shipping product**, and the added WAN round-trip could plausibly threaten the <150 ms budget.

### 4e. Stereoscopic video

~2x 2D bitrate, near-zero compute, full photorealism, but zero motion parallax, needs a stereo camera and a single-viewer eye-tracked display ([Acer SpatialLabs](https://www.acer.com/us-en/spatiallabs)). No mainstream service ships it.

### Comparison table

| Approach | Latency | Bandwidth | Capture | Display needed | Min endpoint spec / reach | Group scaling | Authenticity |
|---|---|---|---|---|---|---|---|
| **Live splat prediction** (VoluMe) | ~36 ms/frame stage; E2E unpublished | 2D-video rates (reconstruct at receiver) | 1 webcam ✓ | 2D or 3D | RTX 4090 Mobile ≈ 1–3% of endpoints ✗ | ~1 feed/GPU ✗ | ✓ live-authentic |
| **Driven splat avatar** (Mon3tr / LAM / SqueezeMe) | ~80 ms E2E ✓ | <0.2 Mbps ✓✓ | Enrollment (multi-view → one photo) ⚠ | 2D, VR ✓ | Flagship phone/Quest 3 receiver; ~20–40% could render ⚠ | 3–5 avatars on mobile; SFU-trivial ✓ | ✗ deepfake-adjacent |
| **2D generative codec** (GFVC/LivePortrait/VASA) | real-time on 4090; standardizing | 10–30 kbps ✓✓✓ | 1 webcam/photo ✓✓ | 2D ✓ | 4090-class receiver today; distillable ⚠ | Like 2D video ✓ | ✗ hallucination risk (VASA-1 withheld) |
| **NeRF avatar** | 10–15 FPS ✗ | kbps ✓ | monocular video ✓ | any | consumer GPU | ✗ | ⚠ |
| **Codec avatar / Persona** | real-time on-device ✓ | kbps ✓✓ | headset sensors + enrollment ✗ | headset both ends ✗ | <1M Vision Pros; Quests lack face tracking ✗ | 3–5/mobile device ⚠ | ⚠ on-device-bound (safest shipped design) |
| **Cloud volumetric** (Beam) | 105.8 ms on 2021 local prototype; Beam E2E unpublished ⚠ | 30–100 Mbps ✗ | 6-camera rig ✗ | 8K light field ✗ | $24,999/room + license ✗ | 1:1 only (group = 2D compositing) ✗ | ✓ live capture |
| **RGBD mesh** | real-time | 30 Mbps–2 Gbps ✗ | depth cam(s) ⚠ | any 3D | prosumer | ✗ | ✓ |
| **Stereo video** | trivial ✓ | ~2x 2D ✓ | stereo pair ⚠ | eye-tracked 3D, 1 viewer ✗ | niche monitors ✗ | ✗ | ✓✓ |

---

## 5. What the Industry Leaders Chose, and What It Reveals

- **Google** chose multi-camera RGB + **cloud** "AI volumetric video model" + light-field display, enterprise-only, one-to-one ([Beam](https://blog.google/technology/research/project-starline-google-beam-update/); [HP Dimension $24,999](https://www.svconline.com/proav-today/first-impressions-slimmed-down-project-starline), 50–100 Mbps symmetric — sources conflict on the exact minimum). Google has never disclosed whether the model is splat-based. Revelation: even the best-resourced team could not make endpoint-local 3D calls consumer-viable — and its ship date slipped from late 2025 to "later in 2026" ([Pocket-lint](https://www.pocket-lint.com/hp-imagine-2026-google-beam/)).
- **Apple** chose on-device pre-enrolled neural avatars driven by headset sensors ([Personas](https://applemagazine.com/vision-pro-persona/)) — the only *shipped consumer* photoreal 3D call feature — bound to live driving by the enrolled owner's own face, the structurally safest anti-deepfake design shipped. Revelation: enrollment + parameters beats streaming geometry when both ends are premium hardware; but the addressable base is sub-1M.
- **Meta** pivoted codec avatars to Gaussians ([SqueezeMe](https://arxiv.org/abs/2412.15171)) and is mass-collecting training data ([Project Warhol, $50/hr](https://www.uploadvr.com/meta-project-warhol-codec-avatars-training-paying/)), yet has shipped nothing photoreal because **no current Quest has face/eye tracking** ([UploadVR](https://www.uploadvr.com/meta-codec-avatars-might-be-coming-to-quest/)). Revelation: the blocker is driving sensors, not splat rendering.
- **NVIDIA** chose 2D-to-3D lifting (triplane NeRF from one webcam, sub-100 ms on RTX; [Maxine 3D](https://research.nvidia.com/labs/nxp/3dvc-siggraph-etech/)) — still early-access, not GA. **Microsoft** retired Mesh (Dec 1, 2025; [Computerworld](https://www.computerworld.com/article/4100319/microsoft-retires-mesh-app-launches-immersive-spaces-for-teams.html)) and handed its volumetric capture business to Arcturus; its research bets are VoluMe (splats) and VASA (withheld). **Cisco** quietly shelved Webex Hologram. **Zoom** built nothing 3D and positioned itself as the fabric on others' endpoints ([Zoom on Vision Pro](https://news.zoom.com/zoom-launches-new-app-for-apple-vision-pro-to-make-hybrid-collaboration-more-immersive/)).

The pattern: nobody streams splats; splats won as the local rendering format (Meta avatars, Hyperscape, Gracia's 4DGS pipeline, likely Apple Personas), while the transmission layer is either ordinary 2D video, driving parameters, or a proprietary cloud reconstruction. And two majors (Microsoft, Cisco) simply exited, which is itself evidence about consumer demand at current cost.

---

## 6. The Honest Answer by Scenario

**Consumer 2D screens (laptops/phones — ~99% of calls):** *Splatting loses; 2D wins — with one honest caveat.* On a flat monitor the deliverable 3D benefits reduce to eye contact, framing, and (single-viewer) head-tracked parallax. Eye contact and framing already ship as near-free 2D AI features ([NVIDIA Broadcast Eye Contact](https://www.nvidia.com/en-us/geforce/news/jan-2023-nvidia-broadcast-update), Windows Studio Effects). The caveat: the claim that these 2D features capture "most of the measurable benefit" is an **assumption, not a finding** — Google's measured gains (+28% recall, −31% fatigue, 87% preference; [research.google](https://research.google/blog/how-project-starline-improves-remote-communication/)) were obtained on a light-field display and were never decomposed into eye-contact vs parallax vs stereo components. What is *not* in doubt: the demonstrated splat sender spec (RTX 4090 Mobile) excludes ~97–99% of endpoints, and no end-to-end splat pipeline has ever been shown on Apple Silicon, iGPUs, or mid-range Android. The rational consumer product in 2026 is 1080p video + GFVC-style compression + 2D AI polish.

**Enterprise booths:** *Cloud volumetric wins, splats optional inside the black box.* Beam is the only product whose display can actually deliver the measured presence gains, and it justifies $25k/room only for high-stakes recurring 1:1s. Its unknowns are real: undisclosed license pricing, unpublished production latency, one-to-one-only true 3D, and a group mode that falls back to 2D compositing. [Tele-Aloha](https://dl.acm.org/doi/10.1145/3641519.3657491) (4 RGB cameras, one consumer GPU, <150 ms, autostereo display) shows this class can scale down toward prosumer cost — using a splatting rasterizer, notably.

**VR/headset-to-headset:** *Driven Gaussian avatars win outright — this is splatting's home turf.* Free-viewpoint stereo-consistent rendering is mandatory here (2D methods can't provide it), bandwidth is sub-Mbps, and rendering fits mobile silicon (72–90 FPS). Apple ships it today (sub-1M users); Meta ships it when a face-tracked headset arrives (expected ~2026). The trade is authenticity: these are enrolled puppets, which drags in the deepfake problem — the emerging design line is binding avatars to live on-device driving by the enrolled owner, since watermarks are already dismissed as [security theatre](https://ironscales.com/blog/zooms-ai-avatar-watermark-is-security-theatre-and-attackers-already-know-it).

### What would have to change for splatting to win outright

1. **Sender-side feed-forward reconstruction must get ~10x cheaper** — VoluMe-class quality at 30 FPS on M-series/iGPU/NPU silicon (no port even exists today; the closest 2D analog runs ~20x below real time on Apple Silicon).
2. **A real-time dynamic-splat codec at <5 Mbps**, or maturation of joint coding ([GS-SCNet](https://arxiv.org/abs/2604.25330) is at ~19 FPS and −75% BD-rate — close but below the bar); until then splats remain a rendering format, not a transmission format.
3. **A consumer 3D display or face-tracked headset install base** — the single biggest blocker; without it, 3D reconstruction feeds a 2D monitor and competes with free 2D AI features. Watch Looking Glass's $2k line and Meta's next headset.
4. **One-shot avatar quality closing on rig quality** (the [UIKA](https://zijian-wu.github.io/uika-page/)/FlexAvatar universal-prior wave is the path), plus driver-side tracking under ~30 ms, so enrollment stops being a market filter.
5. **An authenticity standard** (live-capture binding or C2PA-style provenance for calls) so photoreal 3D avatars are deployable at scale without becoming fraud infrastructure.

Items 1–2 look like 2–4 years of normal research progress. Item 3 is the genuinely uncertain one — and it, not Gaussian splatting, is what currently decides whether "3D Zoom" exists at all.

---

## Appendix: Source Conflicts Surfaced (Not Silently Resolved)

- **GAGAvatar:** 246 FPS (RTX 3090) vs 67 FPS (A100) — likely rendering-only vs full reenactment; presented as unresolved above.
- **LAM enrollment:** 1.4 s (A100) vs "<1 s" — seconds-scale either way.
- **Beam latency:** 105.8 ms belongs to the 2021 local-GPU prototype only; production Beam (cloud reconstruction) has **no published end-to-end latency**.
- **Beam bandwidth:** 50 vs 100 Mbps minimum across partner documents; camera count 6 (Google/HP) vs 7 (SVConline specs listing).
- **Beam shipping:** "late 2025" (announcement) vs "later in 2026" (HP Imagine hands-on) — a real slip.
- **Proto M price:** $5,900 vs $6,900 across sources — unresolved; either way it is a 2D-video depth illusion, not volumetric.
- **Persona dates:** announced June 9, 2025; shipped-quality claims carry the Sept 15, 2025 release date.
- **Quest 3S "12M+ units":** unsourced; treated as an unverified estimate.
- **GaussianTalker:** 120 vs 130 FPS across sources, and two distinct 2024 papers share the name.
