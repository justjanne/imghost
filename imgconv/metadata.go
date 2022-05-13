package imgconv

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	AspectRatio          *Ratio             `json:"aspectRatio,omitempty"`
	Make                 string             `json:"make,omitempty"`
	Model                string             `json:"model,omitempty"`
	LensMake             string             `json:"lensMake,omitempty"`
	LensModel            string             `json:"lensModel,omitempty"`
	Software             string             `json:"software,omitempty"`
	Copyright            string             `json:"copyright,omitempty"`
	Description          string             `json:"description,omitempty"`
	CreatedAt            *time.Time         `json:"createdAt,omitempty"`
	DigitizedAt          *time.Time         `json:"digitizedAt,omitempty"`
	OriginalAt           *time.Time         `json:"originalAt,omitempty"`
	DigitalZoom          *Ratio             `json:"digitalZoom,omitempty"`
	Exposure             *Ratio             `json:"exposure,omitempty"`
	ExposureMode         ExposureMode       `json:"exposureMode,omitempty"`
	ExposureProgram      ExposureProgram    `json:"exposureProgram,omitempty"`
	ShutterSpeed         *Ratio             `json:"shutterSpeed,omitempty"`
	Aperture             *Ratio             `json:"aperture,omitempty"`
	Brightness           *Ratio             `json:"brightness,omitempty"`
	MaxAperture          *Ratio             `json:"maxAperture,omitempty"`
	Flash                *Flash             `json:"flash,omitempty"`
	FocalLength          *Ratio             `json:"focalLength,omitempty"`
	FocalLengthFF        *Ratio             `json:"focalLengthFF,omitempty"`
	IsoSpeedRating       *int64             `json:"isoSpeedRating,omitempty"`
	LightSource          LightSource        `json:"lightSource,omitempty"`
	MeteringMode         MeteringMode       `json:"meteringMode,omitempty"`
	WhiteBalance         WhiteBalance       `json:"whiteBalance,omitempty"`
	SceneMode            SceneMode          `json:"scene,omitempty"`
	ISO                  *int64             `json:"iso,omitempty"`
	Orientation          Orientation        `json:"orientation,omitempty"`
	Contrast             ContrastMode       `json:"contrast,omitempty"`
	Sharpness            SharpnessMode      `json:"sharpness,omitempty"`
	SubjectDistance      *Ratio             `json:"subjectDistance,omitempty"`
	SubjectDistanceRange DistanceRange      `json:"subjectDistanceRange,omitempty"`
	FileSource           FileSource         `json:"source,omitempty"`
	Saturation           *int64             `json:"saturation,omitempty"`
	SensorType           SensorType         `json:"sensor,omitempty"`
	LensSpecification    *LensSpecification `json:"lensSpecification,omitempty"`
	Location             *Location          `json:"location,omitempty"`
	Resolution           *Resolution        `json:"resolution,omitempty"`
}

func parseMetadata(wand *imagick.MagickWand) Metadata {
	keys := exifKeyMap(wand)

	get := func(key string) string {
		originalKey, ok := keys[key]
		if ok {
			return wand.GetImageProperty(originalKey)
		} else {
			return ""
		}
	}

	return Metadata{
		AspectRatio: &Ratio{
			int64(wand.GetImageWidth()),
			int64(wand.GetImageHeight()),
		},
		Make:                 strings.TrimSpace(get("Make")),
		Model:                strings.TrimSpace(get("Model")),
		LensMake:             strings.TrimSpace(get("LensMake")),
		LensModel:            strings.TrimSpace(get("LensModel")),
		Software:             strings.TrimSpace(get("Software")),
		Copyright:            strings.TrimSpace(get("Copyright")),
		Description:          strings.TrimSpace(get("ImageDescription")),
		CreatedAt:            parseTime(get("DateTime"), get("SubSecTime")),
		DigitizedAt:          parseTime(get("DateTimeDigitized"), get("SubSecTimeDigitized")),
		OriginalAt:           parseTime(get("DateTimeOriginal"), get("SubSecTimeOriginal")),
		DigitalZoom:          parseRatio(get("DigitalZoomRatio")),
		Exposure:             parseRatio(get("ExposureBiasValue")),
		ExposureMode:         parseExposureMode(parseNumber(get("ExposureMode"))),
		ExposureProgram:      parseExposureProgram(parseNumber(get("ExposureProgram"))),
		ShutterSpeed:         parseShutterSpeed(parseRatio(get("ExposureTime")), parseRatio(get("ShutterSpeedValue"))),
		Aperture:             parseRatio(get("FNumber")),
		Brightness:           parseRatio(get("BrightnessValue")),
		MaxAperture:          parseRatio(get("MaxApertureValue")),
		Flash:                parseFlash(parseNumber(get("Flash")), get("FlashEnergy")),
		FocalLength:          parseRatio(get("FocalLength")),
		FocalLengthFF:        parseRatio(get("FocalLengthIn35mmFilm")),
		ISO:                  parseNumber(get("PhotographicSensitivity")),
		LightSource:          parseLightSource(parseNumber(get("LightSource"))),
		MeteringMode:         parseMeteringMode(parseNumber(get("MeteringMode"))),
		Orientation:          parseOrientation(parseNumber(get("Orientation"))),
		WhiteBalance:         parseWhiteBalance(parseNumber(get("WhiteBalance"))),
		SceneMode:            parseSceneMode(parseNumber(get("SceneMode"))),
		Contrast:             parseContrastMode(parseNumber(get("Contrast"))),
		Sharpness:            parseSharpnessMode(parseNumber(get("Sharpness"))),
		SubjectDistance:      parseRatio(get("SubjectDistance")),
		SubjectDistanceRange: parseDistanceRange(parseNumber(get("SubjectDistanceRange"))),
		FileSource:           parseFileSource(parseNumber(get("FileSource"))),
		Saturation:           parseNumber(get("Saturation")),
		SensorType:           parseSensorType(parseNumber(get("SensingMethod"))),
		LensSpecification:    parseLensSpecification(get("LensSpecification")),
		Location: parseLocation(
			parseCoordinate(get("GPSLatitude"), get("GPSLatitudeRef")),
			parseCoordinate(get("GPSLongitude"), get("GPSLongitudeRef")),
		),
		Resolution: parseResolution(
			parseNumber(get("ResolutionUnit")),
			parseRatio(get("XResolution")),
			parseRatio(get("YResolution")),
		),
	}
}

func exifKeyMap(wand *imagick.MagickWand) map[string]string {
	metadata := make(map[string]string)
	for _, key := range wand.GetImageProperties("exif:*") {
		if strings.HasPrefix(key, "exif:thumbnail:") {
			continue
		}
		trimmedKey := strings.TrimPrefix(key, "exif:")
		metadata[trimmedKey] = key
	}
	return metadata
}

const ExifTime = "2006:01:02 15:04:05"

func parseTime(value string, subSec string) *time.Time {
	result, err := time.Parse(ExifTime, value)
	if err != nil {
		return nil
	} else {
		microseconds := parseNumber(subSec)
		if microseconds != nil {
			result = result.Add(time.Duration(*microseconds) * time.Microsecond)
		}
		return &result
	}
}

func parseNumber(value string) *int64 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil
	} else {
		return &result
	}
}

func parseShutterSpeed(exposure *Ratio, shutterSpeed *Ratio) *Ratio {
	if shutterSpeed != nil {
		return &Ratio{
			shutterSpeed.Denominator,
			shutterSpeed.Numerator,
		}
	}
	return exposure
}

type Flash struct {
	Available       bool            `json:"available"`
	Fired           bool            `json:"fired"`
	StrobeDetection StrobeDetection `json:"strobeDetection"`
	Mode            FlashMode       `json:"mode,omitempty"`
	RedEyeReduction bool            `json:"redEyeReduction"`
	Strength        *Ratio          `json:"strength,omitempty"`
}

type StrobeDetection struct {
	Available bool `json:"available"`
	Detected  bool `json:"detected"`
}

const (
	maskFired                    = 0x0001
	maskStrobeDetected           = 0x0002
	maskStrobeDetectionAvailable = 0x0004
	maskMode                     = 0x0003
	maskUnavailable              = 0x0020
	maskRedEye                   = 0x0040
)

func parseFlash(flash *int64, strength string) *Flash {
	if flash == nil {
		return nil
	}

	return &Flash{
		Available: *flash&maskUnavailable == 0,
		Fired:     *flash&maskFired != 0,
		StrobeDetection: StrobeDetection{
			Available: *flash&maskStrobeDetectionAvailable != 0,
			Detected:  *flash&maskStrobeDetected != 0,
		},
		Mode:            parseFlashMode((*flash >> 3) & maskMode),
		RedEyeReduction: *flash&maskRedEye != 0,
		Strength:        parseRatio(strength),
	}
}

type Resolution struct {
	X Ratio `json:"x"`
	Y Ratio `json:"y"`
}

func parseResolution(unit *int64, x *Ratio, y *Ratio) *Resolution {
	if unit == nil {
		defaultUnit := int64(2)
		unit = &defaultUnit
	}
	if x == nil || y == nil {
		return nil
	}
	return &Resolution{
		X: Ratio{
			x.Numerator,
			x.Denominator * *unit,
		}.reduce(),
		Y: Ratio{
			y.Numerator,
			y.Denominator * *unit,
		}.reduce(),
	}
}

type LensSpecification struct {
	WideFocalLength *Ratio `json:"wideFocalLength,omitempty"`
	WideAperture    *Ratio `json:"wideAperture,omitempty"`
	TeleFocalLength *Ratio `json:"teleFocalLength,omitempty"`
	TeleAperture    *Ratio `json:"teleAperture,omitempty"`
}

func parseLensSpecification(value string) *LensSpecification {
	split := strings.Split(value, ", ")
	if len(split) != 4 {
		return nil
	}
	return &LensSpecification{
		WideFocalLength: parseRatio(split[0]),
		TeleFocalLength: parseRatio(split[1]),
		WideAperture:    parseRatio(split[2]),
		TeleAperture:    parseRatio(split[3]),
	}
}

type Location struct {
	Longitude Coordinate `json:"longitude"`
	Latitude  Coordinate `json:"latitude"`
}

func parseLocation(longitude *Coordinate, latitude *Coordinate) *Location {
	if longitude == nil || latitude == nil {
		return nil
	}
	return &Location{*longitude, *latitude}
}

type Coordinate struct {
	Degree            Ratio             `json:"degree"`
	Minute            Ratio             `json:"minute"`
	Second            Ratio             `json:"second"`
	CardinalDirection CardinalDirection `json:"reference"`
}

func parseCoordinate(value string, reference string) *Coordinate {
	parts := strings.Split(value, ", ")
	if len(parts) != 3 {
		return nil
	}
	degree := parseRatio(parts[0])
	minute := parseRatio(parts[1])
	second := parseRatio(parts[2])
	direction := parseCardinalDirection(reference)
	if degree == nil || minute == nil || second == nil || direction == "" {
		return nil
	}
	return &Coordinate{
		*degree,
		*minute,
		*second,
		direction,
	}
}

type CardinalDirection string

const (
	CardinalDirectionNorth CardinalDirection = "N"
	CardinalDirectionWest  CardinalDirection = "W"
	CardinalDirectionSouth CardinalDirection = "S"
	CardinalDirectionEast  CardinalDirection = "E"
)

func parseCardinalDirection(value string) CardinalDirection {
	switch value {
	case "N":
		return CardinalDirectionNorth
	case "W":
		return CardinalDirectionWest
	case "S":
		return CardinalDirectionSouth
	case "E":
		return CardinalDirectionEast
	default:
		return ""
	}
}

type Ratio struct {
	Numerator   int64 `json:"num"`
	Denominator int64 `json:"den"`
}

func (ratio *Ratio) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%d/%d", ratio.Numerator, ratio.Denominator))
}

func (ratio *Ratio) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	split := strings.Split(raw, "/")
	if len(split) != 2 {
		fallback := parseNumber(raw)
		if fallback == nil {
			return fmt.Errorf("could not deserialize ratio: ratio is neither a ratio nor a plain number")
		} else {
			ratio.Numerator = *fallback
			ratio.Denominator = 1
			return nil
		}
	}
	numerator := parseNumber(split[0])
	if numerator == nil {
		return fmt.Errorf("could not deserialize ratio: numerator is not a valid number")
	}
	denominator := parseNumber(split[1])
	if denominator == nil {
		return fmt.Errorf("could not deserialize ratio: denominator is not a valid number")
	}

	ratio.Numerator = *numerator
	ratio.Denominator = *denominator
	return nil
}

func (ratio Ratio) reduce() Ratio {
	if ratio.Numerator > 0 && ratio.Denominator > 0 && ratio.Numerator%ratio.Denominator == 0 {
		ratio.Numerator = ratio.Numerator / ratio.Denominator
		ratio.Denominator = 1
	} else if ratio.Numerator == 0 {
		ratio.Numerator = 0
		ratio.Denominator = 1
	}
	return ratio
}

func parseRatio(value string) *Ratio {
	split := strings.Split(value, "/")
	if len(split) != 2 {
		return nil
	}
	numerator := parseNumber(split[0])
	denominator := parseNumber(split[1])
	if numerator == nil || denominator == nil {
		return nil
	}
	result := Ratio{*numerator, *denominator}.reduce()
	return &result
}

type SensorType string

const (
	SensorTypeOther                 SensorType = "other"
	SensorTypeSingleChipColorArea   SensorType = "single_chip_color_area"
	SensorTypeDualChipColorArea     SensorType = "dual_chip_color_area"
	SensorTypeTripleChipColorArea   SensorType = "triple_chip_color_area"
	SensorTypeColorSequentialArea   SensorType = "color_sequential_area"
	SensorTypeTrilinear             SensorType = "trilinear"
	SensorTypeColorSequentialLinear SensorType = "color_sequential_linear"

	valueSensorTypeOther                 = 1
	valueSensorTypeSingleChipColorArea   = 2
	valueSensorTypeDualChipColorArea     = 3
	valueSensorTypeTripleChipColorArea   = 4
	valueSensorTypeColorSequentialArea   = 5
	valueSensorTypeTrilinear             = 7
	valueSensorTypeColorSequentialLinear = 8
)

func parseSensorType(value *int64) SensorType {
	if value == nil {
		return ""
	}

	switch *value {
	case valueSensorTypeOther:
		return SensorTypeOther
	case valueSensorTypeSingleChipColorArea:
		return SensorTypeSingleChipColorArea
	case valueSensorTypeDualChipColorArea:
		return SensorTypeDualChipColorArea
	case valueSensorTypeTripleChipColorArea:
		return SensorTypeTripleChipColorArea
	case valueSensorTypeColorSequentialArea:
		return SensorTypeColorSequentialArea
	case valueSensorTypeTrilinear:
		return SensorTypeTrilinear
	case valueSensorTypeColorSequentialLinear:
		return SensorTypeColorSequentialLinear
	default:
		return ""
	}
}

type FileSource string

const (
	FileSourceOther               FileSource = "other"
	FileSourceTransmissiveScanner FileSource = "scanner_transmissive"
	FileSourceReflectiveScanner   FileSource = "scanner_reflective"
	FileSourceDigitalCamera       FileSource = "digital_camera"

	valueFileSourceOther               = 0
	valueFileSourceTransmissiveScanner = 1
	valueFileSourceReflectiveScanner   = 2
	valueFileSourceDigitalCamera       = 3
)

func parseFileSource(value *int64) FileSource {
	if value == nil {
		return ""
	}

	switch *value {
	case valueFileSourceOther:
		return FileSourceOther
	case valueFileSourceTransmissiveScanner:
		return FileSourceTransmissiveScanner
	case valueFileSourceReflectiveScanner:
		return FileSourceReflectiveScanner
	case valueFileSourceDigitalCamera:
		return FileSourceDigitalCamera
	default:
		return ""
	}
}

type Orientation string

const (
	OrientationTopLeft     Orientation = "top_left"
	OrientationTopRight    Orientation = "top_right"
	OrientationBottomRight Orientation = "bottom_right"
	OrientationBottomLeft  Orientation = "bottom_left"
	OrientationLeftTop     Orientation = "left_top"
	OrientationRightTop    Orientation = "right_top"
	OrientationRightBottom Orientation = "right_bottom"
	OrientationLeftBottom  Orientation = "left_bottom"

	valueOrientationTopLeft     = 1
	valueOrientationTopRight    = 2
	valueOrientationBottomRight = 3
	valueOrientationBottomLeft  = 4
	valueOrientationLeftTop     = 5
	valueOrientationRightTop    = 6
	valueOrientationRightBottom = 7
	valueOrientationLeftBottom  = 8
)

func parseOrientation(value *int64) Orientation {
	if value == nil {
		return ""
	}

	switch *value {
	case valueOrientationTopLeft:
		return OrientationTopLeft
	case valueOrientationTopRight:
		return OrientationTopRight
	case valueOrientationBottomRight:
		return OrientationBottomRight
	case valueOrientationBottomLeft:
		return OrientationBottomLeft
	case valueOrientationLeftTop:
		return OrientationLeftTop
	case valueOrientationRightTop:
		return OrientationRightTop
	case valueOrientationRightBottom:
		return OrientationRightBottom
	case valueOrientationLeftBottom:
		return OrientationLeftBottom
	default:
		return ""
	}
}

type ContrastMode string

const (
	ContrastModeNormal ContrastMode = "normal"
	ContrastModeSoft   ContrastMode = "soft"
	ContrastModeHard   ContrastMode = "hard"

	valueContrastModeNormal = 0
	valueContrastModeSoft   = 1
	valueContrastModeHard   = 2
)

func parseContrastMode(value *int64) ContrastMode {
	if value == nil {
		return ""
	}

	switch *value {
	case valueContrastModeNormal:
		return ContrastModeNormal
	case valueContrastModeSoft:
		return ContrastModeSoft
	case valueContrastModeHard:
		return ContrastModeHard
	default:
		return ""
	}
}

type DistanceRange string

const (
	DistanceRangeMacro   DistanceRange = "macro"
	DistanceRangeClose   DistanceRange = "close"
	DistanceRangeDistant DistanceRange = "distant"

	valueDistanceRangeMacro   = 1
	valueDistanceRangeClose   = 2
	valueDistanceRangeDistant = 3
)

func parseDistanceRange(value *int64) DistanceRange {
	if value == nil {
		return ""
	}

	switch *value {
	case valueDistanceRangeMacro:
		return DistanceRangeMacro
	case valueDistanceRangeClose:
		return DistanceRangeClose
	case valueDistanceRangeDistant:
		return DistanceRangeDistant
	default:
		return ""
	}
}

type ExposureMode string

const (
	ExposureModeAuto    ExposureMode = "auto"
	ExposureModeManual  ExposureMode = "manual"
	ExposureModeBracket ExposureMode = "bracket"

	valueExposureModeAuto    = 0
	valueExposureModeManual  = 1
	valueExposureModeBracket = 2
)

func parseExposureMode(value *int64) ExposureMode {
	if value == nil {
		return ""
	}

	switch *value {
	case valueExposureModeAuto:
		return ExposureModeAuto
	case valueExposureModeManual:
		return ExposureModeManual
	case valueExposureModeBracket:
		return ExposureModeBracket
	default:
		return ""
	}
}

type ExposureProgram string

const (
	ExposureProgramManual           ExposureProgram = "manual"
	ExposureProgramNormal           ExposureProgram = "normal"
	ExposureProgramAperturePriority ExposureProgram = "aperture_priority"
	ExposureProgramShutterPriority  ExposureProgram = "shutter_priority"
	ExposureProgramCreative         ExposureProgram = "creative"
	ExposureProgramAction           ExposureProgram = "action"
	ExposureProgramPortrait         ExposureProgram = "portrait"
	ExposureProgramLandscape        ExposureProgram = "landscape"

	valueExposureProgramManual           = 1
	valueExposureProgramNormal           = 2
	valueExposureProgramAperturePriority = 3
	valueExposureProgramShutterPriority  = 4
	valueExposureProgramCreative         = 5
	valueExposureProgramAction           = 6
	valueExposureProgramPortrait         = 7
	valueExposureProgramLandscape        = 8
)

func parseExposureProgram(value *int64) ExposureProgram {
	if value == nil {
		return ""
	}

	switch *value {
	case valueExposureProgramManual:
		return ExposureProgramManual
	case valueExposureProgramNormal:
		return ExposureProgramNormal
	case valueExposureProgramAperturePriority:
		return ExposureProgramAperturePriority
	case valueExposureProgramShutterPriority:
		return ExposureProgramShutterPriority
	case valueExposureProgramCreative:
		return ExposureProgramCreative
	case valueExposureProgramAction:
		return ExposureProgramAction
	case valueExposureProgramPortrait:
		return ExposureProgramPortrait
	case valueExposureProgramLandscape:
		return ExposureProgramLandscape
	default:
		return ""
	}
}

type FlashMode string

const (
	FlashModeAlwaysOn  FlashMode = "always_on"
	FlashModeAlwaysOff FlashMode = "always_off"
	FlashModeAuto      FlashMode = "auto"

	valueFlashModeAlwaysOn  = 1
	valueFlashModeAlwaysOff = 2
	valueFlashModeAuto      = 3
)

func parseFlashMode(value int64) FlashMode {
	switch value {
	case valueFlashModeAlwaysOn:
		return FlashModeAlwaysOn
	case valueFlashModeAlwaysOff:
		return FlashModeAlwaysOff
	case valueFlashModeAuto:
		return FlashModeAuto
	default:
		return ""
	}
}

type LightSource string

const (
	LightSourceDaylight      LightSource = "daylight"
	LightSourceFluorescent   LightSource = "fluorescent"
	LightSourceIncandescent  LightSource = "incandescent"
	LightSourceFlash         LightSource = "flash"
	LightSourceFineWeather   LightSource = "weather_fine"
	LightSourceCloudyWeather LightSource = "weather_cloudy"
	LightSourceShade         LightSource = "shade"
	LightSource6400K         LightSource = "6400K"
	LightSource5000K         LightSource = "5000K"
	LightSource4200K         LightSource = "4200K"
	LightSource3450K         LightSource = "3450K"
	LightSourceStandardA     LightSource = "standard_a"
	LightSourceStandardB     LightSource = "standard_b"
	LightSourceStandardC     LightSource = "standard_c"
	LightSourceD55           LightSource = "D55"
	LightSourceD65           LightSource = "D65"
	LightSourceD75           LightSource = "D75"
	LightSourceD50           LightSource = "D50"
	LightSourceIsoStudio     LightSource = "iso_studio"

	valueLightSourceDaylight      = 1
	valueLightSourceFluorescent   = 2
	valueLightSourceIncandescent  = 3
	valueLightSourceFlash         = 4
	valueLightSourceFineWeather   = 9
	valueLightSourceCloudyWeather = 10
	valueLightSourceShade         = 11
	valueLightSource6400K         = 12
	valueLightSource5000K         = 13
	valueLightSource4200K         = 14
	valueLightSource3450K         = 15
	valueLightSourceStandardA     = 17
	valueLightSourceStandardB     = 18
	valueLightSourceStandardC     = 19
	valueLightSourceD55           = 20
	valueLightSourceD65           = 21
	valueLightSourceD75           = 22
	valueLightSourceD50           = 23
	valueLightSourceIsoStudio     = 24
)

func parseLightSource(value *int64) LightSource {
	if value == nil {
		return ""
	}

	switch *value {
	case valueLightSourceDaylight:
		return LightSourceDaylight
	case valueLightSourceFluorescent:
		return LightSourceFluorescent
	case valueLightSourceIncandescent:
		return LightSourceIncandescent
	case valueLightSourceFlash:
		return LightSourceFlash
	case valueLightSourceFineWeather:
		return LightSourceFineWeather
	case valueLightSourceCloudyWeather:
		return LightSourceCloudyWeather
	case valueLightSourceShade:
		return LightSourceShade
	case valueLightSource6400K:
		return LightSource6400K
	case valueLightSource5000K:
		return LightSource5000K
	case valueLightSource4200K:
		return LightSource4200K
	case valueLightSource3450K:
		return LightSource3450K
	case valueLightSourceStandardA:
		return LightSourceStandardA
	case valueLightSourceStandardB:
		return LightSourceStandardB
	case valueLightSourceStandardC:
		return LightSourceStandardC
	case valueLightSourceD55:
		return LightSourceD55
	case valueLightSourceD65:
		return LightSourceD65
	case valueLightSourceD75:
		return LightSourceD75
	case valueLightSourceD50:
		return LightSourceD50
	case valueLightSourceIsoStudio:
		return LightSourceIsoStudio
	default:
		return ""
	}
}

type MeteringMode string

const (
	MeteringModeAverage               MeteringMode = "average"
	MeteringModeCenterWeightedAverage MeteringMode = "center_weighted_average"
	MeteringModeSpot                  MeteringMode = "spot"
	MeteringModeMultiSpot             MeteringMode = "multi_sport"
	MeteringModePattern               MeteringMode = "pattern"
	MeteringModePartial               MeteringMode = "partial"

	valueMeteringModeAverage               = 1
	valueMeteringModeCenterWeightedAverage = 2
	valueMeteringModeSpot                  = 3
	valueMeteringModeMultiSpot             = 4
	valueMeteringModePattern               = 5
	valueMeteringModePartial               = 6
)

func parseMeteringMode(value *int64) MeteringMode {
	if value == nil {
		return ""
	}

	switch *value {
	case valueMeteringModeAverage:
		return MeteringModeAverage
	case valueMeteringModeCenterWeightedAverage:
		return MeteringModeCenterWeightedAverage
	case valueMeteringModeSpot:
		return MeteringModeSpot
	case valueMeteringModeMultiSpot:
		return MeteringModeMultiSpot
	case valueMeteringModePattern:
		return MeteringModePattern
	case valueMeteringModePartial:
		return MeteringModePartial
	default:
		return ""
	}
}

type SceneMode string

const (
	SceneModeStandard   SceneMode = "standard"
	SceneModeLandscape  SceneMode = "landscape"
	SceneModePortrait   SceneMode = "portrait"
	SceneModeNightScene SceneMode = "night"

	valueSceneModeStandard   = 0
	valueSceneModeLandscape  = 1
	valueSceneModePortrait   = 2
	valueSceneModeNightScene = 3
)

func parseSceneMode(value *int64) SceneMode {
	if value == nil {
		return ""
	}

	switch *value {
	case valueSceneModeStandard:
		return SceneModeStandard
	case valueSceneModeLandscape:
		return SceneModeLandscape
	case valueSceneModePortrait:
		return SceneModePortrait
	case valueSceneModeNightScene:
		return SceneModeNightScene
	default:
		return ""
	}
}

type SharpnessMode string

const (
	SharpnessModeNormal SharpnessMode = "normal"
	SharpnessModeSoft   SharpnessMode = "soft"
	SharpnessModeHard   SharpnessMode = "hard"

	valueSharpnessModeNormal = 0
	valueSharpnessModeSoft   = 1
	valueSharpnessModeHard   = 2
)

func parseSharpnessMode(value *int64) SharpnessMode {
	if value == nil {
		return ""
	}

	switch *value {
	case valueSharpnessModeNormal:
		return SharpnessModeNormal
	case valueSharpnessModeSoft:
		return SharpnessModeSoft
	case valueSharpnessModeHard:
		return SharpnessModeHard
	default:
		return ""
	}
}

type WhiteBalance string

const (
	WhiteBalanceAuto   WhiteBalance = "auto"
	WhiteBalanceManual WhiteBalance = "manual"

	valueWhiteBalanceAuto   = 0
	valueWhiteBalanceManual = 1
)

func parseWhiteBalance(value *int64) WhiteBalance {
	if value == nil {
		return ""
	}

	switch *value {
	case valueWhiteBalanceAuto:
		return WhiteBalanceAuto
	case valueWhiteBalanceManual:
		return WhiteBalanceManual
	default:
		return ""
	}
}
