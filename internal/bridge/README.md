# Requirements

1. Stores bridge information at file location provided.
2. Stores and retrieves application key for each bridge.
3. Consumer should store which bridge its assigned to (or if they think there's just one, get the only bridge).
4. If the key is not stored, client should initiate a call to create a key and once retrieved, persist it and return the
   bridge ID.
5. Bridges are identified by their mDNS name OR their address.
6. If the key becomes invalid (TBD) the key should be marked as stale.

Business Process

1. Provide a configuration file location or accept default.
2. Provide the bridge ID if already set and attempt to authenticate to it.
3. If not set, attempt to acquire one and then return the bridge ID.
4. Credentials will be stored in the configuration file location.

OR

1. They provide IP address and hue key.

