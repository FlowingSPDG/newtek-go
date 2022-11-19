package newtek

import "encoding/xml"

// ProductInformation. /v1/version API
type ProductInformation struct {
	XMLName            xml.Name `xml:"product_information"`
	Text               string   `xml:",chardata"`
	ProductModel       string   `xml:"product_model"`
	ProductName        string   `xml:"product_name"`
	ProductVersion     string   `xml:"product_version"`
	ProductID          string   `xml:"product_id"`
	ProductSerialNo    string   `xml:"product_serial_no"`
	ProductBuildNo     string   `xml:"product_build_no"`
	MachineName        string   `xml:"machine_name"`
	SessionXResolution string   `xml:"session_x_resolution"`
	SessionYResolution string   `xml:"session_y_resolution"`
	SessionFielded     string   `xml:"session_fielded"`
	SessionFrameRate   string   `xml:"session_frame_rate"`
	SessionAspectRatio string   `xml:"session_aspect_ratio"`
	SessionColorFormat string   `xml:"session_color_format"`
	SessionColorCoding string   `xml:"session_color_coding"`
	SessionName        string   `xml:"session_name"`
	OutputCount        string   `xml:"output_count"`
}

type Tally struct {
	XMLName xml.Name `xml:"tally"`
	Text    string   `xml:",chardata"`
	Column  []struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"name,attr"`
		Index  string `xml:"index,attr"`
		OnPgm  string `xml:"on_pgm,attr"`
		OnPrev string `xml:"on_prev,attr"`
		NdiID  string `xml:"ndi_id,attr"`
	} `xml:"column"`
}

type ShortcutStates struct {
	XMLName       xml.Name `xml:"shortcut_states"`
	Text          string   `xml:",chardata"`
	ShortcutState []struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"name,attr"`
		Value  string `xml:"value,attr"`
		Type   string `xml:"type,attr"`
		Sender string `xml:"sender,attr"`
	} `xml:"shortcut_state"`
}

type Metadata struct {
	XMLName xml.Name `xml:"metadata"`
	Text    string   `xml:",chardata"`
	Clips   []Clip   `xml:"clips"`
}

type Clip struct {
	Text string `xml:",chardata"`
	Mark []struct {
		Text   string   `xml:",chardata"`
		ID     string   `xml:"id,attr"`
		In     string   `xml:"in,attr"`
		Out    string   `xml:"out,attr"`
		Camera []Camera `xml:"camera"`
	} `xml:"mark"`
}

type Camera struct {
	Text  string `xml:",chardata"`
	Angle string `xml:"angle,attr"`
	Tag   string `xml:"tag,attr"`
	Path  string `xml:"path,attr"`
}
