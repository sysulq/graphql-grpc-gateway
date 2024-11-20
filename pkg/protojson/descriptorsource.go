package protojson

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/fullstorydev/grpcurl"
	// nolint
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
)

type Method struct {
	HttpMethod string
	PathNames  []string
	HttpPath   string
	RpcPath    string
}

// GetMethods returns all methods of the given grpcurl.DescriptorSource.
func GetMethods(source grpcurl.DescriptorSource) ([]Method, error) {
	svcs, err := source.ListServices()
	if err != nil {
		return nil, err
	}

	var methods []Method
	for _, svc := range svcs {
		d, err := source.FindSymbol(svc)
		if err != nil {
			return nil, err
		}

		switch val := d.(type) {
		case *desc.ServiceDescriptor:
			svcMethods := val.GetMethods()
			for _, method := range svcMethods {
				rpcPath := fmt.Sprintf("%s/%s", svc, method.GetName())
				ext := proto.GetExtension(method.GetMethodOptions(), annotations.E_Http)
				switch rule := ext.(type) {
				case *annotations.HttpRule:
					if rule == nil {
						methods = append(methods, Method{
							RpcPath: rpcPath,
						})
						continue
					}

					switch httpRule := rule.GetPattern().(type) {
					case *annotations.HttpRule_Get:
						methods = append(methods, Method{
							HttpMethod: http.MethodGet,
							PathNames:  extractFieldNames(httpRule.Get),
							HttpPath:   adjustHttpPath(httpRule.Get),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Post:
						methods = append(methods, Method{
							HttpMethod: http.MethodPost,
							PathNames:  extractFieldNames(httpRule.Post),
							HttpPath:   adjustHttpPath(httpRule.Post),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Put:
						methods = append(methods, Method{
							HttpMethod: http.MethodPut,
							PathNames:  extractFieldNames(httpRule.Put),
							HttpPath:   adjustHttpPath(httpRule.Put),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Delete:
						methods = append(methods, Method{
							HttpMethod: http.MethodDelete,
							PathNames:  extractFieldNames(httpRule.Delete),
							HttpPath:   adjustHttpPath(httpRule.Delete),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Patch:
						methods = append(methods, Method{
							HttpMethod: http.MethodPatch,
							PathNames:  extractFieldNames(httpRule.Patch),
							HttpPath:   adjustHttpPath(httpRule.Patch),
							RpcPath:    rpcPath,
						})
					default:
						methods = append(methods, Method{
							RpcPath: rpcPath,
						})
					}
				default:
					methods = append(methods, Method{
						RpcPath: rpcPath,
					})
				}
			}
		}
	}

	return methods, nil
}

func adjustHttpPath(path string) string {
	// path = strings.ReplaceAll(path, "{", ":")
	// path = strings.ReplaceAll(path, "}", "")
	return path
}

// extractFieldNames extracts all field names (e.g., "name", "id") from the given path template.
func extractFieldNames(template string) []string {
	// Regular expression to match both {name=...} and {name}
	re := regexp.MustCompile(`\{([^=}]+)(=[^}]*)?\}`)

	// Find all matches
	matches := re.FindAllStringSubmatch(template, -1)
	if len(matches) == 0 {
		return nil
	}

	// Extract the field names
	var fieldNames []string
	for _, match := range matches {
		if len(match) > 1 {
			fieldNames = append(fieldNames, match[1])
		}
	}

	return fieldNames
}
