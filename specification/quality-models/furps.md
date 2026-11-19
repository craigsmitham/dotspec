---
generated_by: dotspec
type: quality-model-spec-report
documents: "../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#furps-starter"
---
# FURPS starter

[dotspec](<../../SPEC.md>)

## Description

A dotspec starter covering the five FURPS categories and all 32 subcategory placements (31 distinct subcategories) listed in the cited Rational Unified Process guidance. This is a documented RUP variant, not a claim that FURPS has one universal fixed taxonomy. Explanations are original authoring prompts. Functionality can classify general requirements as well as quality requirements; these categories do not impose targets. The separate FURPS+ constraint categories are outside this starter. Accuracy appears under two categories in the source and is represented once, organized under both reliability and performance.

## Edition

RUP category variant (1987–2001 guidance), dotspec starter 2

## Sources

[Rational Unified Process: Requirements](<https://galileo.softlab.ntua.gr/rup/process/workflow/requirem/co_req.htm>)

[IBM Rational FURPS overview](<https://public.dhe.ibm.com/software/it/pdf/TestingAndQuality.pdf>)

## Quality attributes

### Functionality

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#functionality"} -->

**ID:** functionality

Consider the capabilities and results people need from the system.

### Feature sets

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#feature-sets"} -->

**ID:** feature-sets\
**Sub-attribute of:** functionality

Which groups of functions are needed to support the intended work?

### Capabilities

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#capabilities"} -->

**ID:** capabilities\
**Sub-attribute of:** functionality

What actions and results must the system make possible?

### Security

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#security"} -->

**ID:** security\
**Sub-attribute of:** functionality

What protection and access controls must the system provide?

### Usability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#usability"} -->

**ID:** usability

Consider how people discover, learn, and use the system, including its documentation.

### Human factors

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#human-factors"} -->

**ID:** human-factors\
**Sub-attribute of:** usability

How should the interaction accommodate the people and work environments involved?

### Aesthetics

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#aesthetics"} -->

**ID:** aesthetics\
**Sub-attribute of:** usability

What expectations apply to the visual and sensory presentation?

### User-interface consistency

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#user-interface-consistency"} -->

**ID:** user-interface-consistency\
**Sub-attribute of:** usability

Which interaction conventions should behave consistently across the product?

### Online and context-sensitive help

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#online-and-context-sensitive-help"} -->

**ID:** online-and-context-sensitive-help\
**Sub-attribute of:** usability

What help must be available during a task and tailored to the current situation?

### Wizards and agents

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#wizards-and-agents"} -->

**ID:** wizards-and-agents\
**Sub-attribute of:** usability

Which tasks need guided steps or delegated assistance?

### User documentation

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#user-documentation"} -->

**ID:** user-documentation\
**Sub-attribute of:** usability

What instructions and reference information do users need?

### Training materials

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#training-materials"} -->

**ID:** training-materials\
**Sub-attribute of:** usability

What learning resources are needed for users to become effective?

### Reliability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#reliability"} -->

**ID:** reliability

Consider consistent behavior, availability, and recovery when things go wrong.

### Frequency and severity of failure

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#frequency-and-severity-of-failure"} -->

**ID:** failure-frequency-and-severity\
**Sub-attribute of:** reliability

How often may failures occur, and which consequences distinguish their severity?

### Recoverability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#recoverability"} -->

**ID:** recoverability\
**Sub-attribute of:** reliability

What must be restored after failure, and how will successful recovery be demonstrated?

### Accuracy

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#accuracy"} -->

**ID:** accuracy\
**Sub-attribute of:** reliability, performance

What accuracy must results retain during continued operation and under the specified workload?

### Predictability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#predictability"} -->

**ID:** predictability\
**Sub-attribute of:** reliability

Which behaviors must remain consistent across repeated operations and changing conditions?

### Mean time between failures

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#mean-time-between-failures"} -->

**ID:** mean-time-between-failures\
**Sub-attribute of:** reliability

What observation period and failure definition will be used to assess intervals between failures?

### Performance

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#performance"} -->

**ID:** performance

Consider response times, throughput, and resource use for realistic workloads.

### Speed

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#speed"} -->

**ID:** speed\
**Sub-attribute of:** performance

How quickly must the system complete its work under the expected load?

### Efficiency

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#efficiency"} -->

**ID:** efficiency\
**Sub-attribute of:** performance

What useful output is expected for the resources consumed?

### Availability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#availability"} -->

**ID:** availability\
**Sub-attribute of:** performance

What readiness for use must be maintained under the specified workload?

### Throughput

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#throughput"} -->

**ID:** throughput\
**Sub-attribute of:** performance

How much work must the system complete per unit of time?

### Response time

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#response-time"} -->

**ID:** response-time\
**Sub-attribute of:** performance

How long may a caller wait for a response?

### Recovery time

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#recovery-time"} -->

**ID:** recovery-time\
**Sub-attribute of:** performance

How quickly must useful operation resume following a disruption?

### Resource usage

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#resource-usage"} -->

**ID:** resource-usage\
**Sub-attribute of:** performance

What resource budgets apply while the system is operating?

### Supportability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#supportability"} -->

**ID:** supportability

Consider how the system will be tested, configured, maintained, and supported throughout its life.

### Testability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#testability"} -->

**ID:** testability\
**Sub-attribute of:** supportability

What controls and observations are needed to test the system?

### Extensibility

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#extensibility"} -->

**ID:** extensibility\
**Sub-attribute of:** supportability

What additions should be possible without extensive rework?

### Adaptability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#adaptability"} -->

**ID:** adaptability\
**Sub-attribute of:** supportability

Which changes in environment or use should the system accommodate?

### Maintainability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#maintainability"} -->

**ID:** maintainability\
**Sub-attribute of:** supportability

What effort and confidence are expected when diagnosing and changing the system?

### Compatibility

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#compatibility"} -->

**ID:** compatibility\
**Sub-attribute of:** supportability

Which existing products and versions must continue to work with the system?

### Configurability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#configurability"} -->

**ID:** configurability\
**Sub-attribute of:** supportability

What behavior must operators or users be able to configure?

### Serviceability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#serviceability"} -->

**ID:** serviceability\
**Sub-attribute of:** supportability

What support and repair tasks must be practical during the product life?

### Installability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#installability"} -->

**ID:** installability\
**Sub-attribute of:** supportability

What is required to install, update, and remove the system successfully?

### Localizability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#localizability"} -->

**ID:** localizability\
**Sub-attribute of:** supportability

Which languages and regional conventions must be supported without redesign?

[Edit source](<../../internal/qualitymodel/specs/quality-models/furps.quality-model.spec.md#furps-starter>)
