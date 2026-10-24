package pantech

// Ptr returns a pointer to v. Optional request fields are pointers, so a field
// left nil is not sent, and Ptr sets one inline:
//
//	pantech.CreateVolumeRequest{Name: "data", DiskOfferingSlug: "small-5gb", SizeGB: pantech.Ptr(50)}
func Ptr[T any](v T) *T { return &v }
