# gospn
[Save Page Now](https://web.archive.org/save) client in Go

Set environment variable `GOSPN_DEBUG` to `1` to enable debug output.

## Usage

```go
c, _ := spn.Init("YOUR_ACCESS_KEY", "YOUR_SECRET_KEY")
defer c.Close()
options := spn.CaptureOptions{
    SkipFirstArchive:    true,
    IfNotArchivedWithin: "3d",
    ... More options ...
}
url := "https://example.com"
captureResp, err = c.Capture(url, options)
```

Some possible capture responses:

```json
{"url":"https://example.com/","job_id":"spn2-0123456789abcdef0123456789abcdef12345678"}
{"message": "Cannot resolve host nxdomain.fake.tld.", "status": "error", "status_ext": "error:invalid-host-resolution"}
{"url":"https://example.com/","job_id":null,"message":"The same snapshot had been made 3 minutes ago. You can make new capture of this URL after 2 hours."}
```

>[!NOTE]
> `Capture()` will return immediately after sending the request to the Save Page Now API. The actual capture process may take a while to complete. You can use the `GetCaptureStatus()` method to check the status of the capture job.
> 
