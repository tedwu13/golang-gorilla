
import "encoding/json"
import "fmt"
import "os"

type JSONResponse1 struct {
    Page    int
    Fruits  []string
}

type JSONResponse2 struct {
    Page    int `json:"page"`
    Fruits  []string `json:"fruits"`
}

