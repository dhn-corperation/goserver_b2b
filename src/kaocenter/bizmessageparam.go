package kaocenter

// 다이렉트 메시지 기본형 API -> 보류
type Direct_send_basic struct {
	Sender_key        string               `json:"sender_key" binding:"required"`
	Template_code     string               `json:"template_code" binding:"required"`
	Phone_number      string               `json:"phone_number,omitempty"`
	App_user_id       interface{}          `json:"app_user_id,omitempty"`
	Push_alarm        string               `json:"push_alarm,omitempty"`
	Message_variable  Message_variables    `json:"message_variable,omitempty"`
	Button_variable   Button_variables     `json:"button_variable,omitempty"`
	Coupon_variable   Coupon_variables     `json:"coupon_variable,omitempty"`
	Image_variable    []Image_variables    `json:"image_variable,omitempty"`
	Video_variable    Video_variables      `json:"video_variable,omitempty"`
	Commerce_variable Commerce_variables   `json:"commerce_variable,omitempty"`
	Carousel_variable []Carousel_variables `json:"carousel_variable,omitempty"`
}

type Message_variables struct {
}

type Button_variables struct {
}

type Coupon_variables struct {
}

type Image_variables struct {
	ImgURL  string `json:"img_url"`
	ImgLink string `json:"img_link,omitempty"`
}

type Video_variables struct {
}

type Commerce_variables struct {
}

type Carousel_variables struct {
}
