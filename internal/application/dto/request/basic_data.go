package request

// BasicDataRequest represents the request body for creating/updating BasicData.
type BasicDataRequest struct {
	FirstName         string  `json:"firstName" validate:"required"`
	OthersName        *string `json:"othersName,omitempty"`
	FirstSurname      string  `json:"firstSurName" validate:"required"`
	OthersSurname     *string `json:"othersSurname,omitempty"`
	DateBirth         string  `json:"dateBirth" validate:"required"`
	Located           *string `json:"located,omitempty"`
	LocatedEng        *string `json:"locatedEng,omitempty"`
	StartWorkingDate  *string `json:"startWorkingDate,omitempty"`
	Greeting          *string `json:"greeting,omitempty"`
	GreetingEng       *string `json:"greetingEng,omitempty"`
	Email             string  `json:"email" validate:"required,email"`
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
}
