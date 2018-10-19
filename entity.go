package yext

type EntityType string

type Entity interface {
	GetEntityId() string
	GetEntityType() EntityType
}

type EntityMeta struct {
	Id          *string           `json:"id,omitempty"`
	AccountId   *string           `json:"accountId,omitempty"`
	EntityType  EntityType        `json:"entityType,omitempty"`
	FolderId    *string           `json:"folderId,omitempty"`
	LabelIds    *UnorderedStrings `json:"labelIds,omitempty"`
	CategoryIds *[]string         `json:"categoryIds,omitempty"`
	Language    *string           `json:"language,omitempty"`
	CountryCode *string           `json:"countryCode,omitempty"`
}

type BaseEntity struct {
	Meta       *EntityMeta `json:"meta,omitempty"`
	nilIsEmpty bool
}

func (b *BaseEntity) GetEntityId() string {
	if b.Meta != nil && b.Meta.Id != nil {
		return *b.Meta.Id
	}
	return ""
}

func (b *BaseEntity) GetEntityType() EntityType {
	if b.Meta != nil {
		return b.Meta.EntityType
	}
	return ""
}

func (b *BaseEntity) GetNilIsEmpty() bool {
	return b.nilIsEmpty
}

func (b *BaseEntity) SetNilIsEmpty(val bool) {
	b.nilIsEmpty = val
}

type RawEntity map[string]interface{}

func (r *RawEntity) GetEntityId() string {
	if m, ok := (*r)["meta"]; ok {
		meta := m.(map[string]interface{})
		if id, ok := meta["id"]; ok {
			return id.(string)
		}
	}
	return ""
}

func (r *RawEntity) GetEntityType() EntityType {
	if m, ok := (*r)["meta"]; ok {
		meta := m.(map[string]interface{})
		if t, ok := meta["entityType"]; ok {
			return EntityType(t.(string))
		}
	}
	return EntityType("")
}
