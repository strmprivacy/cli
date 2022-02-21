
Print the current (JWT) access token to the terminal that can be used in a http header. Note that the token is printed on `stdout`, and the Expiry and billing-id are on
`stderr` so it’s easy to capture the token for scripting use with

```bash
export token=$(strm auth access-token)
```

Note that this token might be expired, so a refresh may be required.
Use token as follows:
'Authorization: Bearer &lt;token&gt;' 

### Usage