# Propuesta: CLI de Profiling para R2Lang

**Versión:** 1.0  
**Fecha:** 2025-07-15  
**Estado:** Propuesta  

## Resumen Ejecutivo

Esta propuesta presenta un sistema integral de profiling para R2Lang que permite analizar el rendimiento de código R2 a través de comandos CLI especializados, incluyendo profiling de memoria, CPU, tiempo de ejecución y análisis de cuellos de botella.

## Problema Actual

R2Lang actualmente no tiene herramientas integradas para:

- **Análisis de rendimiento de código**
- **Profiling de uso de memoria**
- **Análisis de tiempo de CPU**
- **Detección de cuellos de botella**
- **Métricas de llamadas a funciones**
- **Análisis de allocaciones**

## Solución Propuesta

### 1. CLI de Profiling

#### 1.1 Comandos Principales
```bash
# Profiling de memoria
r2 profile mem script.r2
r2 profile memory script.r2 --output memory.html
r2 profile mem script.r2 --threshold 1MB --live

# Profiling de CPU
r2 profile cpu script.r2
r2 profile cpu script.r2 --duration 30s --output cpu.prof
r2 profile cpu script.r2 --samples 1000 --flame-graph

# Profiling combinado
r2 profile all script.r2
r2 profile full script.r2 --output report.html

# Profiling en tiempo real
r2 profile live script.r2
r2 profile watch script.r2 --interval 1s
```

#### 1.2 Opciones Avanzadas
```bash
# Filtering y sampling
r2 profile cpu script.r2 --filter "UserService.*" --rate 100
r2 profile mem script.r2 --ignore-small --min-size 1KB

# Output formats
r2 profile cpu script.r2 --format json
r2 profile mem script.r2 --format pprof
r2 profile all script.r2 --format html --template custom.html

# Comparación
r2 profile compare baseline.prof current.prof
r2 profile diff old_script.r2 new_script.r2
```

### 2. Sistema de Profiling de Memoria

#### 2.1 Memory Profiler
```go
// pkg/r2core/memory_profiler.go
type MemoryProfiler struct {
    Enabled         bool
    SampleRate      int
    Threshold       int64
    Allocations     map[string]*AllocationSite
    LiveObjects     map[uintptr]*ObjectInfo
    TotalAllocated  int64
    TotalFreed      int64
    PeakMemory      int64
    GCStats         *GCStats
}

type AllocationSite struct {
    Location     string
    Function     string
    LineNumber   int
    Count        int64
    TotalSize    int64
    AverageSize  float64
    Stack        []string
}

type ObjectInfo struct {
    Type         string
    Size         int64
    AllocatedAt  time.Time
    Stack        []string
    StillAlive   bool
}

func (mp *MemoryProfiler) RecordAllocation(objType string, size int64, location string) {
    if !mp.Enabled {
        return
    }
    
    site := mp.Allocations[location]
    if site == nil {
        site = &AllocationSite{
            Location:   location,
            Function:   extractFunction(location),
            LineNumber: extractLineNumber(location),
            Stack:      getCurrentStack(),
        }
        mp.Allocations[location] = site
    }
    
    site.Count++
    site.TotalSize += size
    site.AverageSize = float64(site.TotalSize) / float64(site.Count)
    
    mp.TotalAllocated += size
    if mp.TotalAllocated-mp.TotalFreed > mp.PeakMemory {
        mp.PeakMemory = mp.TotalAllocated - mp.TotalFreed
    }
}
```

#### 2.2 Integración con R2Lang Runtime
```go
// Instrumentar allocaciones en tipos R2Lang
func (env *Environment) AllocateObject(objType string, size int64) interface{} {
    if profiler := env.GetMemoryProfiler(); profiler != nil {
        profiler.RecordAllocation(objType, size, env.GetCurrentLocation())
    }
    
    // Allocación normal
    return allocateObject(objType, size)
}

// Instrumentar arrays
func (al *ArrayLiteral) Eval(env *Environment) interface{} {
    elements := make([]interface{}, len(al.Elements))
    
    for i, elem := range al.Elements {
        elements[i] = elem.Eval(env)
    }
    
    // Record allocation
    if profiler := env.GetMemoryProfiler(); profiler != nil {
        estimatedSize := int64(len(elements) * 8) // 8 bytes per pointer
        profiler.RecordAllocation("Array", estimatedSize, al.GetLocation())
    }
    
    return elements
}
```

### 3. Sistema de Profiling de CPU

#### 3.1 CPU Profiler
```go
// pkg/r2core/cpu_profiler.go
type CPUProfiler struct {
    Enabled       bool
    SampleRate    time.Duration
    Samples       []*Sample
    FunctionTimes map[string]*FunctionStats
    StartTime     time.Time
    TotalSamples  int64
}

type Sample struct {
    Timestamp time.Time
    Stack     []string
    Goroutine int
    Duration  time.Duration
}

type FunctionStats struct {
    Name           string
    TotalTime      time.Duration
    SelfTime       time.Duration
    CallCount      int64
    AverageTime    time.Duration
    MaxTime        time.Duration
    MinTime        time.Duration
    Callers        map[string]int64
    Callees        map[string]int64
}

func (cp *CPUProfiler) StartSampling() {
    cp.StartTime = time.Now()
    cp.Enabled = true
    
    go func() {
        ticker := time.NewTicker(cp.SampleRate)
        defer ticker.Stop()
        
        for cp.Enabled {
            select {
            case <-ticker.C:
                cp.takeSample()
            }
        }
    }()
}

func (cp *CPUProfiler) takeSample() {
    stack := getCurrentStack()
    sample := &Sample{
        Timestamp: time.Now(),
        Stack:     stack,
        Goroutine: getGoroutineID(),
        Duration:  cp.SampleRate,
    }
    
    cp.Samples = append(cp.Samples, sample)
    cp.TotalSamples++
    
    // Update function statistics
    cp.updateFunctionStats(stack)
}
```

#### 3.2 Function Call Tracing
```go
// Instrumentar llamadas a funciones
func (fc *FunctionCall) Eval(env *Environment) interface{} {
    profiler := env.GetCPUProfiler()
    
    var start time.Time
    if profiler != nil && profiler.Enabled {
        start = time.Now()
        profiler.EnterFunction(fc.FunctionName, fc.GetLocation())
    }
    
    defer func() {
        if profiler != nil && profiler.Enabled {
            duration := time.Since(start)
            profiler.ExitFunction(fc.FunctionName, duration)
        }
    }()
    
    // Ejecución normal de la función
    return fc.executeFunction(env)
}
```

### 4. Generación de Reportes

#### 4.1 Memory Reports
```go
// pkg/r2core/memory_reporter.go
type MemoryReporter struct {
    Profiler *MemoryProfiler
    Options  *ReportOptions
}

type ReportOptions struct {
    Format      string // "html", "json", "text", "pprof"
    OutputFile  string
    Threshold   int64
    TopN        int
    ShowStacks  bool
    GroupBy     string // "function", "type", "location"
}

func (mr *MemoryReporter) GenerateReport() error {
    switch mr.Options.Format {
    case "html":
        return mr.generateHTMLReport()
    case "json":
        return mr.generateJSONReport()
    case "text":
        return mr.generateTextReport()
    case "pprof":
        return mr.generatePProfReport()
    default:
        return fmt.Errorf("unsupported format: %s", mr.Options.Format)
    }
}

func (mr *MemoryReporter) generateHTMLReport() error {
    data := struct {
        TotalAllocated  int64
        TotalFreed      int64
        PeakMemory      int64
        TopAllocations  []*AllocationSite
        Timeline        []MemorySnapshot
        HeapMap         map[string]int64
    }{
        TotalAllocated: mr.Profiler.TotalAllocated,
        TotalFreed:     mr.Profiler.TotalFreed,
        PeakMemory:     mr.Profiler.PeakMemory,
        TopAllocations: mr.getTopAllocations(),
        Timeline:       mr.getMemoryTimeline(),
        HeapMap:        mr.getHeapMap(),
    }
    
    return mr.renderTemplate("memory_report.html", data)
}
```

#### 4.2 CPU Reports
```go
// pkg/r2core/cpu_reporter.go
type CPUReporter struct {
    Profiler *CPUProfiler
    Options  *ReportOptions
}

func (cr *CPUReporter) GenerateFlameGraph() error {
    // Generar datos para flame graph
    flameData := cr.buildFlameGraphData()
    
    if cr.Options.Format == "svg" {
        return cr.generateSVGFlameGraph(flameData)
    }
    
    return cr.generateInteractiveFlameGraph(flameData)
}

func (cr *CPUReporter) generateTextReport() error {
    report := &strings.Builder{}
    
    report.WriteString("CPU Profiling Report\n")
    report.WriteString("===================\n\n")
    
    report.WriteString(fmt.Sprintf("Total samples: %d\n", cr.Profiler.TotalSamples))
    report.WriteString(fmt.Sprintf("Duration: %v\n", time.Since(cr.Profiler.StartTime)))
    report.WriteString(fmt.Sprintf("Sample rate: %v\n\n", cr.Profiler.SampleRate))
    
    // Top functions by total time
    report.WriteString("Top Functions by Total Time:\n")
    topFunctions := cr.getTopFunctions()
    for i, fn := range topFunctions {
        percentage := float64(fn.TotalTime) / float64(cr.getTotalTime()) * 100
        report.WriteString(fmt.Sprintf("%d. %s: %v (%.2f%%)\n", 
            i+1, fn.Name, fn.TotalTime, percentage))
    }
    
    return ioutil.WriteFile(cr.Options.OutputFile, []byte(report.String()), 0644)
}
```

### 5. Integración con CLI

#### 5.1 Memory Profiling Commands
```go
// cmd/profile_mem.go
var memProfileCmd = &cobra.Command{
    Use:   "mem [script.r2]",
    Short: "Profile memory usage of R2Lang script",
    Long:  `Profile memory allocations, deallocations, and memory usage patterns`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("Error: script file required")
            os.Exit(1)
        }
        
        // Parse flags
        outputFile, _ := cmd.Flags().GetString("output")
        format, _ := cmd.Flags().GetString("format")
        threshold, _ := cmd.Flags().GetInt64("threshold")
        live, _ := cmd.Flags().GetBool("live")
        
        // Setup profiler
        profiler := r2core.NewMemoryProfiler()
        profiler.Enabled = true
        profiler.Threshold = threshold
        
        // Run script with profiling
        env := r2core.NewEnvironment()
        env.SetMemoryProfiler(profiler)
        
        if live {
            runLiveMemoryProfiling(args[0], env)
        } else {
            runMemoryProfiling(args[0], env, outputFile, format)
        }
    },
}

func runMemoryProfiling(scriptFile string, env *r2core.Environment, outputFile, format string) {
    // Execute script
    source, err := ioutil.ReadFile(scriptFile)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }
    
    parser := r2core.NewParser(string(source))
    program := parser.ParseProgram()
    
    start := time.Now()
    program.Eval(env)
    duration := time.Since(start)
    
    // Generate report
    profiler := env.GetMemoryProfiler()
    reporter := &r2core.MemoryReporter{
        Profiler: profiler,
        Options: &r2core.ReportOptions{
            Format:     format,
            OutputFile: outputFile,
        },
    }
    
    fmt.Printf("Memory profiling completed in %v\n", duration)
    fmt.Printf("Total allocated: %d bytes\n", profiler.TotalAllocated)
    fmt.Printf("Peak memory: %d bytes\n", profiler.PeakMemory)
    
    if err := reporter.GenerateReport(); err != nil {
        fmt.Printf("Error generating report: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Report saved to: %s\n", outputFile)
}
```

#### 5.2 CPU Profiling Commands
```go
// cmd/profile_cpu.go
var cpuProfileCmd = &cobra.Command{
    Use:   "cpu [script.r2]",
    Short: "Profile CPU usage of R2Lang script",
    Long:  `Profile CPU time, function calls, and performance bottlenecks`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("Error: script file required")
            os.Exit(1)
        }
        
        // Parse flags
        outputFile, _ := cmd.Flags().GetString("output")
        format, _ := cmd.Flags().GetString("format")
        duration, _ := cmd.Flags().GetDuration("duration")
        samples, _ := cmd.Flags().GetInt("samples")
        flameGraph, _ := cmd.Flags().GetBool("flame-graph")
        
        // Setup profiler
        profiler := r2core.NewCPUProfiler()
        profiler.Enabled = true
        profiler.SampleRate = time.Millisecond * 10 // 10ms sample rate
        
        // Run script with profiling
        env := r2core.NewEnvironment()
        env.SetCPUProfiler(profiler)
        
        if flameGraph {
            runCPUProfilingWithFlameGraph(args[0], env, outputFile)
        } else {
            runCPUProfiling(args[0], env, outputFile, format, duration)
        }
    },
}
```

### 6. Live Profiling

#### 6.1 Real-time Memory Monitoring
```go
// pkg/r2core/live_profiler.go
type LiveProfiler struct {
    MemoryProfiler *MemoryProfiler
    CPUProfiler    *CPUProfiler
    UpdateInterval time.Duration
    Dashboard      *Dashboard
}

type Dashboard struct {
    MemoryChart    *Chart
    CPUChart       *Chart
    FunctionList   *FunctionList
    AllocationList *AllocationList
}

func (lp *LiveProfiler) StartLiveMonitoring(scriptFile string) {
    // Start web dashboard
    go lp.startWebDashboard()
    
    // Start profiling
    lp.MemoryProfiler.Enabled = true
    lp.CPUProfiler.StartSampling()
    
    // Update dashboard periodically
    ticker := time.NewTicker(lp.UpdateInterval)
    go func() {
        for range ticker.C {
            lp.updateDashboard()
        }
    }()
    
    // Execute script
    lp.executeScript(scriptFile)
}

func (lp *LiveProfiler) startWebDashboard() {
    http.HandleFunc("/", lp.handleDashboard)
    http.HandleFunc("/api/memory", lp.handleMemoryAPI)
    http.HandleFunc("/api/cpu", lp.handleCPUAPI)
    http.HandleFunc("/api/functions", lp.handleFunctionsAPI)
    
    fmt.Println("Live profiling dashboard: http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### 6.2 Interactive Dashboard
```html
<!-- templates/dashboard.html -->
<!DOCTYPE html>
<html>
<head>
    <title>R2Lang Live Profiler</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        .container { display: flex; flex-wrap: wrap; }
        .chart-container { width: 48%; margin: 1%; }
        .stats { width: 100%; background: #f5f5f5; padding: 10px; }
        .function-list { width: 48%; margin: 1%; }
    </style>
</head>
<body>
    <div class="container">
        <div class="stats">
            <h2>R2Lang Live Profiler</h2>
            <div id="stats">
                <span>Memory: <span id="memory-usage">-</span></span>
                <span>CPU: <span id="cpu-usage">-</span></span>
                <span>Functions: <span id="function-count">-</span></span>
            </div>
        </div>
        
        <div class="chart-container">
            <canvas id="memory-chart"></canvas>
        </div>
        
        <div class="chart-container">
            <canvas id="cpu-chart"></canvas>
        </div>
        
        <div class="function-list">
            <h3>Top Functions</h3>
            <table id="functions-table">
                <thead>
                    <tr>
                        <th>Function</th>
                        <th>Calls</th>
                        <th>Total Time</th>
                        <th>Avg Time</th>
                    </tr>
                </thead>
                <tbody id="functions-body">
                </tbody>
            </table>
        </div>
    </div>
    
    <script>
        // Initialize charts and real-time updates
        const memoryChart = new Chart(document.getElementById('memory-chart'), {
            type: 'line',
            data: { labels: [], datasets: [{ label: 'Memory Usage', data: [] }] },
            options: { responsive: true, animation: false }
        });
        
        const cpuChart = new Chart(document.getElementById('cpu-chart'), {
            type: 'line',
            data: { labels: [], datasets: [{ label: 'CPU Usage', data: [] }] },
            options: { responsive: true, animation: false }
        });
        
        // Update every second
        setInterval(updateDashboard, 1000);
        
        function updateDashboard() {
            fetch('/api/memory')
                .then(response => response.json())
                .then(data => updateMemoryChart(data));
                
            fetch('/api/cpu')
                .then(response => response.json())
                .then(data => updateCPUChart(data));
                
            fetch('/api/functions')
                .then(response => response.json())
                .then(data => updateFunctionTable(data));
        }
    </script>
</body>
</html>
```

### 7. Análisis Comparativo

#### 7.1 Diff Profiling
```bash
# Comparar dos versiones de script
r2 profile compare baseline.r2 optimized.r2

# Comparar perfiles guardados
r2 profile diff baseline.prof current.prof --format html

# Análisis de regresión
r2 profile regression old_script.r2 new_script.r2 --threshold 10%
```

#### 7.2 Implementación de Comparación
```go
// pkg/r2core/profile_comparator.go
type ProfileComparator struct {
    BaselineProfile *Profile
    CurrentProfile  *Profile
    Options         *CompareOptions
}

type CompareOptions struct {
    Threshold     float64 // Percentage threshold for reporting differences
    MetricType    string  // "memory", "cpu", "all"
    OutputFormat  string
    ShowOnlyDiffs bool
}

type ProfileDiff struct {
    Function        string
    BaselineValue   float64
    CurrentValue    float64
    AbsoluteDiff    float64
    PercentageDiff  float64
    Significance    string // "improved", "regressed", "unchanged"
}

func (pc *ProfileComparator) Compare() (*ComparisonReport, error) {
    report := &ComparisonReport{
        Summary: pc.generateSummary(),
        Diffs:   make([]*ProfileDiff, 0),
    }
    
    // Compare memory metrics
    if pc.Options.MetricType == "memory" || pc.Options.MetricType == "all" {
        memoryDiffs := pc.compareMemoryMetrics()
        report.Diffs = append(report.Diffs, memoryDiffs...)
    }
    
    // Compare CPU metrics
    if pc.Options.MetricType == "cpu" || pc.Options.MetricType == "all" {
        cpuDiffs := pc.compareCPUMetrics()
        report.Diffs = append(report.Diffs, cpuDiffs...)
    }
    
    // Filter by threshold
    if pc.Options.Threshold > 0 {
        report.Diffs = pc.filterByThreshold(report.Diffs)
    }
    
    return report, nil
}
```

### 8. Integración con CI/CD

#### 8.1 Performance Regression Detection
```yaml
# .github/workflows/performance.yml
name: Performance Testing
on: [push, pull_request]

jobs:
  performance-test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Setup R2Lang
      run: |
        # Install R2Lang
        
    - name: Run performance tests
      run: |
        # Run baseline
        r2 profile cpu benchmark.r2 --output baseline.prof
        
        # Run current
        r2 profile cpu benchmark.r2 --output current.prof
        
        # Compare
        r2 profile compare baseline.prof current.prof --threshold 10% --format json > diff.json
        
    - name: Check for regressions
      run: |
        # Fail if significant regressions detected
        if grep -q "regressed" diff.json; then
          echo "Performance regression detected!"
          exit 1
        fi
```

#### 8.2 Automated Reporting
```go
// cmd/profile_ci.go
var ciProfileCmd = &cobra.Command{
    Use:   "ci [script.r2]",
    Short: "Run profiling in CI mode",
    Long:  `Optimized profiling for CI/CD pipelines`,
    Run: func(cmd *cobra.Command, args []string) {
        // CI-optimized profiling
        // - Shorter duration
        // - JSON output
        // - Exit codes for regressions
        // - Minimal output
        
        baseline, _ := cmd.Flags().GetString("baseline")
        threshold, _ := cmd.Flags().GetFloat64("threshold")
        
        runCIProfile(args[0], baseline, threshold)
    },
}
```

## Plan de Implementación

### Fase 1: Core Profiling Infrastructure
- [ ] Memory profiler básico
- [ ] CPU profiler con sampling
- [ ] CLI commands básicos (`r2 profile mem`, `r2 profile cpu`)
- [ ] Text report generation

### Fase 2: Advanced Profiling
- [ ] Function call tracing
- [ ] Allocation tracking
- [ ] Heap analysis
- [ ] Multiple output formats (HTML, JSON, pprof)

### Fase 3: Real-time Profiling
- [ ] Live monitoring dashboard
- [ ] Web interface
- [ ] Real-time charts
- [ ] Interactive exploration

### Fase 4: Analysis and Comparison
- [ ] Profile comparison tools
- [ ] Regression detection
- [ ] Flame graph generation
- [ ] Performance optimization suggestions

### Fase 5: CI/CD Integration
- [ ] Automated performance testing
- [ ] Regression detection in CI
- [ ] Performance budgets
- [ ] Integration with popular CI systems

## Beneficios

1. **Optimización de Performance:** Identificar cuellos de botella fácilmente
2. **Debugging Avanzado:** Entender el comportamiento del código en runtime
3. **Monitoreo Continuo:** Detectar regresiones de performance
4. **Desarrollo Guiado:** Tomar decisiones basadas en datos reales
5. **Escalabilidad:** Preparar aplicaciones para producción

## Conclusión

Este sistema de profiling proporcionará a R2Lang capacidades de análisis de performance de nivel profesional, permitiendo a los desarrolladores optimizar código, detectar problemas de rendimiento y mantener aplicaciones escalables y eficientes.