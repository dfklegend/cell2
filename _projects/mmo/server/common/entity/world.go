package entity

type World struct {
	IWorld

	nextId   EntityID
	entities map[EntityID]IEntity

	events  IWorldEventListener
	context IWorldContext
}

func NewWorld() *World {
	return &World{
		nextId:   EntityID(1),
		entities: map[EntityID]IEntity{},
	}
}

func (w *World) SetContext(ctx IWorldContext) {
	w.context = ctx
}

func (w *World) GetContext() IWorldContext {
	return w.context
}

func (w *World) SetEventListener(listener IWorldEventListener) {
	w.events = listener
}

func (w *World) AllocId() EntityID {
	id := w.nextId
	w.nextId++
	return id
}

func (w *World) AddEntity(e IEntity) {
	w.entities[e.GetId()] = e

	if w.events != nil {
		w.events.OnAddEntity(e)
	}
}

func (w *World) GetEntity(id EntityID) IEntity {
	e := w.entities[id]
	if e != nil && !e.IsDestroyed() {
		return e
	}
	return nil
}

func (w *World) removeEntity(e IEntity) {
	delete(w.entities, e.GetId())
}

func (w *World) Update() {
	var toDelete []IEntity

	for _, v := range w.entities {
		if !v.IsDestroyed() {
			v.Update()
		} else {
			if toDelete == nil {
				toDelete = make([]IEntity, 0)
			}
			toDelete = append(toDelete, v)
		}
	}

	for _, v := range w.entities {
		if !v.IsDestroyed() {
			v.LateUpdate()
		}
	}

	if toDelete == nil {
		return
	}

	for _, v := range toDelete {
		w.removeEntity(v)
	}
}

func (w *World) DestroyEntity(id EntityID) {
	e := w.GetEntity(id)
	if e == nil {
		return
	}
	if e.IsDestroyed() {
		return
	}

	if w.events != nil {
		w.events.OnDestroyEntity(e)
	}
	e.DoDestroy()
}

func (w *World) Destroy() {
	for k, _ := range w.entities {
		w.DestroyEntity(k)
	}

	w.Update()
}

func (w *World) Filter(doFunc func(entity IEntity)) {
	for _, v := range w.entities {
		if !v.IsDestroyed() {
			doFunc(v)
		}
	}
}
