package skill

// 技能过程
// Hit只是为了定义清楚，实际上不存在这个阶段
// 可以用子技能来实现复杂技能阶段
// 如果技能有子技能，那么所有子技能执行完毕，才算结束
const (
	Init int = iota
	Prefire
	Hit
	SubSkillRunning // 子技能执行过程中
	Postfire
	Over
	Failed // 执行失败
)

/*
	prefire->hit->postfire
				->SubSkillRunning
									->Over
									(cancel or break)->Failed
*/

// failed reason
const (
	ReasonCancel int = iota // 自己取消
	ReasonSubFailed
	ReasonBreak // 被打断
)
