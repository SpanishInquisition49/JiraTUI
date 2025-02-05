package app

// TODO: add the size in the components
type Size struct {
	width         uint8 // The width of the terminal
	height        uint8 // The height of the terminal
	widthPercent  uint8 // The perentage of the width used by the component
	heightPercent uint8 // The perentage of the height used by the component
}

func (s *Size) SetWidth(width uint8) {
	s.width = width
}

func (s *Size) SetHeight(height uint8) {
	s.height = height
}

func (s *Size) SetDimensions(width uint8, height uint8) {
	s.SetWidth(width)
	s.SetHeight(height)
}

func (s *Size) SetWidthPercent(widthPercent uint8) {
	s.widthPercent = widthPercent
}

func (s *Size) SetHeightPercent(heightPercent uint8) {
	s.heightPercent = heightPercent
}

func (s *Size) SetDimensionsPercent(widthPercent uint8, heightPercent uint8) {
	s.SetWidthPercent(widthPercent)
	s.SetHeightPercent(heightPercent)
}

/**
  * GetWidth returns the width of to use for the component based on the given percentage
  * @return uint8 - The width of the component
  */
func (s *Size) GetWidth() uint8 {
  return s.width * s.widthPercent / 100
}

/**
  * GetHeight returns the height of to use for the component based on the given percentage
  * @return uint8 - The height of the component
  */
func (s *Size) GetHeight() uint8 {
  return s.height * s.heightPercent / 100
}

/**
  * GetDimensions returns the width and height of to use for the component based on the given percentage
  * @return uint8, uint8 - The width and height of the component
  */
func (s *Size) GetDimensions() (uint8, uint8) {
  return s.GetWidth(), s.GetHeight()
}
