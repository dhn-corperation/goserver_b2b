package kakaojson

type Friendtalk struct {
	Message_type  string     `json:"message_type"`
	Serial_number string     `json:"serial_number"`
	Sender_key    string     `json:"sender_key"`
	Phone_number  string     `json:"phone_number,omitempty"`
	App_user_id   string     `json:"app_user_id,omitempty"`
	User_key      string     `json:"user_key,omitempty"`
	Message       string     `json:"message"`
	Ad_flag       string     `json:"ad_flag,omitempty"`
	Attachment    *FTAttachment `json:"attachment,omitempty"`
	Header        string     `json:"header,omitempty"`
	Carousel      *FTCarousel `json:"carousel,omitempty"`
}

type FTAButton struct {
		Name          string `json:"name,omitempty"`
		Type          string `json:"type,omitempty"`
		SchemeAndroid string `json:"scheme_android,omitempty"`
		URLPc         string `json:"url_pc,omitempty"`
		URLMobile     string `json:"url_mobile,omitempty"`
		SchemeIos     string `json:"scheme_ios,omitempty"`
	}

type FTAImage struct {
		ImgURL  string `json:"img_url,omitempty"`
		ImgLink string `json:"img_link,omitempty"`
	}

type FTAItem struct {
		List *[]FTAIList `json:"list,omitempty"`
	}

type FTAIList struct {
			Title         string `json:"title,omitempty"`
			ImgURL        string `json:"img_url,omitempty"`
			SchemeAndroid string `json:"scheme_android,omitempty"`
			URLPc         string `json:"url_pc,omitempty"`
			URLMobile     string `json:"url_mobile,omitempty"`
			SchemeIos     string `json:"scheme_ios,omitempty"`
		}

type FTACoupon struct {
		Title         string `json:"title,omitempty"`
		Description   string `json:"description,omitempty"`
		SchemeAndroid string `json:"scheme_android,omitempty"`
		URLPc         string `json:"url_pc,omitempty"`
		URLMobile     string `json:"url_mobile,omitempty"`
		SchemeIos     string `json:"scheme_ios,omitempty"`
	}


type FTACommerce struct {
		Title         string `json:"title,omitempty"`
		RegularPrice  string `json:"regular_price,omitempty"`
		DiscountPrice string `json:"discount_price,omitempty"`
		DiscountRate  string `json:"discount_rate,omitempty"`
		DiscountFixed string `json:"discount_fixed,omitempty"`
	}


type FTAVideo struct {
		VideoURL     string `json:"video_url,omitempty"`
		ThumbnailURL string `json:"thumbnail_url,omitempty"`
	}


		
type FTAttachment struct {
	Button *[]FTAButton `json:"button,omitempty"`
	Image    *FTAImage  `json:"image,omitempty"`
	Item     *FTAItem   `json:"item,omitempty"`
	Coupon   *FTACoupon `json:"coupon,omitempty"`
	Commerce *FTACommerce `json:"commerce,omitempty"`
	Video    *FTAVideo  `json:"video,omitempty"`
}

type FTCHead struct {
		Header        string `json:"header,omitempty"`
		Content       string `json:"content,omitempty"`
		ImageURL      string `json:"image_url,omitempty"`
		URLMobile     string `json:"url_mobile,omitempty"`
		URLPc         string `json:"url_pc,omitempty"`
		SchemeAndroid string `json:"scheme_android,omitempty"`
		SchemeIos     string `json:"scheme_ios,omitempty"`
	}


type FTCLAttachment        struct {
			Button    *FTAButton   `json:"button,omitempty"`
			Image     *FTAImage    `json:"image,omitempty"`
			Coupon    *FTACoupon   `json:"coupon,omitempty"`
			Commerce  *FTACommerce `json:"commerce,omitempty"`
		}

type FTCList struct {
		Header            string `json:"header,omitempty"`
		Message           string `json:"message,omitempty"`
		AdditionalContent string `json:"additional_content,omitempty"`
		Attachment       *FTCLAttachment `json:"attachment,omitempty"`
	}


type FTCTail struct {
		URLPc         string `json:"url_pc,omitempty"`
		URLMobile     string `json:"url_mobile,omitempty"`
		SchemeIos     string `json:"scheme_ios,omitempty"`
		SchemeAndroid string `json:"scheme_android,omitempty"`
	}

type FTCarousel struct {
	Head     *FTCHead    `json:"head,omitempty"`
	List     *FTCList    `json:"list,omitempty"`
	Tail     *FTCTail    `json:"tail,omitempty"`
}


type KakaoResponse struct {
	Code        string
	Received_at string
	Message     string
}

type PollingResponse struct {
	Code         string
	Response_id  int
	Response     PResponse
	Responsed_at string
	Message      string
}

type PResponse struct {
	Success []PResult
	Fail    []PResult
}

type PResult struct {
	Serial_number string
	Status        string
	Received_at   string
}
