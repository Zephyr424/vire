package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var vars = make(map[string]interface{})
var handles = make(map[int]interface{})
var nextHandle = 1

func trim(s string) string {
	return strings.TrimSpace(s)
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isSyscall(s string) bool {
	return strings.HasPrefix(s, "syscall(")
}

func extractSyscallArgs(expr string) []interface{} {
	start := strings.Index(expr, "(")
	end := strings.LastIndex(expr, ")")
	if start == -1 || end == -1 || start >= end {
		return nil
	}
	params := expr[start+1 : end]
	parts := strings.Split(params, ",")
	args := []interface{}{}
	for _, p := range parts {
		p = trim(p)
		if p == "" {
			continue
		}
		if p[0] == '"' {
			args = append(args, p[1:len(p)-1])
		} else if isNumber(p) {
			val, _ := strconv.Atoi(p)
			args = append(args, val)
		} else {
			if val, ok := vars[p]; ok {
				args = append(args, val)
			} else {
				args = append(args, nil)
			}
		}
	}
	return args
}

func syscallDispatch(args []interface{}) interface{} {
	if len(args) < 1 {
		return nil
	}
	num, ok := args[0].(int)
	if !ok {
		return nil
	}

	switch num {
	case 0:
		return "Windows"
	case 1:
		return os.Getpid()
	case 101:
		path, _ := args[1].(string)
		mode, _ := args[2].(string)
		var file *os.File
		var err error
		if mode == "r" {
			file, err = os.OpenFile(path, os.O_RDONLY, 0644)
		} else {
			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		}
		if err != nil {
			return nil
		}
		handles[nextHandle] = file
		handle := nextHandle
		nextHandle++
		return handle
	case 102:
		handle, ok := args[1].(int)
		if !ok {
			return nil
		}
		file, ok := handles[handle].(*os.File)
		if !ok {
			return nil
		}
		file.Seek(0, 0)
		data, err := os.ReadFile(file.Name())
		if err != nil {
			return nil
		}
		return string(data)
	case 103:
		handle, ok := args[1].(int)
		if !ok {
			return nil
		}
		file, ok := handles[handle].(*os.File)
		if !ok {
			return nil
		}
		data, _ := args[2].(string)
		file.WriteString(data)
		return 1
	case 104:
		handle, ok := args[1].(int)
		if !ok {
			return nil
		}
		if file, ok := handles[handle].(*os.File); ok {
			file.Close()
			delete(handles, handle)
		}
		return 0
	case 105:
		path, _ := args[1].(string)
		_, err := os.Stat(path)
		return err == nil
	case 106:
		path, _ := args[1].(string)
		err := os.Remove(path)
		return err == nil
	case 200:
		jsonStr, _ := args[1].(string)
		key, _ := args[2].(string)
		var data map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			return nil
		}
		if val, ok := data[key]; ok {
			return val
		}
		return nil
	case 301:
		return "HTTP GET placeholder"
	case 401:
		return "sha256 placeholder"
	case 501:
		return time.Now().Unix()
	case 502:
		return time.Now().Format("Mon Jan 2 15:04:05 2006")
	case 601:
		return "gzip placeholder"
	case 602:
		return "gunzip placeholder"
	case 701:
		return "regex placeholder"
	default:
		return nil
	}
}

func executeLine(line string) {
	line = trim(line)
	if line == "" {
		return
	}

	if strings.HasPrefix(line, "print(") && strings.HasSuffix(line, ")") {
		inner := line[6 : len(line)-1]
		inner = trim(inner)
		if isNumber(inner) {
			fmt.Println(inner)
		} else if strings.HasPrefix(inner, "\"") && strings.HasSuffix(inner, "\"") {
			fmt.Println(inner[1 : len(inner)-1])
		} else {
			if val, ok := vars[inner]; ok {
				fmt.Println(val)
			} else {
				fmt.Println("0")
			}
		}
		return
	}

	if isSyscall(line) {
		args := extractSyscallArgs(line)
		result := syscallDispatch(args)
		if result != nil {
			fmt.Println(result)
		} else {
			fmt.Println("0")
		}
		return
	}

	if strings.Contains(line, "=") {
		parts := strings.SplitN(line, "=", 2)
		lhs := trim(parts[0])
		rhs := trim(parts[1])
		var val interface{}
		if strings.HasPrefix(rhs, "\"") && strings.HasSuffix(rhs, "\"") {
			raw := rhs[1 : len(rhs)-1]
			raw = strings.ReplaceAll(raw, "\\\"", "\"")
			raw = strings.ReplaceAll(raw, "\\\\", "\\")
			val = raw
		} else if isNumber(rhs) {
			n, _ := strconv.Atoi(rhs)
			val = n
		} else if isSyscall(rhs) {
			args := extractSyscallArgs(rhs)
			val = syscallDispatch(args)
		} else {
			if v, ok := vars[rhs]; ok {
				val = v
			} else {
				val = 0
			}
		}
		vars[lhs] = val
		return
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Vire 引擎 - 用法: vire.exe 脚本.vire")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		executeLine(line)
	}
}