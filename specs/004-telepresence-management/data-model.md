# Data Model: Telepresence Management

## Devtool Item (Telepresence)

This feature relies on the existing devtool list model. We will represent the Telepresence management integration as a `DevtoolItem`.

### Entities

1. **Telepresence Tool**
   - **ID**: `telepresence`
   - **Name**: `Telepresence`
   - **Type**: `CLI`
   - **Connection Status**: Enum / Boolean (Connected, Disconnected)
   - **Actions**:
     - `Connect`
     - `Disconnect`
     - `Restart`
     - `Remove`

### UI State (Bubbletea)

- **telepresenceInstalled**: `bool` - Tracks whether the binary is available.
- **telepresenceStatus**: `string` - Records the current stdout output of `telepresence status`.
- **isConnecting**: `bool` - Tracks if a connection attempt is currently in progress.
- **connectionError**: `error` - Any errors encountered during the connection process.

## State Transitions

- **Disconnected -> Connect Initiated**: User selects Connect action.
- **Connect Initiated -> Connected**: `telepresence connect` finishes successfully.
- **Connect Initiated -> Error**: `telepresence connect` fails or times out.
- **Connected -> Disconnect Initiated**: User selects Disconnect action.
- **Disconnect Initiated -> Disconnected**: `telepresence quit` finishes successfully.
