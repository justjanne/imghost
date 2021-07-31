# EXIF metadata format

* `Make`  
  Manufacturer
* `Model`  
  Device
* `DateTime`
* `DateTimeDigitized`
* `DateTimeOriginal`
* `DigitalZoomRatio`  
  Zoom  
  - 0 unused
* `ExposureBiasValue`  
  Exposure  
  -100.0 to +100.0
* `ExposureMode`  
  Exposure Mode  
  - 0 = auto
  - 1 = manual
  - 2 = bracket
* `ExposureProgram`  
  Exposure Program
  - 0 undefined 
  - 1 manual
  - 2 normal
  - 3 aperture priority 
  - 4 shutter priority
  - 5 creative (depth of field priority)
  - 6 action (fast shutter priority)
  - 7 portrait (object separation priority)
  - 8 landscape (background in focus priority)
* `ExposureTime`  
  Shutter in seconds
* `FNumber`
  Aperture
* `Flash`  
  Bitfield of Flash metadata (bits counted from LSB to MSB)
  - bit 0: flashFired
  - bit 1: flashStrobeDetectionAvailable
  - bit 2: flashStrobeDetected
  - bit 3-4
    - 0 undefined 
    - 1 always on
    - 2 always off
    - 3 auto 
  - bit 5: flashAvailable
  - bit 6: redEyeReductionAvailable
* `FlashEnergy`  
  Strobe energy in BCPS
* `FocalLength`  
  Focal Length in mm
* `FocalLengthIn35mmFilm`
  Focal Length in mm compared to a 35mm Film equivalent
* `ISOSpeedRatings`
  ISO exposure/speed rating
* `LightSource`  
  Type of lightsource for white balance
  - 0 undefined
  - 1 daylight
  - 2 fluorescent
  - 3 tungsten / incandescent
  - 4 flash
  - 9 fine weather
  - 10 cloudy weather
  - 11 shade
  - 12 daylight fluorescent 5700-7100K
  - 13 day white fluorescent 4600-5400K
  - 14 cool white fluorescent 3900-4500K
  - 15 white fluorescent 3200-3700K
  - 17 standard light A
  - 18 standard light B
  - 19 standard light C
  - 20 D55
  - 21 D65
  - 22 D75
  - 23 D50
  - 24 ISO studio tungsten
  - 255 other
* `MeteringMode`  
  exposure metering mode
  - 0 undefined
  - 1 average
  - 2 center weighted average
  - 3 spot
  - 4 multispot
  - 5 pattern
  - 6 partial
  - 255 other
* `WhiteBalance`  
  White balance
  - 0 auto
  - 1 manual
* `SceneCaptureType`  
  Scene Mode
  - 0 standard
  - 1 landscape
  - 2 portrait
  - 3 night scene
* `Contrast`  Contrast Processing
  - 0 normal
  - 1 soft
  - 2 hard
* `Sharpness`  
  Sharpness Processing
  - 0 normal
  - 1 soft
  - 2 hard
* `SubjectDistance`  
  Distance in meters
* `SubjectDistanceRange`  
  Distance Type
  - 0 unknown
  - 1 macro
  - 2 close
  - 3 distant
* `Software`  
  Application and version used to generate the image
* `Copyright`  
  copyright information
