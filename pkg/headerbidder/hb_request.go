package headerbidder

// HBRequest
//
// A header-bidder listener request is structure that is sent to kafka as a event.
//
// swagger:model
type HBRequest struct {
	// This is the ID for a HB Request
	//
	// required: true
	// example: 24234dfa-ef05-45fc-9dd8-7f37898294b6
	Id string `json:"id"`
	// swagger:allOf
	Imp []Imp `json:"imp"`
	// swagger:allOf
	App App `json:"app"`
	// swagger:allOf
	Device Device `json:"device"`
	// This object contains information known or derived about the human user of the device (i.e., the
	// audience for advertising). The user id is an exchange artifact and may be subject to rotation or other
	// privacy policies. However, this user ID must be stable long enough to serve reasonably as the basis for
	// frequency capping and retargeting.
	User struct {
		// Exchange-specific ID for the user. At least one of id or buyeruid is recommended.
		//
		// required: false
		// example: 575883e2-8fb1-4147-b6d3-0a5896b94e54
		Id string `json:"id"`
	} `json:"user"`
	// Auction type, where 1 = First Price, 2 = Second Price Plus.
	// Exchange-specific auction types can be defined using values
	// greater than 500.
	//
	// required: false
	// example: 1
	At int `json:"at"`
	// This parameter allows the timeout limit for SpotX to respond to be changed to best fit your desired user experience.
	// Values are in milliseconds with a max value of 650 being allowed.
	//
	// required: false
	// example: 2800
	Tmax int `json:"tmax"`
	// This object describes the nature and behavior of the entity that is the source of the bid request
	// upstream from the exchange. The primary purpose of this object is to define post-auction or upstream
	// decisioning when the exchange itself does not control the final decision. A common example of this is
	// header bidding, but it can also apply to upstream server entities such as another RTB exchange, a
	// mediation platform, or an ad server combines direct campaigns with 3rd party demand in decisioning.
	Source struct {
		// Entity responsible for the final impression sale decision, where 0 = exchange, 1 = upstream source.
		//
		// required: false
		// example: 0
		Fd int `json:"fd"`
		// Payment ID chain string containing embedded syntax described in the TAG Payment ID Protocol v1.0.
		//
		// required: false
		// example:""
		Pchain string `json:"pchain"`
	} `json:"source"`
	// Extra Data for Samsung Ads
	Ext struct {
		// Samsung HB Tag
		// HB TagID define in patner ad server.
		//
		// required: false
		// example: 1234
		Sam_HB_Tag string `json:"sam_hb_tag"`
		// Samsung App ID
		// HB App ID define in patner ad server.
		//
		// required: false
		// example: TEST-App-Id01
		Sam_APP_Id string `json:"sam_app_id"`
		// Samsung Session ID
		// Session ID generated on App open.
		//
		// required: false
		// example: 888888e2-621f-4f6a-844c-d21aed450c2a
		Sam_Session_ID string `json:"sam_session_id"`
	} `json:"ext"`
}

// This is impression of user in HB Request
type Imp struct {
	// Impression ID for User's Impression on a Ad.
	//
	// required: false
	// example: 1
	Id string `json:"id"`
	// swagger:allOf
	Video Video `json:"video"`
	// Identifier for specific ad placement or ad tag that was used to
	// initiate the auction. This can be useful for debugging of any
	// issues, or for optimization by the buyer.
	//
	// required: false
	// example: 666310
	Tagid string `json:"tagid"`
	// Flag to indicate if the impression requires secure HTTPS URL
	// creative assets and markup, where 0 = non-secure, 1 = secure.
	// If omitted, the secure state is unknown, but non-secure HTTP
	// support can be assumed.
	//
	// required: false
	// example: 1
	Secure int `json:"secure"`
	// Minimum bid for this impression expressed in CPM.
	//
	// required: false
	// example: 0.01
	BidFloor float64 `json:"bidfloor"`
}

type Video struct {
	// Content MIME types supported (e.g., “video/x-ms-wmv”, “video/mp4”).
	//
	// required: false
	// example: [ "video/mp4", "video/ogg","video/webm"]
	Mimes []string `json:"mimes"`
	// Minimum video ad duration in seconds.
	//
	// required: false
	// example: 1
	Minduration int `json:"minduration"`
	// Maximum video ad duration in seconds.
	//
	// required: false
	// example: 300
	Maxduration int `json:"maxduration"`
	// Array of supported video protocols.
	//
	// 1 - VAST 1.0
	// 2 - VAST 2.0
	// 3 - VAST 3.0
	// 4 - VAST 1.0 Wrapper
	// 5 - VAST 2.0 Wrapper
	// 6 - VAST 3.0 Wrapper
	// 7 - VAST 4.0
	// 8 - VAST 4.0 Wrapper
	// 9 - DAAST 1.0
	// 10 - DAAST 1.0 Wrapper
	//
	// required: false
	// example: [1, 2, 3, 4, 5, 6]
	Protocols []int `json:"protocols"`
	// Width of the video player in device independent pixels (DIPS).
	//
	// required: false
	// example: 1920
	W int `json:"w"`
	// Height of the video player in device independent pixels (DIPS).
	//
	// required: false
	// example: 1080
	H int `json:"h"`
	// Indicates the start delay in seconds for pre-roll, mid-roll, or post-roll ad placements.
	//
	// required: false
	// example: -1
	Startdelay int `json:"startdelay"`
	// Indicates if the impression must be linear, nonlinear, etc. If none specified, assume all are allowed.
	//
	// required: false
	// example: 1
	Linearity int `json:"linearity"`
	// If multiple ad impressions are offered in the same bid request,
	// the sequence number will allow for the coordinated delivery
	// of multiple creatives.
	//
	// required: false
	// example: 1
	Sequence int `json:"sequence"`
	// Minimum bit rate in Kbps.
	//
	// required: false
	// example: 1
	Minbitrate int `json:"minbitrate"`
	// Maximum bit rate in Kbps.
	//
	// required: false
	// example: 280000
	Maxbitrate int `json:"maxbitrate"`
}

// This object should be included if the ad supported content is a non-browser application
// (typically in mobile) as opposed to a website.
type App struct {
	// Exchange-specific app ID.
	//
	// required: false
	// example: 666310
	Id string `json:"id"`
	// App name (may be aliased at the publisher’s request).
	//
	// required: false
	// example: test_app
	Name string `json:"name"`
	// A platform-specific application identifier intended to be
	// unique to the app and independent of the exchange. On
	// Android, this should be a bundle or package name (e.g.,
	// com.foo.mygame). On iOS, it is typically a numeric ID.
	//
	//required: false
	//example: test
	Bundle string `json:"bundle"`
	// Domain of the app (e.g., “mygame.foo.com”).
	//
	// required: false
	// example: test.com
	Domain string `json:"domain"`
	// This object describes the publisher of the media in which the ad will be displayed.
	// The publisher is typically the seller in an OpenRTB transaction.
	Publisher struct {
		// Exchange-specific publisher ID.
		//
		// required: false
		// example: 1681
		Id string `json:"id"`
	} `json:"publisher"`
}

// This object encapsulates various methods for specifying a geographic location. When subordinate to a
// Device object, it indicates the location of the device which can also be interpreted as the user’s current
// location. When subordinate to a User object, it indicates the location of the user’s home base (i.e., not
// necessarily their current location).
type Geo struct {
	// Latitude from -90.0 to +90.0, where negative is south.
	//
	// required: false
	// example: 37.751
	Lat float64 `json:"lat"`
	// Longitude from -180.0 to +180.0, where negative is west.
	//
	// required: false
	// example: -97.822
	Lon float64 `json:"lon"`
	// Source of location data; recommended when passing lat/lon.
	//
	// 1 - GPS/Location Services
	// 2 - IP Address
	// 3 - User provided (e.g., registration data)
	//
	// required: false
	// example: 2
	Type string `json:"type"`
	// Country code using ISO-3166-1-alpha-3.
	//
	// required: false
	// example: USA
	Country string `json:"country"`
	// Service or provider used to determine geolocation from IP address if applicable (i.e., type = 2).
	//
	// 1 - ip2location
	// 2 - Neustar (Quova)
	// 3 - MaxMind
	// 4 - NetAcuity (Digital Element)
	//
	// required: false
	// example: 3
	Ipservice int `json:"ipservice"`
}

type Device struct {
	// Browser user agent string.
	//
	// required: false
	// example: Mozilla/5.0 (SMART-TV; LINUX; Tizen 3.0) AppleWebKit/538.1 (KHTML, like Gecko) Version/3.0 TV Safari/538.1
	Ua string `json:"ua"`
	// swagger:allOf
	Geo Geo `json:"geo"`
	// Standard “Do Not Track” flag as set in the header by the
	// browser, where 0 = tracking is unrestricted, 1 = do not track.
	//
	// required: false
	// example: 0
	Dnt int `json:"dnt"`
	// “Limit Ad Tracking” signal commercially endorsed (e.g., iOS,
	// Android), where 0 = tracking is unrestricted, 1 = tracking must
	// be limited per commercial guidelines.
	//
	//required: false
	// example: 0
	Lmt int `json:"lmt"`
	// IPv4 address closest to device.
	//
	// required: false
	// example: 11.22.33.44
	Ip string `json:"ip"`
	//The general type of device.
	//
	// 1 - Mobile/Tablet Version : 2.0
	// 2 - Personal Computer Version : 2.0
	// 3 - Connected TV Version : 2.0
	// 4 - Phone New for Version : 2.2
	// 5 - Tablet New for Version : 2.2
	// 6 - Connected Device New for Version : 2.2
	// 7 - Set Top Box New for Version : 2.2
	//
	// required: false
	// example: 3
	Devicetype int `json:"devicetype"`
	// Device make (e.g., “Apple”).
	//
	// required: false
	// example: Samsung
	Make string `json:"make"`
	// Device model (e.g., “iPhone”).
	//
	// required: false
	// example: Tizen TV 2017
	Model string `json:"model"`
	// Device operating system (e.g., “iOS”).
	//
	// required: false
	// example: Tizen
	Os string `json:"os"`
	// Browser language using ISO-639-1-alpha-2.
	//
	// required: false
	// example: en
	Language string `json:"language"`
	// ID sanctioned for advertiser use in the clear (i.e., not hashed).
	//
	// required: false
	// example: 11112222-3333-4444-5555-666677778888
	Ifa string `json:"ifa"`
}

// swagger:parameters create-header-bidder
type _ struct {
	// HB Request Body.
	//
	// in:body
	// required: true
	Body HBRequest
}
