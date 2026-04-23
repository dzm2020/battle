package config

// BuffType 表示增益/减益的类型
type BuffType string

const (
	BuffTypeBuff   BuffType = "buff"   // 增益（正面效果）
	BuffTypeDebuff BuffType = "debuff" // 减益（负面效果）
)

type BufferEffectType int32

const (
	BufferEffectUndefined  BufferEffectType = iota
	BufferEffectStatChange                  // 属性变化（攻防等，数值见 Params）
	BufferEffectControl                     // 控制效果（眩晕、沉默等，子类型可由 ParamsString[0] 或 Params 约定）
	BufferEffectDamage                      // 伤害（如 DoT，强度）
	BufferEffectHeal                        // 治疗（如 HoT，强度）
)

type BufferEffect struct {
}

// BuffStackBehavior 同 BuffId 再次施加时的叠加策略（与 JSON/YAML 中的 stack_behavior 字符串对应）。
type BuffStackBehavior = int32

const (
	BuffStackUndefined BuffStackBehavior = iota
	BuffStackRefresh                     // 刷新持续时间，层数不变
	BuffStackAdd                         // 层数 +1（封顶见 MaxStack）并刷新持续
	BuffStackReplace                     // 用新实例替换已有槽位
	BuffStackIgnore                      // 已有同 Id 则忽略本次施加
)

// BuffConfig 表示一个独立的增益/减益配置
type BuffConfig struct {
	// 基础信息
	ID          uint32   `json:"id" yaml:"id"`                   // 唯一标识
	Name        string   `json:"name" yaml:"name"`               // 显示名称
	Type        BuffType `json:"type" yaml:"type"`               // buff 或 debuff
	Icon        string   `json:"icon" yaml:"icon"`               // UI 图标路径
	Description string   `json:"description" yaml:"description"` // 描述文本
	Dispellable bool     `json:"dispellable" yaml:"dispellable"` // 是否可被驱散
	// 生命周期
	DurationFrame int `json:"duration" yaml:"duration"` // 持续时间（帧），0 表示永久 用于控制buff生命周期
	// 叠加
	MaxStack      int               `json:"max_stack" yaml:"max_stack"`           // 最大叠加层数，1 表示不可叠加
	StackBehavior BuffStackBehavior `json:"stack_behavior" yaml:"stack_behavior"` // 叠加行为，见 [BuffStackRefresh] 等常量
	// 效果
	EffectType   BufferEffectType `json:"effect_type" yaml:"effect_type"`     // 属性名（如 attack_damage）
	Params       []float64        `json:"params" yaml:"params"`               // 效果参数  不同类型参数不同
	ParamsString []string         `json:"params_string" yaml:"params_string"` // 效果参数  不同类型参数不同
	CoolingFrame int              `json:"cooling_frame" yaml:"cooling_frame"` // 效果冷却帧每N帧生效一次

}
