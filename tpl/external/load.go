package external

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/clarkmcc/go-typescript"
	"github.com/dop251/goja"
	"github.com/gohugoio/hugo/helpers"
)

const maxFuncFileLoadTime = time.Second

type FunctionDetails struct {
	Name     string
	Func     func(args ...reflect.Value) (string, error)
	Examples [][2]string
}

type Namespace struct {
	funcs map[string]map[string]func(...reflect.Value) (string, error)
}

func (ns *Namespace) Function(name string, args ...reflect.Value) (string, error) {
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid function name: %s. Function names must be of the form fileName.FunctionName", name)
	}

	mFuncs, ok := ns.funcs[parts[0]]
	if !ok {
		return "", fmt.Errorf("the function file named %s has not been loaded", parts[0])
	}

	fn, ok := mFuncs[parts[1]]
	if !ok {
		return "", fmt.Errorf("the function named %s does not exist in %s", parts[1], parts[0])
	}

	return fn(args...)
}

func LoadFunctionFiles(funcsPath string) (*Namespace, error) {
	tsPaths, err := filepath.Glob(path.Join(funcsPath, "*.ts"))
	if err != nil {
		return nil, err
	}

	ns := &Namespace{
		funcs: make(map[string]map[string]func(...reflect.Value) (string, error)),
	}

	for _, tsPath := range tsPaths {
		mName := strings.TrimSuffix(path.Base(tsPath), ".ts")
		funcs, warns := loadTsFunctionsFile(tsPath)

		for _, warn := range warns {
			helpers.DistinctWarnLog.Warnf("Issue loading functions from %s: %v\n", tsPath, warn)
		}

		if len(funcs) == 0 {
			continue
		}

		ns.funcs[mName] = make(map[string]func(...reflect.Value) (string, error))

		for _, fn := range funcs {
			ns.funcs[mName][fn.Name] = fn.Func
			helpers.DistinctWarnLog.Infof("Loaded function: %s.%s\n", mName, fn.Name)
		}
	}

	return ns, nil
}

func loadTsFunctionsFile(tsPath string) ([]FunctionDetails, []error) {
	vm, exports, err := executeTS(tsPath)
	if err != nil {
		return nil, []error{err}
	}

	return extractFunctions(vm, exports)
}

func executeTS(tsPath string) (*goja.Runtime, goja.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxFuncFileLoadTime)
	defer cancel()

	tsFile, err := os.Open(tsPath)
	if err != nil {
		return nil, nil, err
	}

	var vm *goja.Runtime
	exportRuntime := func(cfg *typescript.EvaluateConfig) { vm = cfg.Runtime }

	_, err = typescript.EvaluateCtx(ctx, tsFile,
		typescript.WithTranspile(),
		typescript.WithAlmondModuleLoader(),
		exportRuntime,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not evaluate function file: %v", err)
	}

	exports, err := vm.RunString("exports")
	if err != nil {
		return nil, nil, fmt.Errorf("the function file is missing exports: %v", err)
	}

	return vm, exports, nil
}

func extractFunctions(vm *goja.Runtime, exports goja.Value) ([]FunctionDetails, []error) {
	var funcs []FunctionDetails
	var warns []error

	exportMap, ok := exports.Export().(map[string]any)
	if !ok {
		return nil, []error{fmt.Errorf("the function file's 'exports' variable is not a map")}
	}

	for name, obj := range exportMap {
		fn, ok := obj.(func(goja.FunctionCall) goja.Value)
		if !ok {
			// Only looking for exported functions
			continue
		}

		detail := FunctionDetails{
			Name: name,
			Func: func(args ...reflect.Value) (string, error) {
				valueArgs := make([]goja.Value, len(args))
				for i, arg := range args {
					val, err := toGojaType(vm, arg)
					if err != nil {
						return "", err
					}
					valueArgs[i] = val
				}

				val := fn(goja.FunctionCall{Arguments: valueArgs})

				var out string
				err := vm.ExportTo(val, &out)
				return out, err
			},
		}

		if obj, ok := exportMap[name+"Examples"]; ok {
			examples, err := extractExamples(obj.([]interface{}))
			if err == nil {
				detail.Examples = examples
			} else {
				warns = append(warns, err)
			}
		}

		funcs = append(funcs, detail)
	}

	return funcs, warns
}

func extractExamples(ifaces []interface{}) ([][2]string, error) {
	var examples [][2]string

	_ = ifaces

	return examples, nil
}

func toGojaType(vm *goja.Runtime, val reflect.Value) (goja.Value, error) {
	iface := val.Interface()

	if t, ok := iface.(time.Time); ok {
		return vm.New(vm.Get("Date").ToObject(vm), vm.ToValue(t.UnixNano()/1e6))
	}

	return vm.ToValue(val.Interface()), nil
}
