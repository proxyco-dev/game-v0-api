package common

import (
	"os"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Exclude(data interface{}, excludeFields []string) interface{} {
    excludeMap := make(map[string]map[string]bool)
    for _, field := range excludeFields {
        parts := strings.Split(field, ".")
        if len(parts) == 2 {
            if excludeMap[parts[0]] == nil {
                excludeMap[parts[0]] = make(map[string]bool)
            }
            excludeMap[parts[0]][parts[1]] = true
        } else {
            if excludeMap[""] == nil {
                excludeMap[""] = make(map[string]bool)
            }
            excludeMap[""][field] = true
        }
    }

    val := reflect.ValueOf(data)
    
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

	switch val.Kind() {
    case reflect.Slice:
        result := make([]map[string]interface{}, val.Len())
        for i := 0; i < val.Len(); i++ {
            item := val.Index(i)
            if item.Kind() == reflect.Struct {
                result[i] = excludeStruct(item, excludeMap)
            }
        }
        return result

    case reflect.Struct:
        return excludeStruct(val, excludeMap)

    default:
        return data
    }
}

func excludeStruct(val reflect.Value, excludeMap map[string]map[string]bool) map[string]interface{} {
    result := make(map[string]interface{})
    typ := val.Type()
    
    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i)
        jsonTag := field.Tag.Get("json")
        if jsonTag == "" {
            jsonTag = field.Name
        }
        jsonName := strings.Split(jsonTag, ",")[0]
        
        if excludeMap[""] != nil && excludeMap[""][jsonName] {
            continue
        }

        fieldValue := val.Field(i).Interface()

        if nestedExcludes := excludeMap[jsonName]; nestedExcludes != nil {
            if structVal := reflect.ValueOf(fieldValue); structVal.Kind() == reflect.Struct {
                nestedMap := map[string]map[string]bool{"": nestedExcludes}
                fieldValue = excludeStruct(structVal, nestedMap)
            }
        }

        result[jsonName] = fieldValue
    }
    return result
}

func AuthMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if len(auth) < 7 {
		return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
	}
	tokenString := auth[7:] // Remove "Bearer " prefix

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	c.Locals("user", token)
	return c.Next()
}
