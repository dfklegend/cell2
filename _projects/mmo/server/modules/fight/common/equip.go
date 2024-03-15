package common

// 所有的可装备物体
// 要么增加属性，要么增加buf(一些动态机制)
// buf本身也可以增加属性

// IEquipable
// 比如装备，铭文，都可以使用此接口来实际增加属性
// 可以将装备对象构建一个临时此对象，调用
type IEquipable interface {
	Equip(character ICharacter)
	Unequip(character ICharacter)
}

// IEquipGroup
// 装备组，比如 所有装备，所有buf都可以作为一个group添加到character身上
// character初始化，只要重置所有属性，同时调一下所有group的equip即可重新计算属性
// 成组，避免细碎的装备添加
// groups
//		初始属性
//		装备
//		buf
//		...
// 装备组内部的具体装备的Equip/Unequip自己维护(使用上面的IEquipable)
type IEquipGroup interface {
	OnAdded(character ICharacter)
	Equip(character ICharacter)
	Unequip(character ICharacter)
}

type IEquipSlots interface {
	IEquipGroup
	SetEquip(index int, id EquipId)
	RemoveEquip(index int)
}
