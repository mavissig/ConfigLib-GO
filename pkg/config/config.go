package config

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	tagName = "envconfig"
)

type Options struct {
	printLogLoading bool
	prefix          string
	files           map[string]struct{}
}

type LoadOption func(*Options)

func LoadConfig[T any](options ...LoadOption) (*T, error) {
	opts := &Options{files: make(map[string]struct{})}
	for _, opt := range options {
		opt(opts)
	}

	for fileName := range opts.files {
		log.Printf("Loading .env file: %s", fileName)
		if err := godotenv.Load(fileName); err != nil {
			log.Printf("[WARN] Could not load file (%s): %v", fileName, err)
		}
	}

	fmt.Println("Environment variables:")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	var instance T
	fmt.Printf("Before envconfig.Process: %+v\n", instance)

	if err := envconfig.Process(opts.prefix, &instance); err != nil {
		log.Fatalf("Failed to process environment variables: %v\n", err)
	}
	fmt.Printf("After envconfig.Process: %+v\n", instance)

	if opts.printLogLoading {
		printConfig(instance, opts.prefix)
	}

	return &instance, nil
}
func AddFile(filePath string) func(*Options) {
	return func(opt *Options) {
		opt.files[filePath] = struct{}{}
	}
}

func AddFiles(files []string) func(*Options) {
	return func(opt *Options) {
		for _, filePath := range files {
			opt.files[filePath] = struct{}{}
		}
	}
}

func WithPrefix(prefix string) func(*Options) {
	return func(opt *Options) {
		opt.prefix = prefix
	}
}

func WithPrintConfig() func(*Options) {
	return func(opt *Options) {
		opt.printLogLoading = true
	}
}

func printConfig[T any](data T, prefix string) {
	fmt.Printf("--------------------------------\n")
	printConfigUtil(data, prefix)
	fmt.Printf("--------------------------------\n\n")
}

func printConfigUtil[T any](data T, prefix string) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct")
		return
	}

	fmt.Printf("Load %v\n", val.Type().String())

	typ := val.Type()

	maxLen := 0
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if len(field.Tag.Get(tagName)) > maxLen {
			maxLen = len(field.Tag.Get(tagName))
		}
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.Struct {
			fieldValue := val.Field(i).Interface()
			printConfigUtil(fieldValue, prefix)
			continue
		}

		value := val.Field(i)
		fmt.Printf("%-*s : %v\n", maxLen, field.Tag.Get(tagName), value)
	}

}
