package r2libs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterConsole(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"log": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			timestamp := time.Now().Format("15:04:05")
			consoleWrite(fmt.Sprintf("[%s] %s\n", timestamp, joinConsoleArgs(args)))
			return nil
		}),

		"info": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			timestamp := time.Now().Format("15:04:05")
			consoleWrite(fmt.Sprintf("\033[36m[%s] INFO:\033[0m %s\n", timestamp, joinConsoleArgs(args)))
			return nil
		}),

		"warn": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			timestamp := time.Now().Format("15:04:05")
			consoleWrite(fmt.Sprintf("\033[33m[%s] WARN:\033[0m %s\n", timestamp, joinConsoleArgs(args)))
			return nil
		}),

		"error": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			timestamp := time.Now().Format("15:04:05")
			consoleWrite(fmt.Sprintf("\033[31m[%s] ERROR:\033[0m %s\n", timestamp, joinConsoleArgs(args)))
			return nil
		}),

		"debug": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			timestamp := time.Now().Format("15:04:05")
			consoleWrite(fmt.Sprintf("\033[35m[%s] DEBUG:\033[0m %s\n", timestamp, joinConsoleArgs(args)))
			return nil
		}),

		"clear": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			fmt.Print("\033[2J\033[H")
			return nil
		}),

		"group": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			label := "Group"
			if len(args) > 0 {
				label = fmt.Sprintf("%v", args[0])
			}
			fmt.Printf("\033[1m▼ %s\033[0m\n", label)
			return nil
		}),

		"groupEnd": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			fmt.Println()
			return nil
		}),

		"table": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}

			switch data := args[0].(type) {
			case []interface{}:
				printArrayTable(data)
			case map[string]interface{}:
				printObjectTable(data)
			default:
				fmt.Printf("| %v |\n", data)
			}
			return nil
		}),

		"time": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			label := "default"
			if len(args) > 0 {
				label = fmt.Sprintf("%v", args[0])
			}
			startTime := time.Now()
			setConsoleTimer(label, startTime)
			fmt.Printf("Timer '%s' started\n", label)
			return nil
		}),

		"timeEnd": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			label := "default"
			if len(args) > 0 {
				label = fmt.Sprintf("%v", args[0])
			}
			startTime := getConsoleTimer(label)
			if startTime != nil {
				elapsed := time.Since(*startTime)
				fmt.Printf("Timer '%s': %v\n", label, elapsed)
				removeConsoleTimer(label)
			} else {
				fmt.Printf("Timer '%s' not found\n", label)
			}
			return nil
		}),

		"count": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			label := "default"
			if len(args) > 0 {
				label = fmt.Sprintf("%v", args[0])
			}
			count := incrementConsoleCounter(label)
			fmt.Printf("%s: %d\n", label, count)
			return nil
		}),

		"countReset": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			label := "default"
			if len(args) > 0 {
				label = fmt.Sprintf("%v", args[0])
			}
			resetConsoleCounter(label)
			fmt.Printf("Counter '%s' reset\n", label)
			return nil
		}),

		"assert": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}

			condition := false
			if val, ok := args[0].(bool); ok {
				condition = val
			} else if val, ok := args[0].(float64); ok {
				condition = val != 0
			} else if val, ok := args[0].(string); ok {
				condition = val != ""
			}

			if !condition {
				var b strings.Builder
				b.WriteString("\033[31mAssertion failed:")
				if len(args) > 1 {
					b.WriteString(" ")
					b.WriteString(joinConsoleArgs(args[1:]))
				}
				b.WriteString("\033[0m\n")
				consoleWrite(b.String())
			}
			return nil
		}),

		"dir": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}

			obj := args[0]
			var b strings.Builder
			switch v := obj.(type) {
			case map[string]interface{}:
				b.WriteString("Object {\n")
				for key, value := range v {
					fmt.Fprintf(&b, "  %s: %v\n", key, formatValue(value))
				}
				b.WriteString("}\n")
			case []interface{}:
				b.WriteString("Array [\n")
				for i, value := range v {
					fmt.Fprintf(&b, "  %d: %v\n", i, formatValue(value))
				}
				b.WriteString("]\n")
			default:
				fmt.Fprintf(&b, "%T: %v\n", obj, obj)
			}
			consoleWrite(b.String())
			return nil
		}),

		"trace": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			consoleWrite("\033[35mTrace:" + joinConsoleArgs(args) + "\033[0m\n")
			return nil
		}),

		"prompt": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			message := "Enter input:"
			if len(args) > 0 {
				message = fmt.Sprintf("%v", args[0])
			}

			fmt.Print(message + " ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return ""
			}

			return strings.TrimSpace(input)
		}),

		"confirm": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			message := "Confirm (y/n):"
			if len(args) > 0 {
				message = fmt.Sprintf("%v", args[0])
			}

			fmt.Print(message + " ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return false
			}

			input = strings.TrimSpace(strings.ToLower(input))
			return input == "y" || input == "yes"
		}),

		"read": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			prompt := ""
			if len(args) > 0 {
				prompt = fmt.Sprintf("%v", args[0])
			}

			if prompt != "" {
				fmt.Print(prompt)
			}

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return ""
			}

			return strings.TrimSpace(input)
		}),

		"readLine": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return ""
			}

			return strings.TrimSpace(input)
		}),

		"readPassword": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			prompt := "Password:"
			if len(args) > 0 {
				prompt = fmt.Sprintf("%v", args[0])
			}

			fmt.Print(prompt + " ")

			// Simple implementation - in production, you'd want to use terminal libraries
			// for proper password hiding
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return ""
			}

			return strings.TrimSpace(input)
		}),

		"print": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			consoleWrite(joinConsoleArgs(args))
			return nil
		}),

		"println": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			consoleWrite(joinConsoleArgs(args) + "\n")
			return nil
		}),

		"printf": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}

			format, ok := args[0].(string)
			if !ok {
				return nil
			}

			if len(args) > 1 {
				fmt.Printf(format, args[1:]...)
			} else {
				fmt.Print(format)
			}
			return nil
		}),

		"getChar": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			// Simple implementation - reads one character
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()
			if err != nil {
				return ""
			}
			return string(char)
		}),

		"pause": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			message := "Press Enter to continue..."
			if len(args) > 0 {
				message = fmt.Sprintf("%v", args[0])
			}

			fmt.Print(message)
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
			return nil
		}),

		"beep": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			fmt.Print("\a")
			return nil
		}),

		"setTitle": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}
			title := fmt.Sprintf("%v", args[0])
			fmt.Printf("\033]0;%s\007", title)
			return nil
		}),

		"color": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				return nil
			}

			color := fmt.Sprintf("%v", args[0])
			text := fmt.Sprintf("%v", args[1])

			colorCode := getColorCode(color)
			if colorCode != "" {
				fmt.Printf("%s%s\033[0m", colorCode, text)
			} else {
				fmt.Print(text)
			}
			return nil
		}),

		"colorLine": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				return nil
			}

			color := fmt.Sprintf("%v", args[0])
			text := fmt.Sprintf("%v", args[1])

			colorCode := getColorCode(color)
			if colorCode != "" {
				fmt.Printf("%s%s\033[0m\n", colorCode, text)
			} else {
				fmt.Println(text)
			}
			return nil
		}),

		"bold": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}
			text := fmt.Sprintf("%v", args[0])
			fmt.Printf("\033[1m%s\033[0m", text)
			return nil
		}),

		"italic": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}
			text := fmt.Sprintf("%v", args[0])
			fmt.Printf("\033[3m%s\033[0m", text)
			return nil
		}),

		"underline": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}
			text := fmt.Sprintf("%v", args[0])
			fmt.Printf("\033[4m%s\033[0m", text)
			return nil
		}),

		"moveCursor": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				return nil
			}

			row, ok1 := args[0].(float64)
			col, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return nil
			}

			fmt.Printf("\033[%d;%dH", int(row), int(col))
			return nil
		}),

		"clearLine": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			fmt.Print("\033[2K\r")
			return nil
		}),

		"hideCursor": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			fmt.Print("\033[?25l")
			return nil
		}),

		"showCursor": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			fmt.Print("\033[?25h")
			return nil
		}),

		"progressBar": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				return nil
			}

			progress, ok := args[0].(float64)
			if !ok {
				return nil
			}

			width := 40
			if len(args) > 1 {
				if w, ok := args[1].(float64); ok {
					width = int(w)
				}
			}

			filled := int(progress * float64(width))
			if filled > width {
				filled = width
			}

			var b strings.Builder
			b.WriteString("\r[")
			for i := 0; i < width; i++ {
				if i < filled {
					b.WriteString("=")
				} else {
					b.WriteString(" ")
				}
			}
			fmt.Fprintf(&b, "] %.1f%%", progress*100)
			consoleWrite(b.String())
			return nil
		}),

		"spinner": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			spinChars := []string{"|", "/", "-", "\\"}
			step := 0
			if len(args) > 0 {
				if s, ok := args[0].(float64); ok {
					step = int(s) % len(spinChars)
					if step < 0 {
						step += len(spinChars)
					}
				}
			}

			fmt.Printf("\r%s", spinChars[step])
			return nil
		}),
	}

	RegisterModule(env, "console", functions)
}

// Global state for console functionality
var (
	consoleStateMu  sync.Mutex
	consoleTimers   = make(map[string]time.Time)
	consoleCounters = make(map[string]int)
)

// consoleWrite serializes writes to stdout so that a single log/table/etc.
// call is never torn apart by another goroutine's concurrent console output.
func consoleWrite(s string) {
	consoleStateMu.Lock()
	defer consoleStateMu.Unlock()
	fmt.Print(s)
}

func joinConsoleArgs(args []interface{}) string {
	parts := make([]string, len(args))
	for i, arg := range args {
		parts[i] = formatValue(arg)
	}
	return strings.Join(parts, " ")
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		if v == float64(int(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "null"
	case []interface{}:
		return fmt.Sprintf("[%d items]", len(v))
	case map[string]interface{}:
		return fmt.Sprintf("{%d keys}", len(v))
	case *r2core.DateValue:
		return v.Time.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func printArrayTable(data []interface{}) {
	if len(data) == 0 {
		return
	}

	var b strings.Builder
	b.WriteString("┌─────────┬─────────────────────────┐\n")
	b.WriteString("│ (index) │         Values          │\n")
	b.WriteString("├─────────┼─────────────────────────┤\n")

	for i, item := range data {
		fmt.Fprintf(&b, "│ %7d │ %-23s │\n", i, formatValue(item))
	}

	b.WriteString("└─────────┴─────────────────────────┘\n")
	consoleWrite(b.String())
}

func printObjectTable(data map[string]interface{}) {
	if len(data) == 0 {
		return
	}

	var b strings.Builder
	b.WriteString("┌─────────────────────────┬─────────────────────────┐\n")
	b.WriteString("│          Key            │         Value           │\n")
	b.WriteString("├─────────────────────────┼─────────────────────────┤\n")

	for key, value := range data {
		fmt.Fprintf(&b, "│ %-23s │ %-23s │\n", key, formatValue(value))
	}

	b.WriteString("└─────────────────────────┴─────────────────────────┘\n")
	consoleWrite(b.String())
}

func setConsoleTimer(label string, startTime time.Time) {
	consoleStateMu.Lock()
	defer consoleStateMu.Unlock()
	consoleTimers[label] = startTime
}

func getConsoleTimer(label string) *time.Time {
	consoleStateMu.Lock()
	defer consoleStateMu.Unlock()
	if startTime, exists := consoleTimers[label]; exists {
		return &startTime
	}
	return nil
}

func removeConsoleTimer(label string) {
	consoleStateMu.Lock()
	defer consoleStateMu.Unlock()
	delete(consoleTimers, label)
}

func incrementConsoleCounter(label string) int {
	consoleStateMu.Lock()
	defer consoleStateMu.Unlock()
	consoleCounters[label]++
	return consoleCounters[label]
}

func resetConsoleCounter(label string) {
	consoleStateMu.Lock()
	defer consoleStateMu.Unlock()
	consoleCounters[label] = 0
}

func getColorCode(color string) string {
	colors := map[string]string{
		"black":   "\033[30m",
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
		"reset":   "\033[0m",
	}

	if code, exists := colors[strings.ToLower(color)]; exists {
		return code
	}
	return ""
}
