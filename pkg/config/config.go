package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	tagName = "envconfig"
)

type Options struct {
	printLogLoading bool
	files           map[string]struct{}
}

type LoadOption func(*Options)

func LoadConfig[T any](options ...LoadOption) (*T, error) {
	opts := &Options{files: make(map[string]struct{})}
	for _, opt := range options {
		opt(opts)
	}

	for fileName := range opts.files {
		err := godotenv.Load(fileName)
		if err != nil {
			log.Println("[API][TRANSPORT][CONFIG] WARN ", err)
		}
	}

	var instance T
	val := reflect.ValueOf(&instance).Elem()

	if val.Kind() != reflect.Struct {

	}

	if err := envconfig.Process("", &instance); err != nil {
		log.Fatalln(err)
	}

	if opts.printLogLoading {
		printConfig(instance)
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

func WithPrintConfig() func(*Options) {
	return func(opt *Options) {
		opt.printLogLoading = true
	}
}

func printConfig[T any](data T) {
	fmt.Printf("--------------------------------\n")
	printConfigUtil(data)
	fmt.Printf("--------------------------------\n\n")
}

func printConfigUtil[T any](data T) {
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
			printConfigUtil(fieldValue)
			continue
		}

		value := val.Field(i)
		fmt.Printf("%-*s : %v\n", maxLen, field.Tag.Get(tagName), value)
	}

}
