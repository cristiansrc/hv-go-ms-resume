package response

// BasicDataResponse represents the API response for BasicData.
type BasicDataResponse struct {
	ID                int64   `json:"id"`
	FirstName         string  `json:"firstName"`
	OthersName        *string `json:"othersName,omitempty"`
	FirstSurname      string  `json:"firstSurName"`
	OthersSurname     *string `json:"othersSurname,omitempty"`
	DateBirth         string  `json:"dateBirth"`
	Located           *string `json:"located,omitempty"`
	LocatedEng        *string `json:"locatedEng,omitempty"`
	StartWorkingDate  *string `json:"startWorkingDate,omitempty"`
	Greeting          *string `json:"greeting,omitempty"`
	GreetingEng       *string `json:"greetingEng,omitempty"`
	Email             string  `json:"email"`
	Instagram         *string `json:"instagram,omitempty"`
	Linkedin          *string `json:"linkedin,omitempty"`
	X                 *string `json:"x,omitempty"`
	Github            *string `json:"github,omitempty"`
	Description       *string `json:"description,omitempty"`
	DescriptionEng    *string `json:"descriptionEng,omitempty"`
	DescriptionPdf    *string `json:"descriptionPdf,omitempty"`
	DescriptionPdfEng *string `json:"descriptionPdfEng,omitempty"`
	Wrapper           *string `json:"wrapper,omitempty"`
	WrapperEng        *string `json:"wrapperEng,omitempty"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}
