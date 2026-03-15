# Feature Specification: Add and Remove Devtool

**Feature Branch**: `002-devtool-add-remove`  
**Created**: 2026-03-15  
**Status**: Draft  
**Input**: User description: "devtool add and remove - Add and remove devtool from the manage list, from predefined supported tools. Which tools that singleton like telepresence or Docker, which can be added once; and for tools that binded by port like postgres or mysql, can be add more than 1 for management"

## Clarifications

### Session 2026-03-15 (Pass 2)

- Q: If a user specifies a `--port` that is already assigned to a managed instance, how should the CLI behave? → A: Fail and demand input (Option B) - Note: The user provided 'B' for both questions. For Q1, B means "Prompt to change". For Q2, B means "Fail and demand input". The updated logic will reflect this.
- Q: What happens if the system cannot find an available port automatically? → A: Fail and demand input (Option B)

### Session 2026-03-15

- Q: If a user attempts to add a predefined supported tool that is not actually installed on their host system, how should the CLI behave? → A: Just add it (Option A)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Add a Singleton Tool (Priority: P1)

As a developer, I want to add a singleton tool (like Docker or Telepresence) to my managed list so that I can easily prepare and manage my environment.

**Why this priority**: Core functionality for managing fundamental single-instance tools required for the environment.

**Independent Test**: Can be independently tested by running the add command for a singleton tool and verifying it appears exactly once in the tool list, with subsequent additions being rejected.

**Acceptance Scenarios**:

1. **Given** an empty managed list, **When** the user attempts to add a predefined singleton tool (e.g., Docker), **Then** the tool is successfully added to the managed list.
2. **Given** a singleton tool is already in the managed list, **When** the user attempts to add the same tool again, **Then** the system presents an error message indicating the tool cannot be added multiple times.

---

### User Story 2 - Add a Port-bound Tool (Priority: P1)

As a developer, I want to add a port-bound tool (like PostgreSQL or MySQL) to my managed list multiple times so that I can manage different instances for different services.

**Why this priority**: Essential for developing microservices that require separate databases or caching layers running simultaneously.

**Independent Test**: Can be validated by adding the same port-bound tool multiple times (with distinct identifiers/ports) and ensuring all instances are tracked correctly without conflict.

**Acceptance Scenarios**:

1. **Given** the user selects a predefined port-bound tool (e.g., PostgreSQL), **When** they add it for the first time, **Then** the tool is successfully added with its default configuration/port.
2. **Given** a port-bound tool is already active, **When** the user adds it again, **Then** the system prompts for or automatically assigns a new port/identifier, and successfully adds the second instance.

---

### User Story 3 - Remove a Tool (Priority: P1)

As a developer, I want to remove a tool from the managed list when I no longer need it, so that my environment remains clean and avoids wasting resources.

**Why this priority**: To allow lifecycle management of managed tools and free up ports or memory.

**Independent Test**: Can be tested by adding a tool and then removing it, validating the tool is no longer returned in the active list.

**Acceptance Scenarios**:

1. **Given** a specific instance of a tool (singleton or port-bound) exists in the managed list, **When** the user removes it by identifier, **Then** it is successfully removed from the managed list.
2. **Given** a user attempts to remove a tool that is not in the managed list, **When** they execute the remove action, **Then** an error message is presented indicating the tool is not managed.

### Edge Cases

- What happens when the user tries to add a tool that is not in the predefined supported list? (Reject with an informative error listing supported tools)
- How does system handle removing a tool that is currently running or locked by another process? Upon removal, if the tool is currently running, the system will prompt the user to confirm whether to forceful stop it before removing. If it is run non-interactively, then the default behavior is to fail and demand the user to stop the tool first.
- What happens when a user explicitly provides a `--port` that is already occupied? The system will interactively inform the user of the conflict and ask if they want to use a different port instead.
- What happens when a user adds a port-bound tool, but all available/preferred ports are occupied? The system will fail the command and ask the user to explicitly provide a `--port` flag.
- How to uniquely identify multiple instances of the same port-bound tool during removal? When removing a port-bound tool, if there are multiple instances running, the system will list all existing instances and ask the user to interactively select which to remove. If it is run non-interactively, then the command must fail and demand an explicit identifier.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a command-line interface to add a tool to the user's managed list, without verifying if the tool's executable is actually installed on the host.
- **FR-002**: System MUST provide a command-line interface to remove a tool from the user's managed list.
- **FR-003**: System MUST reject adding tools that do not exist in the predefined supported list.
- **FR-004**: System MUST differentiate between singleton tools (max 1 instance) and port-bound tools (multiple instances allowed).
- **FR-005**: System MUST prevent adding more than one instance of a singleton tool.
- **FR-006**: System MUST allow adding multiple instances of a port-bound tool.
- **FR-007**: System MUST clearly identify each managed tool instance so users can selectively remove them.

### Key Entities

- **Tool Definition**: Represents a predefined supported tool, categorized as either `singleton` or `port-bound`.
- **Managed Instance**: Represents a specific runtime allocation of a tool in the user's managed list, including its identifier (e.g., port or alias) and its parent Tool Definition.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully add a supported tool into their local managed configuration file in under 3 seconds.
- **SC-002**: System 100% reliably prevents duplication of singleton tools in the configuration.
- **SC-003**: Users are able to add 5 distinct instances of a port-bound tool without configuration collision.
- **SC-004**: Users can reliably target and remove a single specific instance of a port-bound tool without affecting other instances.
