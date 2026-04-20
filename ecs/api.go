package ecs

// ========== 系统接口 ==========

// System 系统接口
type System interface {
	Initialize(w *World)
	Update(dt float64)
}

// ========== 查询器 ==========
