package v2

// ImportsTemplate ...
const ImportsTemplate = `import (
    "fmt"
    "context"
    "encoding/json"
    "strings"
    "strconv"
    "sort"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-lambda-go/events"
)`
