using System.Collections.Generic;

namespace Phoenix.Game.FightEmulator.BagSystem
{	
    public class PlayerBags
    {
        private List<Bag> _bags = new List<Bag>();

        public void Init(int maxBag)
        {
            for (var i = 0; i < maxBag; i++)
                _bags.Add(null);
        }
        
        public void SetBag(int index, Bag bag)
        {
            if (index < 0 || index >= _bags.Count)
                return;
            _bags[index] = bag;
        }

        public Bag GetBag(int index)
        {
            if (index < 0 || index >= _bags.Count)
                return null;
            return _bags[index];
        }        
    }
} // namespace Phoenix
