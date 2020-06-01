package nvidia

// TODO: Optimization: cache most recent info w/ timestamp, and only update if older than X

import (
	"strconv"
	"strings"

	"github.com/rai-project/nvidia-smi"
	//"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/xxxserxxx/gotop/v4/devices"
)

func init() {
	devices.RegisterTemp(updateNvidiaTemp)
	devices.RegisterMem(updateNvidiaMem)
	devices.RegisterCPU(updateNvidiaUsage)
}

func updateNvidiaTemp(temps map[string]int) map[string]error {
	errs := make(map[string]error)
	info, err := nvidiasmi.New()
	if err != nil {
		errs["nvidia"] = err
		return errs
	}
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
			if gpu.GpuTemp == "N/A" {
				// The GPU does not export a temperature measure
				continue
			}
			name := gpu.ProductName + " " + strconv.Itoa(i)
			temperature, err := strconv.ParseFloat(strings.ReplaceAll(gpu.GpuTemp, " C", ""), 10)
			if err != nil {
				errs[name] = err
				continue
			}
			temps[name] = int(temperature)
		}
	}
	return errs
}

func updateNvidiaMem(mems map[string]devices.MemoryInfo) map[string]error {
	errs := make(map[string]error)
	info, err := nvidiasmi.New()
	if err != nil {
		errs["nvidia"] = err
		return errs
	}
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
			if gpu.MemoryUtil == "N/A" || gpu.Total == "N/A" || gpu.Used == "N/A" {
				// The GPU does not export sufficient memory measures
				continue
			}
			name := gpu.ProductName + strconv.Itoa(i)
			mem, err := strconv.Atoi(gpu.MemoryUtil)
			if err != nil {
				errs[name+"Mem"] = err
				continue
			}
			total, err := strconv.Atoi(gpu.Total)
			if err != nil {
				errs[name+"Total"] = err
				continue
			}
			used, err := strconv.Atoi(gpu.Used)
			if err != nil {
				errs[name+"Used"] = err
				continue
			}
			if total == 0 && used == 0 {
				total = 100
				used = mem
			} else if total != 0 && used == 0 {
				used = int(float64(total) * (float64(mem) / 100))
			} else if total == 0 && used != 0 {
				total = int(float64(used) / (float64(mem) / 100))
			}
			dev := devices.MemoryInfo{
				Total: uint64(total),
				Used:  uint64(used),
			}
			dev.UsedPercent = float64(mem)
			mems[name] = dev
		}
	}
	return errs
}

func updateNvidiaUsage(cpus map[string]int, _ bool) map[string]error {
	errs := make(map[string]error)
	info, err := nvidiasmi.New()
	if err != nil {
		errs["nvidia"] = err
		return errs
	}
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
			if gpu.GpuUtil == "N/A" {
				// The GPU does not export sufficient memory measures
				continue
			}
			name := gpu.ProductName + " " + strconv.Itoa(i)
			usage, err := strconv.Atoi(gpu.GpuUtil)
			if err != nil {
				errs[name] = err
				continue
			}
			cpus[name] = usage
		}
	}
	return errs
}
