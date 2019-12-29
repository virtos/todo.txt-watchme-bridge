package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

type configurationData struct {
	XMLName     xml.Name      `xml:"ConfigurationData"`
	Type        string        `xml:"xmlns,attr"`
	Timers      []timers      `xml:"Timers"`
	Settings    []settings    `xml:"Settings"`
	TimerAlerts []timerAlerts `xml:"TimerAlerts"`
	Categories  []categories  `xml:"Categories"`
}

type timers struct {
	XMLName      xml.Name `xml:"Timers"`
	TimerGUID    string   `xml:"TimerGUID"`
	CurrentTicks string   `xml:"CurrentTicks"`
	Name         string   `xml:"Name"`
	Note         string   `xml:"Note"`
	Autostar     string   `xml:"Autostart"`
	Autoreset    string   `xml:"Autoreset"`
	Hotkey       string   `xml:"Hotkey"`
	Category     string   `xml:"Category"`
	NoteHeight   string   `xml:"NoteHeight"`
	NoteWidth    string   `xml:"NoteWidth"`
	FlagIcon     string   `xml:"FlagIcon"`
	IsCountdown  string   `xml:"IsCountdown"`
	DefaultTicks string   `xml:"DefaultTicks"`
	LastStart    string   `xml:"LastStart"`
	LastStop     string   `xml:"LastStop"`
	LastReset    string   `xml:"LastReset"`
	ControlType  string   `xml:"ControlType"`
}

type settings struct {
	XMLName xml.Name `xml:"Settings"`
	Name    string   `xml:"Name"`
	Value   string   `xml:"Value"`
}

type timerAlerts struct {
	XMLName                          xml.Name `xml:"TimerAlerts"`
	TimerAlertGUID                   string   `xml:"TimerAlertGUID"`
	Hours                            string   `xml:"Hours"`
	Minutes                          string   `xml:"Minutes"`
	Seconds                          string   `xml:"Seconds"`
	DisplayMessageOnAlert            string   `xml:"DisplayMessageOnAlert"`
	StopAfterAlert                   string   `xml:"StopAfterAlert"`
	BeepOnAlert                      string   `xml:"BeepOnAlert"`
	LaunchAppOnAlert                 string   `xml:"LaunchAppOnAlert"`
	LaunchPath                       string   `xml:"LaunchPath"`
	PlaySoundOnAlert                 string   `xml:"PlaySoundOnAlert"`
	SoundPath                        string   `xml:"SoundPath"`
	TimerGUID                        string   `xml:"TimerGUID"`
	StartStopOtherTimerCountdown     string   `xml:"StartStopOtherTimerCountdown"`
	StartStopOtherTimerCountdownGUID string   `xml:"StartStopOtherTimerCountdownGUID"`
	Restart                          string   `xml:"Restart"`
	Reset                            string   `xml:"Reset"`
	IncrementCounter                 string   `xml:"IncrementCounter"`
	IncrementCounterGUID             string   `xml:"IncrementCounterGUID"`
	DecrementCounter                 string   `xml:"DecrementCounter"`
	DecrementCounterGUID             string   `xml:"DecrementCounterGUID"`
}

type categories struct {
	XMLName xml.Name `xml:"Categories"`
	Name    string   `xml:"Name"`
	Order   string   `xml:"Order"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getCurTodos(todoFN string) ([]string, error) {
	todoFile, err := os.Open(todoFN)
	check(err)
	defer todoFile.Close()

	scanner := bufio.NewScanner(todoFile)

	scanner.Split(bufio.ScanLines)

	var lines []string
	dt := " due:" + time.Now().Format("2006-01-02")

	for scanner.Scan() {
		curLine := scanner.Text()
		if strings.Contains(curLine, dt) {
			curLine = strings.Replace(curLine, dt, "", -1)
			lines = append(lines, curLine)
		}
	}
	fmt.Println("The following lines have been added: ", lines)
	return lines, nil
}

func getUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

var todoFN, wmConfigFN, templateFN string

func init() {
	const (
		defaultTodoFN     = "todo.txt"
		todoFNusage       = "A source todo.txt file name"
		defaultWMConfigFN = "WatchMeConfig.xml"
		wmConfigFNUsage   = "WatchMe configuration file name"
		defaultTemplateFN = "template.xml"
		templateFNUsage   = "Todo item template file name"
	)
	flag.StringVar(&todoFN, "td", defaultTodoFN, todoFNusage)
	flag.StringVar(&wmConfigFN, "w", defaultWMConfigFN, wmConfigFNUsage)
	flag.StringVar(&templateFN, "tt", defaultTemplateFN, templateFNUsage)
	flag.Parse()

	//Check if all nessessary files exist
	// fmt.Println("todoFN = ", todoFN)
	// fmt.Println("wmConfigFN = ", wmConfigFN)
	// fmt.Println("templateFN = ", templateFN)

	if _, err := os.Stat(todoFN); os.IsNotExist(err) {
		fmt.Println("Source todo.txt file not found: " + todoFN)
		os.Exit(2)
	}
	if _, err := os.Stat(wmConfigFN); os.IsNotExist(err) {
		fmt.Println("WatchMe config file not found: " + wmConfigFN)
		os.Exit(2)
	}
	if _, err := os.Stat(templateFN); os.IsNotExist(err) {
		fmt.Println("Template file for a todo item not found: " + templateFN)
		os.Exit(2)
	}
}

// Take Todo items from a todo.txt file and write them into a watchme config file
func main() {
	wmTemplateBackup := wmConfigFN + ".bak"
	rand.Seed(time.Now().UnixNano())

	todos, err := getCurTodos(todoFN)
	check(err)
	if len(todos) == 0 {
		os.Exit(1)
	}

	// Form an addition
	var template timers
	templateFile, err := os.Open(templateFN)
	check(err)
	defer templateFile.Close()

	templateByteValue, _ := ioutil.ReadAll(templateFile)
	xml.Unmarshal(templateByteValue, &template)

	template.Name = "Template1"

	xmlFile, err := os.Open(wmConfigFN)
	check(err)
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	xmlFile.Close()

	var configurationData configurationData
	xml.Unmarshal(byteValue, &configurationData)

	for _, curTodo := range todos {

		// Prevent duplicates
		duplicateExist := false
		for _, item := range configurationData.Timers {
			if item.Name == curTodo {
				duplicateExist = true
				break
			}
		}
		if !duplicateExist {
			addTimer := template
			addTimer.Name = curTodo
			addTimer.TimerGUID = getUUID()
			configurationData.Timers = append(configurationData.Timers, addTimer)
		}
	}

	// Sort All the tasks, including existing ones
	sort.Slice(configurationData.Timers, func(i, j int) bool {
		return configurationData.Timers[i].Name < configurationData.Timers[j].Name
	})

	// Write the output
	os.Remove(wmTemplateBackup)
	err = os.Rename(wmConfigFN, wmTemplateBackup)
	check(err)

	byteOut, err := xml.Marshal(configurationData)
	ioutil.WriteFile(wmConfigFN, byteOut, os.ModeExclusive)

}
