package main
import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os/exec"

	"strconv"

	"os"
	"strings"
)

type disk struct{
	diskName string
	model string
	serial string
	device string
	deviceSb string
	temperature int
	status string
}

type ups struct {
	model string
	status string
	linevolts float64
	loadPercentatge float64
	batteryCharge float64
	timeleft int64
	lotrans float64
	hitrans float64
	batteryVolts float64
	nominalVolts float64
	nominalBatteryVolts float64
	nominalPower int64



}
func main() {

	disksPath := flag.String("unraid-disks-path", "/var/local/emhttp/disks.ini", "Path where disks.ini is located, by default /var/local/emhttp/disks.ini")
    modules := flag.String("modules","disk,ups","modules to extract data separated with coma, default(disk,ups)  available disk,ups")
	flag.Parse()
	//*disksPath="D://my.ini"

	modulesSlice := strings.Split(*modules,",")
	for _, module := range modulesSlice {
		switch module {

		case "disk":
			diskMetrics(disksPath)
		case "ups":
			upsMetrics()
		}
	}
}

func upsMetrics() {
	cmd := exec.Command("apcaccess")
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err!=nil {
		log.Fatal("Error executting apcaccess: ",err)
	}
	/*cmd.
	output := string(out)
	if err != nil {
		log.Fatal(err)
	}*/
	//file, _ :=os.Open("D://test")
	scanner := bufio.NewScanner(stdout)
	var ups ups
	for scanner.Scan() {

		entry := strings.Split(scanner.Text(),":")
		switch strings.Trim(entry[0]," ") {
		case "MODEL":
			ups.model = strings.Trim(entry[1], " ")
		case "STATUS":
			ups.status = strings.Trim(entry[1], " ")
		case "LINEV":
			ups.linevolts, _ = strconv.ParseFloat(strings.Split(strings.Trim(entry[1], " "), " ")[0], 32)
		case "LOADPCT":
			ups.loadPercentatge, _ = strconv.ParseFloat(strings.Split(strings.Trim(entry[1], " "), " ")[0], 32)
		case "BCHARGE":
			ups.batteryCharge, _ = strconv.ParseFloat(strings.Split(strings.Trim(entry[1], " "), " ")[0], 32)
		case "TIMELEFT":
			values := strings.Split(strings.Trim(entry[1], " ")," ")
			number,_ := strconv.ParseFloat(values[0],32)
			switch values[1] {
			case "Minutes":
				number = number*60
			case "Hours":
				number = number * 3600
			case "Seconds":
				number = number
			}
			ups.timeleft = int64(number)
		case "NOMPOWER":
			ups.nominalPower, _ = strconv.ParseInt(strings.Split(strings.Trim(entry[1], " "), " ")[0], 10,32)
		}
	}
	cmd.Wait()
	fmt.Printf("unraid_ups,model=%s line_volts=%f,load_percentage=%f,battery_charge=%f,timeleft=%d,nominal_power=%d\n",influxNormalize(ups.model),ups.linevolts,ups.loadPercentatge,ups.batteryCharge,ups.timeleft,ups.nominalPower)

}

func diskMetrics(disksPath *string) {
	cfg, err := ini.Load(*disksPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	var disks []*disk
	for _, section := range cfg.Sections() {
		if strings.Trim(section.Key("device").String()," ") != "" {
			serialSeparator := strings.LastIndex(section.Key("id").String(),"_")
			model :=section.Key("id").String()[:serialSeparator]
			serial :=section.Key("id").String()[serialSeparator+1:]
			disks = append(disks, &disk{
				diskName:section.Key("name").String(),
				model:model,
				serial:serial,
				device:section.Key("device").String(),
				deviceSb: section.Key("deviceSb").String(),
				temperature:section.Key("temp").MustInt(-1),
				status: section.Key("status").String(),
			})
		}
	}

	//Print disks
	for _, disk := range disks {
		if disk.diskName!="DEFAULT"{
			fmt.Printf("unraid_disk_info,diskname=%s,device=%s,deviceSb=%s,model=%s,serial_no=%s value=0\n",influxNormalize(disk.diskName),influxNormalize(disk.device),influxNormalize(disk.deviceSb),influxNormalize(disk.model),influxNormalize(disk.serial))
		}
	}
}

func influxNormalize(s string) interface{} {
	return strings.Replace(s," ", "\\ ",-1)
}


