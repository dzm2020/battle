package attr

// Base 基础属性：养成面板上的原始四维与等级。
// 战斗中攻防等不直接读 Base，一律经 Calculator 得到 Derived。
type Base struct {
	Level int32
	STR   int32 // 力量
	AGI   int32 // 敏捷
	INT   int32 // 智力
	VIT   int32 // 体质
}
