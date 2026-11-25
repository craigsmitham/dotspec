---
generated_by: dotspec
type: interface-spec-report
source: "../../specs/personal-task-list.system.spec.md#task-list-forms"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Interfaces](<../personal-task-list.system.spec-report.md#interfaces>)

# Task list forms

The forms the list owner uses to review and update commitments.

Open offers work the owner can advance. Waiting shows the expected reply and when to chase it. Completed holds settled errands the owner may reopen. After an accepted change the browser returns to the view where the person submitted it.

## Contracts

[Form template](<../../app/list.html>) supplies the actions and fields.

## Adopts

### Redirect after an accepted form submission

<!-- dotspec-doc {"type":"adopted-provision-report","documents":"../../specs/personal-task-list.system.spec.md#redirect-after-an-accepted-form-submission"} -->

[View section source](<../../specs/personal-task-list.system.spec.md#redirect-after-an-accepted-form-submission>)

[303 See Other](<https://www.rfc-editor.org/rfc/rfc9110.html#section-15.4.4>).

#### Edition

RFC 9110, June 2022.

#### Applicability

After accepting a change through a task form, Task service redirects the browser to the originating list view.

## Related specifications

**Applies here**

- [Show which replies need follow-up](<../requirements/show-which-replies-need-follow-up.requirement.spec-report.md>)

---

[View source](<../../specs/personal-task-list.system.spec.md#task-list-forms>)
