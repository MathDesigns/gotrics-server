package hardware

type HardwareInfo struct {
	Hostname        string `json:"hostname"`
	Platform        string `json:"platform"`
	OS              string `json:"os"`
	PlatformFamily  string `json:"platform_family"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	CPUModel        string `json:"cpu_model"`
	CPUCores        int32  `json:"cpu_cores"`
	CPULogicalCores int32  `json:"cpu_logical_cores"`
	TotalMemoryGB   uint64 `json:"total_memory_gb"`
}
