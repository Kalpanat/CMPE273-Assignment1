package main
import (
"github.com/gorilla/rpc/json"
"net/http"
"fmt"
"log"
"bytes"
"regexp"
"strings"
"strconv"
"os"
)

type RPCAPIArguments struct {
 StockMessage string
 Budget float32
}
type RPCAPITRADEArguement struct {
 TradeId int
}
type RPCAPIResponse struct {
 Message string
}
type RPCAPITradeResponse struct {
 Message string
}
func jsonRpcStockCall(method string, args RPCAPIArguments) (reply RPCAPIResponse, err error) {
    buf, _ := json.EncodeClientRequest(method, args)

    resp, err := http.Post("http://localhost:8080/rpc", "application/json", bytes.NewBuffer(buf))
    if err != nil {
        return
    }
    defer resp.Body.Close()
    err = json.DecodeClientResponse(resp.Body, &reply)
    return
}
func jsonRpcTradeCall(method string, args RPCAPITRADEArguement) (reply RPCAPITradeResponse, err error) {
    buf, _ := json.EncodeClientRequest(method, args)
    resp, err := http.Post("http://localhost:8080/rpc", "application/json", bytes.NewBuffer(buf))
    if err != nil {
        return
    }
    defer resp.Body.Close()
    err = json.DecodeClientResponse(resp.Body, &reply)
    return
}
func main() {
	//fmt.Println("Starting client")
    fmt.Print("For Buying stocks Enter 1/ For Checking your portfolio Enter 2: ")
    var inputOption int64
    fmt.Scanf("%d", &inputOption)

    if inputOption==1{
        fmt.Print("Enter stock request in following format--> stockSymbol:Percentage,stockSymbol:Percentage... : ")
        var inputStock string
        fmt.Scanf("%s", &inputStock)

        fmt.Print("Enter budget: ") 
        var inputBudget float32
        fmt.Scanf("%f",&inputBudget)

        inputFormat := strings.Replace(inputStock, "%", "", -1)
        re := regexp.MustCompile("[:,;]")
        result := re.Split(inputFormat, -1)
        var totalPercentage int64
        for i := range(result) {
            if i%2!=0{
                number, _ := strconv.ParseInt(result[i], 10, 0)
                //fmt.Println("number is ",number)
                totalPercentage+=number
                }
            }

        if (totalPercentage!=100){
            fmt.Println("Total stock percentage should be 100, found",totalPercentage)
            os.Exit(3)
        }else{
            reply, err := jsonRpcStockCall("FinanceService.Response", RPCAPIArguments{inputStock,inputBudget})
            if err != nil {
                log.Fatal("Error: ",err)
            }
            fmt.Println(reply.Message)
        }
    }else if inputOption==2{
        
        fmt.Print("Please enter TradeId: ") 
        var inputTradeId int
        fmt.Scanf("%d",&inputTradeId)


        if(inputTradeId==0){
            fmt.Println("Not valid tradeId :",inputTradeId)
            os.Exit(3)
        }else{
            reply, err := jsonRpcTradeCall("FinanceService.TradeRequest", RPCAPITRADEArguement{TradeId:inputTradeId})
            if err != nil {
                fmt.Println("Error occured :")
                log.Fatal("Error: ",err)
            }
            fmt.Println(reply.Message)
        }
    }else{
        fmt.Println("Please select correct option.")
        os.Exit(3)
    }
}