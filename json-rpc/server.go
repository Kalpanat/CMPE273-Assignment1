package main
import (
 "github.com/gorilla/rpc"
 myencoding "encoding/json"
 "github.com/gorilla/rpc/json"
 "net/http"
 "fmt"
 "log"
 "regexp"
 "strings"
 "strconv"
 "math/rand"
 "time"
 "io/ioutil"
 "os"
 "math"
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
func randomNum() int{
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    return r1.Intn(10000)
}
var finalMap= make(map[int]string)

type Target struct {
    Query struct {
        Count int `json:"count"`
        Created time.Time `json:"created"`
        Lang string `json:"lang"`
        Results struct {
            Quote []struct {
                Symbol string `json:"symbol"`
                LastTradePriceOnly string `json:"LastTradePriceOnly"`
            } `json:"quote"`
        } `json:"results"`
    } `json:"query"`
}

type FinanceService struct{}

func (h *FinanceService) TradeRequest(r *http.Request, args *RPCAPITRADEArguement, reply *RPCAPITradeResponse) error {
        var msgString string
        var currentMarketValue float64
        var CurrUnvestedAmt float64
        for key,value := range finalMap {
            //fmt.Println("key***",key,"args.TradeId",args.TradeId)
            if key==args.TradeId{
                formatString := strings.Replace(value, "%", "", -1)
                formatString=strings.Replace(formatString, "$", "", -1)
                formatString=strings.Replace(formatString, " ", "", -1)

                s:=strings.Split(formatString, "#")
                var stockSymbolsFormat string

                res := regexp.MustCompile("[:,]")
                result := res.Split(s[0], -1)
                //fmt.Println("Result length:",len(result))
                lengthArraay := len(result)

                stk := make([]int, 1)
                stk[0]=0
                switch {
                    case lengthArraay == 6:stk = append(stk, 3)
                    case lengthArraay == 9:stk = append(stk, 3);stk = append(stk, 6)
                    case lengthArraay == 12:stk = append(stk, 3);stk = append(stk, 6);stk = append(stk, 9)
                    case lengthArraay == 15:stk = append(stk, 3);stk = append(stk, 6);stk = append(stk, 9);stk = append(stk, 12)
                    case lengthArraay == 18:stk = append(stk, 3);stk = append(stk, 6);stk = append(stk, 9);stk = append(stk, 12);stk = append(stk, 15)
                    case lengthArraay == 21:stk = append(stk, 3);stk = append(stk, 6);stk = append(stk, 9);stk = append(stk, 12);stk = append(stk, 15);stk = append(stk, 18)
                    case lengthArraay == 24:stk = append(stk, 3);stk = append(stk, 6);stk = append(stk, 9);stk = append(stk, 12);stk = append(stk, 15);stk = append(stk, 18);stk = append(stk, 21)
                    case lengthArraay == 27:stk = append(stk, 3);stk = append(stk, 6);stk = append(stk, 9);stk = append(stk, 12);stk = append(stk, 15);stk = append(stk, 18);stk = append(stk, 21);stk = append(stk, 24)
                }
                /*
                for i := range(result){
                    fmt.Println("val****",result[i]," i*** ",i)
                }*/
                for j := range(stk) {
                        stockSymbolsFormat+="%22"+result[stk[j]]+"%22,"
                }
                stockSymbolsFormat=stockSymbolsFormat[0:len(stockSymbolsFormat)-1]
                var tr Target
                var priceMap= make(map[string]string)
                

                url :="http://query.yahooapis.com/v1/public/yql?q=select%20symbol,%20LastTradePriceOnly%20from%20yahoo.finance.quote%20where%20symbol%20in%20("+stockSymbolsFormat+")&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"

                responseValue, err := http.Get(url)
                if err != nil {
                    fmt.Printf("%s", err)
                    os.Exit(1)
                } else {
                    defer responseValue.Body.Close()
                    reply, err := ioutil.ReadAll(responseValue.Body)
                    myencoding.Unmarshal([]byte(reply), &tr)
                    if err != nil {
                        fmt.Printf("%s", err)
                        os.Exit(1)
                    }
                    for i := 0; i < (lengthArraay/3); i++{
                            var signString string
                            signString=tr.Query.Results.Quote[i].Symbol+":"+result[stk[i]+1]

                            currPrice, err := strconv.ParseFloat(tr.Query.Results.Quote[i].LastTradePriceOnly, 64)
                            if err != nil {
                                    fmt.Printf("%s", err)
                                }
                            numStk, err := strconv.ParseFloat(result[stk[i]+1], 64)
                            if err != nil {
                                    fmt.Printf("%s", err)
                                }
                            currentMarketValue+=(currPrice*numStk)
                            prevPrice, err := strconv.ParseFloat(result[stk[i]+2], 64)
                            if err != nil {
                                    fmt.Printf("%s", err)
                                }
                            if(currPrice>prevPrice){
                                signString+=":+$"+tr.Query.Results.Quote[i].LastTradePriceOnly
                            } else if(currPrice<prevPrice){
                                signString+=":-$"+tr.Query.Results.Quote[i].LastTradePriceOnly
                            } else {
                                signString+=":$"+tr.Query.Results.Quote[i].LastTradePriceOnly
                            }
                            priceMap[tr.Query.Results.Quote[i].Symbol]=signString   
                    }
                }
                msgString=""
                inputAmt, err := strconv.ParseFloat(s[1], 64)
                    if err != nil {
                            fmt.Printf("%s", err)
                        }
                CurrUnvestedAmt=inputAmt-currentMarketValue
                for _,value := range priceMap{
                    msgString+=value+" ,"
                }
                msgString=msgString[0:len(msgString)-1]
                break;
            } else {
                msgString=""
                msgString="TradeId doesnt exists in system."
            }
        }
        reply.Message = fmt.Sprintf("StockString: %s CurrentMarketValue: %f UnvestedAmount: %f",msgString,float32(currentMarketValue),float32(CurrUnvestedAmt))
        return nil
}

func (h *FinanceService) Response(r *http.Request, args *RPCAPIArguments, reply *RPCAPIResponse) error {
    inputStock :=args.StockMessage
    inputBudget := args.Budget

    var stockmap map[string]float64
    stockmap = make(map[string]float64)

    inputFormat := strings.Replace(inputStock, "%", "", -1)
    re := regexp.MustCompile("[:,;]")
    result := re.Split(inputFormat, -1)

    for i := range(result) {
        if i%2!=0{
            number, _ := strconv.ParseInt(result[i], 10, 0)
            percentVal := ((float64(inputBudget)*float64(number))/100)
            stockmap[result[i-1]]= percentVal
        }
    }
    var stockSymbols string
    for key := range stockmap {
        stockSymbols+="%22"+key+"%22,"
    }
    formatStockSymbols :=stockSymbols[0:len(stockSymbols)-1]
	var t Target
    url :="http://query.yahooapis.com/v1/public/yql?q=select%20symbol,%20LastTradePriceOnly%20from%20yahoo.finance.quote%20where%20symbol%20in%20("+formatStockSymbols+")&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"
    tradeId:=randomNum()

    var unvestedAmount float64
    var stockString string
    responseVal, err := http.Get(url)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        } else {
            defer responseVal.Body.Close()
            reply, err := ioutil.ReadAll(responseVal.Body)
            myencoding.Unmarshal([]byte(reply), &t)
            if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
            }            
            for key, value := range stockmap {
                for i := 0; i < len(stockmap); i++ {
                    if(t.Query.Results.Quote[i].Symbol==key){
                        f, err := strconv.ParseFloat(t.Query.Results.Quote[i].LastTradePriceOnly, 64)
                        if err != nil {
                            fmt.Printf("%s", err)
                        }
                        numberOfStock := int(value/f)
                        unvestedAmount+=math.Mod(value, f)
                        stockString+=t.Query.Results.Quote[i].Symbol+":"+strconv.Itoa(numberOfStock)+":$"+t.Query.Results.Quote[i].LastTradePriceOnly+","
                    }
                }
            }
        }
        stockString=stockString[0:len(stockString)-1]
        inputBudgetStr:= strconv.FormatFloat(float64(inputBudget), 'E', 2, 64)
        finalMap[tradeId]=stockString+"#"+inputBudgetStr
        reply.Message = fmt.Sprintf("TradeId : %d StockString: %s UnvestedAmount : %f",tradeId,stockString,unvestedAmount) + "\n"
    	return nil
}
func handler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi there, hello %s!", r.URL.Path[1:])
    }
func main() {
        fmt.Println("Starting HTTP server")
        s := rpc.NewServer()
        s.RegisterCodec(json.NewCodec(), "application/json")
        s.RegisterService(new(FinanceService), "")
        http.Handle("/rpc", s)
        http.HandleFunc("/", handler)
        http.ListenAndServe(":8080", nil)
        log.Println("Starting JSON-RPC server on localhost:8080/rpc")
        log.Fatal(http.ListenAndServe(":8080", nil))
}