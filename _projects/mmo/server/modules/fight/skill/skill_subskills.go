package skill

import (
	"github.com/dfklegend/cell2/utils/logger"
)

func (s *Skill) getSubSkills() []string {
	return s.cfg.SubSkills.Skills
}

func (s *Skill) hasSubSkills() bool {
	skills := s.getSubSkills()
	return skills != nil && len(skills) > 0
}

func (s *Skill) hasSubSkill(index int) bool {
	skills := s.getSubSkills()
	return index < len(skills)
}

func (s *Skill) checkCanBeginSubSkills() bool {
	maxStackDepth := 5
	if s.stackDepth >= maxStackDepth {
		logger.L.Warnf("subskills is too deep: %v", s.id)
		return false
	}
	return true
}

func (s *Skill) beginSubSkills() {
	s.startSubSkill(0)
}

func (s *Skill) createSubSkill(index int) *Skill {
	skills := s.getSubSkills()
	if index >= len(skills) {
		return nil
	}
	subSkillId := skills[index]

	return createSubSkill(s.owner, s, subSkillId, s.level, s.stackDepth+1)
}

func (s *Skill) updateSubSkills() {
	cur := s.subRunning
	if cur.IsFailed() {
		s.switchToFailed(ReasonSubFailed)
		return
	}

	if !cur.IsOver() {
		cur.Update()
	}

	if cur.IsOver() {
		if s.hasSubSkill(s.subRunningIndex + 1) {
			s.startSubSkill(s.subRunningIndex + 1)
		} else {
			s.doEnd()
		}
	}
}

func (s *Skill) startSubSkill(index int) {
	subSkill := s.createSubSkill(index)
	if subSkill == nil {
		s.switchToFailed(ReasonSubFailed)
		return
	}

	subSkill.Start()
	s.subRunningIndex = index
	s.subRunning = subSkill
}

func (s *Skill) TestBreakSubSkill() {
	if s.subRunning == nil {
		return
	}
	s.subRunning.Break()
}
