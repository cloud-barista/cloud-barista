package common

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"strconv"
	"github.com/labstack/echo"
	"regexp"
	//for ping
	//"github.com/sparrc/go-ping"
	
)

type benchInfo struct {
	Result string `json:"result"`
	Unit string `json:"unit"`
	Desc string `json:"desc"`
	Elapsed string `json:"elapsed"`
	SpecId string `json:"specid"`
}

type multiInfo struct {
	ResultArray []benchInfo `json:"resultarray"`
}

type request struct {
	Host string `json:"host"`
	Spec string `json:"spec"`
}

type mRequest struct {
	Multihost []request `json:"multihost"`
}


func RestGetInstall(c echo.Context) error {


	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	// wget install script from github install.sh
	cmdStr := "wget https://github.com/cloud-barista/cb-milkyway/raw/master/src/script/install.sh -P ~/script/"
	result, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in installation: wget script " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	// change chmod
	cmdStr = "sudo chmod 755 ~/script/install.sh"
	result2, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in installation: chmod " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	// run script
	cmdStr = "~/script/install.sh"
	result3, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in installation: chmod " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}
	
	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)

	result += result2
	result += result3

	result = "The installation is complete"

	content.Result = result
	content.Elapsed = elapsedStr 

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetInit(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	// Init fileio
	cmdStr := "sysbench fileio --file-total-size=50M prepare"
	outputStr, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: Init fileio " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`([0-9]+) files, ([0-9]+)([a-zA-Z]+) each, ([0-9]+)([a-zA-Z]+) total`)
	parseStr := grepStr.FindStringSubmatch(outputStr)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[0])
		fmt.Printf("File creation result: %s\n", parseStr1)

		outputStr = parseStr1
	}

	// Init DB
	cmdStr = "sysbench /usr/share/sysbench/oltp_read_write.lua --db-driver=mysql --table-size=100000 --mysql-db=sysbench --mysql-user=sysbench --mysql-password=psetri1234ak prepare"
	outputStr2, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: Init DB " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	grepStr = regexp.MustCompile(` ([0-9]+) records into .([a-zA-Z]+).`)
	parseStr = grepStr.FindStringSubmatch(outputStr2)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[0])
		fmt.Printf("Table creation result: %s\n", parseStr1)

		outputStr2 = parseStr1
	}
	
	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)

	outputStr += ", "
	outputStr += outputStr2


	//result = "The init is complete: "

	content.Result = "The init is complete"
	content.Elapsed = elapsedStr 

	content.Desc = outputStr + " are created"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetClean(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	// Clean fileio
	cmdStr := "sysbench fileio --file-total-size=50M cleanup"
	result, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: Clean fileio " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	// Clean DB
	cmdStr = "sysbench /usr/share/sysbench/oltp_read_write.lua --db-driver=mysql --table-size=100000 --mysql-db=sysbench --mysql-user=sysbench --mysql-password=psetri1234ak cleanup"
	result2, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: Clean DB "  + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}
	
	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)

	result += result2

	result = "The cleaning is complete"

	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "The benchmark files and tables are removed"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetCPUM(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	cores := strconv.Itoa(GetNumCPU())

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench cpu --cpu-max-prime=100000 --threads=" + cores + " run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: CPU"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`events per second:(\s+[+-]?([0-9]*[.])?[0-9]+)`)
	//for excution time:`execution time \(avg/stddev\):(\s+[+-]?([0-9]*[.])?[0-9]+)/`
	
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[1])
		fmt.Printf("execution time: %s\n", parseStr1)

		result = parseStr1
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Repeat the calculation (excution) for prime numbers in 100000 using " + cores + "cores"
	content.Unit = "Executions/sec"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetCPUS(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	cores := strconv.Itoa(1)

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench cpu --cpu-max-prime=100000 --threads=" + cores + " run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: CPU"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`events per second:(\s+[+-]?([0-9]*[.])?[0-9]+)`)
	//for excution time:`execution time \(avg/stddev\):(\s+[+-]?([0-9]*[.])?[0-9]+)/`
	
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[1])
		fmt.Printf("execution time: %s\n", parseStr1)

		result = parseStr1
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Repeat the calculation (excution) for prime numbers in 100000 using " + cores + "cores"
	content.Unit = "Executions/sec"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}


func RestGetMEMR(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench memory --memory-block-size=1K --memory-scope=global --memory-total-size=10G --memory-oper=read run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: MEMR"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(` transferred .([+-]?([0-9]*[.])?[0-9]+) `)
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[1])
		fmt.Printf("execution time: %s\n", parseStr1)

		result = parseStr1
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Allocate 10G memory buffer and read (repeat reading a pointer)"
	content.Unit = "MiB/sec"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetMEMW(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench memory --memory-block-size=1K --memory-scope=global --memory-total-size=10G --memory-oper=write run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: MEMW"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(` transferred .([+-]?([0-9]*[.])?[0-9]+) `)
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[1])
		fmt.Printf("execution time: %s\n", parseStr1)

		result = parseStr1
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Allocate 10G memory buffer and write (repeat writing a pointer)"
	content.Unit = "MiB/sec"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetFIOR(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench fileio --file-total-size=50M --file-test-mode=rndrd --max-time=30 --max-requests=0 run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: FIOR"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`read, MiB/s:(\s+[+-]?([0-9]*[.])?[0-9]+)`)
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[1])
		fmt.Printf("Throughput read, MiB/s: %s\n", parseStr1)

		result = parseStr1
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Check read throughput by excuting random reads for files in 50MiB for 30s"
	content.Unit = "MiB/sec"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetFIOW(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench fileio --file-total-size=50M --file-test-mode=rndwr --max-time=30 --max-requests=0 run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: FIOW"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`written, MiB/s:(\s+[+-]?([0-9]*[.])?[0-9]+)`)
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		parseStr1 := strings.TrimSpace(parseStr[1])
		fmt.Printf("Throughput write, MiB/s: %s\n", parseStr1)

		result = parseStr1
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Check write throughput by excuting random writes for files in 50MiB for 30s"
	content.Unit = "MiB/sec"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetDBR(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench /usr/share/sysbench/oltp_read_only.lua --db-driver=mysql --table-size=100000 --mysql-db=sysbench --mysql-user=sysbench --mysql-password=psetri1234ak run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: DBR"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`transactions:(\s+([0-9]*)(\s+)\([+-]?([0-9]*[.])?[0-9]+)`)
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		
		parseStr1 := strings.Split(parseStr[1], "(")
		fmt.Printf("DB Read Transactions/s: %s\n", parseStr1[1])

		result = parseStr1[1]
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Read transactions by simulating transaction loads (OLTP) in DB for 100000 records"
	content.Unit = "Transactions/s"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetDBW(c echo.Context) error {

	if(checkInit() != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: not initialized"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	content := benchInfo{}

	start := time.Now()

	fmt.Println("===============================================")

	cmdStr := "sysbench /usr/share/sysbench/oltp_write_only.lua --db-driver=mysql --table-size=100000 --mysql-db=sysbench --mysql-user=sysbench --mysql-password=psetri1234ak run"
	result, err := SysCall(cmdStr)

	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: DBW"}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`transactions:(\s+([0-9]*)(\s+)\([+-]?([0-9]*[.])?[0-9]+)`)
	parseStr := grepStr.FindStringSubmatch(result)	
	if len(parseStr) > 0 {
		
		parseStr1 := strings.Split(parseStr[1], "(")
		fmt.Printf("DB Write Transactions/s: %s\n", parseStr1[1])

		result = parseStr1[1]
	}
	
	content.Result = result
	content.Elapsed = elapsedStr 

	content.Desc = "Write transactions by simulating transaction loads (OLTP) in DB for 100000 records"
	content.Unit = "Transactions/s"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}


func RestGetRTT(c echo.Context) error {

	content := benchInfo{}

	req := request{}

	start := time.Now()

	fmt.Println("===============================================")

	if err := c.Bind(&req); err != nil {
		mapA := map[string]string{"message": "Error in request binding " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	pingHost := req.Host
	
	// system call for ping
	cmdStr := "ping -c 10 " + pingHost
	outputStr, err := SysCall(cmdStr)
	if(err != nil){
		mapA := map[string]string{"message": "Error in excuting the benchmark: Ping " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	var grepStr = regexp.MustCompile(`(\d+.\d+)/(\d+.\d+)/(\d+.\d+)/(\d+.\d+)`)
	parseStr := grepStr.FindAllStringSubmatch(outputStr, -1)
	if len(parseStr) > 0 {
		vals := parseStr[0]
		fmt.Printf("Ping result: %s\n", vals[1])

		outputStr = vals[2]
	}

	
	elapsed := time.Since(start)
	elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)

	
	content.Result = outputStr
	content.Elapsed = elapsedStr 
	content.Desc = "Average RTT to " + pingHost
	content.Unit = "ms"

	PrintJsonPretty(content)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &content)
}

func RestGetMultiRTT(c echo.Context) error {

	//content := benchInfo{}
	contentArray := multiInfo{}

	mReq := mRequest{}

	start := time.Now()

	fmt.Println("===============================================")

	if err := c.Bind(&mReq); err != nil {
		mapA := map[string]string{"message": "Error in request binding " + err.Error()}
		return c.JSON(http.StatusNotFound, &mapA)
	}

	hostList := mReq.Multihost
	for _, v := range hostList {
		content := benchInfo{}

		pingHost := v.Host
		// system call for ping
		cmdStr := "ping -c 10 " + pingHost
		outputStr, err := SysCall(cmdStr)
		if(err != nil){
			mapA := map[string]string{"message": "Error in excuting the benchmark: Ping " + err.Error()}
			return c.JSON(http.StatusNotFound, &mapA)
		}
	
		var grepStr = regexp.MustCompile(`(\d+.\d+)/(\d+.\d+)/(\d+.\d+)/(\d+.\d+)`)
		parseStr := grepStr.FindAllStringSubmatch(outputStr, -1)
		if len(parseStr) > 0 {
			vals := parseStr[0]
			fmt.Printf("Ping result: %s\n", vals[1])
	
			outputStr = vals[2]
		}

		elapsed := time.Since(start)
		elapsedStr := strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64)

		content.Result = outputStr
		content.Elapsed = elapsedStr 
		content.Desc = "Average RTT to " + pingHost
		content.Unit = "ms"
		content.SpecId = v.Spec

		contentArray.ResultArray = append(contentArray.ResultArray, content)
	
	}


	PrintJsonPretty(contentArray)
	fmt.Println("===============================================")

	return c.JSON(http.StatusOK, &contentArray)
}

func checkInit() error {
	checkPath, err := SysLookPath("sysbench")
	if(err != nil){
		return err
	}
	fmt.Printf("checkPath: %s\n", checkPath)
	return nil
}



func ApiValidation() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Printf("%v\n", "[API request!]")

			/*
			checkPath, err := SysLookPath("sysbench")
			mapA := map[string]string{"message": "Error in excuting the benchmark: no sysbench"}	
			if(err != nil){
				return echo.NewHTTPError(http.StatusNotFound, &mapA)
			}
			fmt.Printf("checkPath: %s\n", checkPath)
			*/

			return next(c)
		}
	}
}