# Data Model: devtool-add-remove

**Feature**: Add and Remove Devtool  
**Spec**: [spec.md](./spec.md)

## Entities

### `ToolDefinition` (Existing Concept)
Represents a statically predefined tool that the system knows how to manage.
- **Fields**:
  - `Name` (string): Unique identifier for the tool (e.g., "docker", "postgres")
  - `Type` (enum: Singleton | PortBound): Identifies if the tool can have multiple instances.
  - `DefaultPort` (integer, optional): The standard port it binds to, if applicable.
- **Constraints**: 
  - Must be statically defined in the codebase or a trusted source.
  - The `add` command verifies the requested tool against this list.

### `ManagedInstance` (Core)
Represents an actual instance of a tool added to the user's managed list.
- **Fields**:
  - `ToolName` (string): Reference to the `ToolDefinition` name.
  - `Identifier` (string): Unique identifier for this specific instance. For singleton tools, this is usually just the `ToolName`. For port-bound tools, this could be the port number or a uniquely generated alias (e.g., `postgres-5432`).
  - `CreatedAt` (timestamp): When the tool was added to the managed list.
- **Constraints**:
  - `Identifier` must be globally unique within the managed list.
  - For Singleton tools, there can only ever be *one* `ManagedInstance` with that `ToolName`.

### `ManagedConfig` (Existing/Augmented)
The root object that gets serialized to the local configuration file.
- **Fields**:
  - `Instances` (array of `ManagedInstance`): The list of currently active managed tools.
- **Operations**:
  - `AddInstance(instance *ManagedInstance) error`: Validates and adds an instance. Returns an error if it's a duplicate singleton or if the identifier conflicts.
  - `RemoveInstance(identifier string) error`: Removes the instance with the specific identifier.

## State Transitions
**Adding a Tool (`devtool add <tool>`)**
1. User provides `ToolName`.
2. System looks up `ToolDefinition`.
   - If not found: Error (Unsupported tool).
3. Check `ToolDefinition.Type`:
   - If Singleton: Check `ManagedConfig` for existing `ToolName`. If exists: Error (Already managed). If not exists: Create instance with `Identifier = ToolName`, add to config.
   - If PortBound: Attempt to use `DefaultPort`. If that identifier (`<ToolName>-<Port>`) exists in config, generate/prompt for a new port/identifier. Add new instance to config.
4. Save config.

**Removing a Tool (`devtool remove <tool>`)**
1. User provides `ToolName`.
2. System queries `ManagedConfig` for all instances matching `ToolName`.
   - If 0 found: Error (Not managed).
   - If exactly 1 found: Prompt to force stop if running. If confirmed/stopped, remove from config and save.
   - If > 1 found (PortBound): Present interactive selection list showing `Identifier`s. User selects one. Prompt to force stop if running. If confirmed/stopped, remove selected from config and save.
