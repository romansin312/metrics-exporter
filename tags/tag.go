package tags

type TagModel struct {
	tagName  string
	tagValue string
}

func NewTagModel(tagName string, tagValue string) *TagModel {
	return &TagModel{
		tagName:  tagName,
		tagValue: tagValue,
	}
}

func (t *TagModel) GetName() string {
	return t.tagName
}

func (t *TagModel) GetValue() string {
	return t.tagValue
}
