package enum

import "encoding/json"

type ImageType int8

const (
	LocationType ImageType = 1
	QiNiuType    ImageType = 2
)

func (image ImageType) String() string {
	switch image {
	case LocationType:
		return "本地"
	case QiNiuType:
		return "七牛云"
	}

	return ""
}

func (image ImageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(image.String())
}
