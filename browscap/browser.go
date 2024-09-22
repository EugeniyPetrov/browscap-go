package browscap

type Browser struct {
	Pattern                    string
	Comment                    string
	Browser                    string
	BrowserType                string
	BrowserBits                int
	BrowserMaker               string
	BrowserModus               string
	Version                    string
	MajorVer                   string
	MinorVer                   string
	Platform                   string
	PlatformVersion            string
	PlatformDescription        string
	PlatformBits               int
	PlatformMaker              string
	Alpha                      bool
	Beta                       bool
	Win16                      bool
	Win32                      bool
	Win64                      bool
	Frames                     bool
	Iframes                    bool
	Tables                     bool
	Cookies                    bool
	BackgroundSounds           bool
	Javascript                 bool
	VBScript                   bool
	JavaApplets                bool
	ActiveXControls            bool
	IsMobileDevice             bool
	IsTablet                   bool
	IsSyndicationReader        bool
	Crawler                    bool
	IsFake                     bool
	IsAnonymized               bool
	IsModified                 bool
	CSSVersion                 int
	AolVersion                 int
	DeviceName                 string
	DeviceMaker                string
	DeviceType                 string
	DevicePointingMethod       string
	DeviceCodeName             string
	DeviceBrandName            string
	RenderingEngineName        string
	RenderingEngineVersion     string
	RenderingEngineDescription string
	RenderingEngineMaker       string
}

type BrowserNode struct {
	ID                         int     `db:"id"`
	Parent                     string  `mapstructure:"Parent" db:"parent"`
	Pattern                    string  `mapstructure:"Pattern" db:"pattern"`
	Comment                    *string `mapstructure:"Comment" db:"comment" browscap:""`
	Browser                    *string `mapstructure:"Browser" db:"browser"`
	BrowserType                *string `mapstructure:"Browser_Type" db:"browser_type"`
	BrowserBits                *int    `mapstructure:"Browser_Bits" db:"browser_bits"`
	BrowserMaker               *string `mapstructure:"Browser_Maker" db:"browser_maker"`
	BrowserModus               *string `mapstructure:"Browser_Modus" db:"browser_modus"`
	Version                    *string `mapstructure:"Version" db:"version"`
	MajorVer                   *string `mapstructure:"MajorVer" db:"major_ver"`
	MinorVer                   *string `mapstructure:"MinorVer" db:"minor_ver"`
	Platform                   *string `mapstructure:"Platform" db:"platform"`
	PlatformVersion            *string `mapstructure:"Platform_Version" db:"platform_version"`
	PlatformDescription        *string `mapstructure:"Platform_Description" db:"platform_description"`
	PlatformBits               *int    `mapstructure:"Platform_Bits" db:"platform_bits"`
	PlatformMaker              *string `mapstructure:"Platform_Maker" db:"platform_maker"`
	Alpha                      *bool   `mapstructure:"Alpha" db:"alpha"`
	Beta                       *bool   `mapstructure:"Beta" db:"beta"`
	Win16                      *bool   `mapstructure:"Win16" db:"win16"`
	Win32                      *bool   `mapstructure:"Win32" db:"win32"`
	Win64                      *bool   `mapstructure:"Win64" db:"win64"`
	Frames                     *bool   `mapstructure:"Frames" db:"frames"`
	Iframes                    *bool   `mapstructure:"IFrames" db:"iframes"`
	Tables                     *bool   `mapstructure:"Tables" db:"tables"`
	Cookies                    *bool   `mapstructure:"Cookies" db:"cookies"`
	BackgroundSounds           *bool   `mapstructure:"BackgroundSounds" db:"background_sounds"`
	Javascript                 *bool   `mapstructure:"JavaScript" db:"javascript"`
	VBScript                   *bool   `mapstructure:"VBScript" db:"vbscript"`
	JavaApplets                *bool   `mapstructure:"JavaApplets" db:"java_applets"`
	ActiveXControls            *bool   `mapstructure:"ActiveXControls" db:"activex_controls"`
	IsMobileDevice             *bool   `mapstructure:"isMobileDevice" db:"is_mobile_device"`
	IsTablet                   *bool   `mapstructure:"isTablet" db:"is_tablet"`
	IsSyndicationReader        *bool   `mapstructure:"isSyndicationReader" db:"is_syndication_reader"`
	Crawler                    *bool   `mapstructure:"Crawler" db:"crawler"`
	IsFake                     *bool   `mapstructure:"isFake" db:"is_fake"`
	IsAnonymized               *bool   `mapstructure:"isAnonymized" db:"is_anonymized"`
	IsModified                 *bool   `mapstructure:"isModified" db:"is_modified"`
	CSSVersion                 *int    `mapstructure:"CssVersion" db:"css_version"`
	AolVersion                 *int    `mapstructure:"AolVersion" db:"aol_version"`
	DeviceName                 *string `mapstructure:"Device_Name" db:"device_name"`
	DeviceMaker                *string `mapstructure:"Device_Maker" db:"device_maker"`
	DeviceType                 *string `mapstructure:"Device_Type" db:"device_type"`
	DevicePointingMethod       *string `mapstructure:"Device_Pointing_Method" db:"device_pointing_method"`
	DeviceCodeName             *string `mapstructure:"Device_Code_Name" db:"device_code_name"`
	DeviceBrandName            *string `mapstructure:"Device_Brand_Name" db:"device_brand_name"`
	RenderingEngineName        *string `mapstructure:"RenderingEngine_Name" db:"rendering_engine_name"`
	RenderingEngineVersion     *string `mapstructure:"RenderingEngine_Version" db:"rendering_engine_version"`
	RenderingEngineDescription *string `mapstructure:"RenderingEngine_Description" db:"rendering_engine_description"`
	RenderingEngineMaker       *string `mapstructure:"RenderingEngine_Maker" db:"rendering_engine_maker"`
}

func String(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func Int(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func Bool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func (n *BrowserNode) ToBrowser() *Browser {
	return &Browser{
		Pattern:                    n.Pattern,
		Comment:                    String(n.Comment),
		Browser:                    String(n.Browser),
		BrowserType:                String(n.BrowserType),
		BrowserBits:                Int(n.BrowserBits),
		BrowserMaker:               String(n.BrowserMaker),
		BrowserModus:               String(n.BrowserModus),
		Version:                    String(n.Version),
		MajorVer:                   String(n.MajorVer),
		MinorVer:                   String(n.MinorVer),
		Platform:                   String(n.Platform),
		PlatformVersion:            String(n.PlatformVersion),
		PlatformDescription:        String(n.PlatformDescription),
		PlatformBits:               Int(n.PlatformBits),
		PlatformMaker:              String(n.PlatformMaker),
		Alpha:                      Bool(n.Alpha),
		Beta:                       Bool(n.Beta),
		Win16:                      Bool(n.Win16),
		Win32:                      Bool(n.Win32),
		Win64:                      Bool(n.Win64),
		Frames:                     Bool(n.Frames),
		Iframes:                    Bool(n.Iframes),
		Tables:                     Bool(n.Tables),
		Cookies:                    Bool(n.Cookies),
		BackgroundSounds:           Bool(n.BackgroundSounds),
		Javascript:                 Bool(n.Javascript),
		VBScript:                   Bool(n.VBScript),
		JavaApplets:                Bool(n.JavaApplets),
		ActiveXControls:            Bool(n.ActiveXControls),
		IsMobileDevice:             Bool(n.IsMobileDevice),
		IsTablet:                   Bool(n.IsTablet),
		IsSyndicationReader:        Bool(n.IsSyndicationReader),
		Crawler:                    Bool(n.Crawler),
		IsFake:                     Bool(n.IsFake),
		IsAnonymized:               Bool(n.IsAnonymized),
		IsModified:                 Bool(n.IsModified),
		CSSVersion:                 Int(n.CSSVersion),
		AolVersion:                 Int(n.AolVersion),
		DeviceName:                 String(n.DeviceName),
		DeviceMaker:                String(n.DeviceMaker),
		DeviceType:                 String(n.DeviceType),
		DevicePointingMethod:       String(n.DevicePointingMethod),
		DeviceCodeName:             String(n.DeviceCodeName),
		DeviceBrandName:            String(n.DeviceBrandName),
		RenderingEngineName:        String(n.RenderingEngineName),
		RenderingEngineVersion:     String(n.RenderingEngineVersion),
		RenderingEngineDescription: String(n.RenderingEngineDescription),
		RenderingEngineMaker:       String(n.RenderingEngineMaker),
	}
}
