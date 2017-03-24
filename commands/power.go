package commands

import (
	"fmt"

	"github.com/vapor-ware/vesh/client"
	"github.com/vapor-ware/vesh/utils"
)

const powerpath = "power/"
const device_id = "power"

type PowerDetails struct {
	InputPower  float64 `json:"input_power"`
	OverCurrent bool    `json:"over_current"`
	PowerOK     bool    `json:"power_ok"`
	PowerStatus string  `json:"power_status"`
}

type PowerResult struct {
	utils.Result
	*PowerDetails
}

// ListPower iterates over the complete list of devices and returns input power,
// over current, power ok, and power status for each `power` device type.
func ListPower(vc *client.VeshClient, filter func(res utils.Result) bool) ([]PowerResult, error) {
	var devices []utils.Result

	var data []PowerResult

	for res := range utils.FilterDevices(filter) {
		devices = append(devices, res)
	}

	progressBar := utils.ProgressBar(len(devices))

	for _, res := range devices {
		power, _ := GetPower(vc, res)
		data = append(data, PowerResult{res, power})
		progressBar.Incr()
	}

	uiprogress.Stop()
	return data, nil
}

func GetPower(vc *client.VeshClient, res utils.Result) (*PowerDetails, error) {
	power := &PowerDetails{}
	path := fmt.Sprintf("%s/%s/%s", res.RackID, res.BoardID, res.DeviceID)
	_, err := vc.Sling.New().Path(powerpath).Get(path).ReceiveSuccess(power)
	if err != nil {
		return power, err
	}

	return power, nil
>>>>>>> Refactored power
}

// PrintListPower takes the output from ListPower and pretty prints it into a table.
// Multiple lights are grouped by board, then by rack. Table format is set to not
// auto merge duplicate entries.
func PrintListPower(vc *client.VeshClient) error {
	filter := func(res utils.Result) bool {
		return res.DeviceType == device_id
	}

	header := []string{"Rack", "Board", "Name", "Input Power", "Power Ok?"}
	powerList, _ := ListPower(vc, filter)

	var data [][]string

	for _, res := range powerList {
		data = append(data, []string{
			res.RackID,
			res.BoardID,
			res.DeviceInfo,
			fmt.Sprintf("%.2f", res.InputPower),
			fmt.Sprintf("%t", res.PowerOK)})
	}

	utils.TableOutput(header, data)

	return nil
}

// PrintGetPower takes the output of GetPower and pretty prints it in table form.
// Multiple entries are not merged.
func PrintGetPower(vc *client.VeshClient, rack_id, board_id string) error {
	filter := func(res utils.Result) bool {
		return res.DeviceType == device_id && res.RackID == rack_id && res.BoardID == board_id
	}

	header := []string{"Rack", "Board", "Device", "Name", "Input Power", "Over Current?", "Power Ok?", "Power Status"}
	powerList, _ := ListPower(vc, filter) // Add error reporting

	var data [][]string

	for _, res := range powerList {
		data = append(data, []string{
			res.RackID,
			res.BoardID,
			res.DeviceID,
			res.DeviceInfo,
			fmt.Sprintf("%.2f", res.InputPower),
			fmt.Sprintf("%t", res.OverCurrent),
			fmt.Sprintf("%t", res.PowerOK),
			res.PowerStatus})
	}

	utils.TableOutput(header, data)

	return nil
}

// SetPower takes a rack and board id as a locator as well as a power status
// string. The power status of the corresponding "power" device is set to the
// given power status.
// Options are: "on", "off", "cycle"
func SetPower(vc *client.VeshClient, rack_id, board_id, power_status string) (string, error) {
	responseData := &PowerDetails{}
	resp, err := vc.Sling.New().Path(powerpath).Path(rack_id + "/").Path(board_id + "/").Path(device_id + "/").Get(power_status).ReceiveSuccess(responseData) // Add error reporting
	if resp.StatusCode != 200 {                                                                                                                               // This is not what I meant by "error reporting"
		return "", err
	}
	if err == nil && power_status == "cycle" { // This should check if successful
		return power_status, err
	}
	return responseData.PowerStatus, err
}

// PrintSetPower takes the output of SetPower and pretty prints whether the
// status was changed successfully.
func PrintSetPower(vc *client.VeshClient, rack_id, board_id, power_status string) error {
	status, err := SetPower(vc, rack_id, board_id, power_status)
	if err == nil && status == "cycle" {
		fmt.Printf("Power successfully %sd\n", status)
	}
	if err == nil && status != "cycle" {
		fmt.Printf("Power set to %s\n", status)
	}
	return err
}
