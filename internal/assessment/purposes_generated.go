// Code generated from model/metamodel.yaml. DO NOT EDIT.

package assessment

const RequirementsVerification Purpose = "requirements_verification"

const RequirementsValidation Purpose = "requirements_validation"

const SystemVerification Purpose = "system_verification"

const SystemValidation Purpose = "system_validation"

const AssumptionValidation Purpose = "assumption_validation"

var Definitions = [...]Definition{
	{RequirementsVerification, "Requirements verification", "Do the requirement statements and set meet the applicable criteria for clarity, consistency and verifiability?"},
	{RequirementsValidation, "Requirements validation", "Do the proposed requirements faithfully express stakeholder needs, higher-level intent and constraints?"},
	{SystemVerification, "System verification", "Does the realization fulfill the specified obligations under the assessed conditions?"},
	{SystemValidation, "System validation", "Can the intended users accomplish their work with the proposed or realized solution in its intended context?"},
	{AssumptionValidation, "Assumption validation", "Does the evidence support the premise under the stated conditions, and is reliance on it justified for the intended use?"},
}
