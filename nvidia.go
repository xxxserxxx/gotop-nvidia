package main

// TODO: Optimization: cache most recent info w/ timestamp, and only update if older than X

import (
	"strconv"
	"strings"
	"time"

	"github.com/rai-project/nvidia-smi"
	//"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/xxxserxxx/gotop/devices"
)

func Init() {
	devices.RegisterTemp(updateNvidiaTemp)
	devices.RegisterMem(updateNvidiaMem)
	devices.RegisterCPU(updateNvidiaUsage)
}

func updateNvidiaTemp(temps map[string]int) map[string]error {
	var errs map[string]error
	info, err := nvidiasmi.New()
	if err != nil {
		errs["nvidia"] = err
	}
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
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
	var errs map[string]error
	info, err := nvidiasmi.New()
	if err != nil {
		errs["nvidia"] = err
	}
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
			name := gpu.ProductName + strconv.Itoa(i)
			mem, err := strconv.Atoi(gpu.MemoryUtil)
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
			dev := devices.MemoryInfo{
				Total: uint64(total),
				Used:  uint64(used),
			}
			dev.UsedPercent = 1.0 / float64(mem)
			mems[name] = dev
		}
	}
	return errs
}

func updateNvidiaUsage(cpus map[string]int, _ time.Duration, _ bool) map[string]error {
	errs := make(map[string]error)
	info, err := nvidiasmi.New()
	if err != nil {
		errs["nvidia"] = err
		return errs
	}
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
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
