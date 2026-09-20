package blob

import "git.sr.ht/~rockorager/go-jmap"

func init() {
	jmap.RegisterMethod("Blob/copy", newCopyResponse)
	jmap.RegisterMethod("Blob/upload", newUploadResponse)
	jmap.RegisterMethod("Blob/get", newGetResponse)
	jmap.RegisterMethod("Blob/lookup", newLookupResponse)
}
