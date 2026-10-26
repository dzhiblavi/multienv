# JSON storage

The server exposes `POST /storage`. Requests use the same action/request envelope as other operations:

```json
{
  "action": "set",
  "request": {
    "key": "example",
    "path": ["devbox", "status"],
    "value": {"ready": true}
  }
}
```

`key` selects a Redis JSON document. `path` is a list of literal object keys;
omit it or use `[]` to select the whole document. Dots, quotes and numeric strings
in path segments are object keys, not JSONPath expressions or array indices.

| Action | Request fields | Response |
| --- | --- | --- |
| `get` | `key`, optional `path` | `{"found": true, "value": ...}` |
| `set` | `key`, optional `path`, `value` | `{}` |
| `set-if-absent` | `key`, optional `path`, `value` | `{"created": true}` or `{"created": false}` |
| `delete` | `key`, optional `path` | `{"deleted": true}` or `{"deleted": false}` |

A missing value returns `{"found": false, "value": null}`. Stored JSON `null` returns
`{"found": true, "value": null}` and counts as existing for `set-if-absent`.

Writes atomically create missing parent objects. An existing parent that is not an object
causes an error without changing the document. Arrays can be stored as values;
individual array elements cannot be selected through `path`.

Updates to different branches preserve each other. Writes to the same location replace its value
in execution order; storage does not compare snapshot timestamps.
