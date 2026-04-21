package config

// BuffType 表示增益/减益的类型
type BuffType string

const (
	BuffTypeBuff   BuffType = "buff"   // 增益（正面效果）
	BuffTypeDebuff BuffType = "debuff" // 减益（负面效果）
)

// StatModifier 表示单个属性的修改
type StatModifier struct {
	Stat      Attribute `json:"stat" yaml:"stat"`             // 属性名（如 attack_damage）
	Delta     float64   `json:"delta" yaml:"delta"`           // 变化量（绝对值或百分比）
	IsPercent bool      `json:"is_percent" yaml:"is_percent"` // 是否为百分比
}

// BuffConfig 表示一个独立的增益/减益配置
type BuffConfig struct {
	ID          string   `json:"id" yaml:"id"`                   // 唯一标识，如 "atk_buff_1"
	Name        string   `json:"name" yaml:"name"`               // 显示名称
	Type        BuffType `json:"type" yaml:"type"`               // buff 或 debuff
	Icon        string   `json:"icon" yaml:"icon"`               // UI 图标路径
	Description string   `json:"description" yaml:"description"` // 描述文本

	// 效果参数
	Modifiers     []StatModifier `json:"modifiers" yaml:"modifiers"`           // 属性修改列表（支持多个）
	DurationFrame int            `json:"duration" yaml:"duration"`             // 持续时间（帧），0 表示永久
	MaxStack      int            `json:"max_stack" yaml:"max_stack"`           // 最大叠加层数，1 表示不可叠加
	StackBehavior string         `json:"stack_behavior" yaml:"stack_behavior"` // 叠加行为：refresh（刷新持续时间）、add（增加层数）、ignore（无效）
	Dispellable   bool           `json:"dispellable" yaml:"dispellable"`       // 是否可被驱散
}
