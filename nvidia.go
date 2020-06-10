package nvidia

import (
	"bytes"
	"encoding/csv"
	"os/exec"
	"strconv"
	"sync"
	"time"

	//"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/xxxserxxx/gotop/v4/devices"
)

func init() {
	_temps = make(map[string]int)
	_mems  = make(map[string]devices.MemoryInfo)
	_cpus  = make(map[string]int)
	errors = make(map[string]error)
	devices.RegisterTemp(updateNvidiaTemp)
	devices.RegisterMem(updateNvidiaMem)
	devices.RegisterCPU(updateNvidiaUsage)

	lock = sync.Mutex{}
	devices.RegisterStartup(startup)
}

func updateNvidiaTemp(temps map[string]int) map[string]error {
	lock.Lock()
	defer lock.Unlock()
	for k, v := range _temps {
		temps[k] = v
	}
	return errors
}

func updateNvidiaMem(mems map[string]devices.MemoryInfo) map[string]error {
	lock.Lock()
	defer lock.Unlock()
	for k, v := range _mems {
		mems[k] = v
	}
	return errors
}

func updateNvidiaUsage(cpus map[string]int, _ bool) map[string]error {
	lock.Lock()
	defer lock.Unlock()
	for k, v := range _cpus {
		cpus[k] = v
	}
	return errors
}

func startup(vars map[string]string) error {
	var err error
	refresh := time.Second
	if v, ok := vars["nvidia-refresh"]; ok {
		if refresh, err = time.ParseDuration(v); err != nil {
			return err
		}
	}
	go func() {
		timer := time.Tick(refresh)
		for range timer {
			update()
		}
	}()
	return nil
}

var (
	_temps map[string]int
	_mems  map[string]devices.MemoryInfo
	_cpus  map[string]int
	errors map[string]error
)

var lock sync.Mutex

// update updates the cached NVidia metric data: name, index,
// temperature.gpu, utilization.gpu, utilization.memory, memory.total, memory.free, memory.used
func update() {
	bs, err := exec.Command(
		"nvidia-smi",
		"--query-gpu=name,index,temperature.gpu,utilization.gpu,memory.total,memory.used",
		"--format=csv,noheader,nounits").Output()
	if err != nil {
		errors["nvidia"] = err
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(bs))
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()
	if err != nil {
		errors["nvidia"] = err
		return
	}

	lock.Lock()
	defer lock.Unlock()
	for _, row := range records {
		name := row[0] + "." + row[1]
		if _temps[name], err = strconv.Atoi(row[2]); err != nil {
			errors[name] = err
		}
		if _cpus[name], err = strconv.Atoi(row[3]); err != nil {
			errors[name] = err
		}
		t, err := strconv.Atoi(row[4])
		if err != nil {
			errors[name] = err
		}
		u, err := strconv.Atoi(row[5])
		if err != nil {
			errors[name] = err
		}
		_mems[name] = devices.MemoryInfo{
			Total:       1048576*uint64(t),
			Used:        1048576*uint64(u),
			UsedPercent: float64(u) / float64(t),
		}
	}
}
