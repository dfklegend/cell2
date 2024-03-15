using System.Collections.Generic;
using Phoenix.Core;

namespace Phoenix.Game
{    
    // 抽象物品显示
    public class BaseItemStyleFactory
    {
        protected Dictionary<eItemStyle, System.Type> _styles =
            new Dictionary<eItemStyle, System.Type>();

        public virtual void Init()
        {           
        }

        public void Register(eItemStyle item, System.Type type)
        {
            _styles[item] = type;
        }

        public T Create<T>(eItemStyle style)
            where T:class
        {
            System.Type type;
            if (!_styles.TryGetValue(style, out type))
                return null;
            return System.Activator.CreateInstance(type) as T;
        }
    }

    public class ItemListStyleFactory : BaseItemStyleFactory
    {
        public static ItemListStyleFactory It 
        { 
            get { return Singleton<ItemListStyleFactory>.It; } 
        }

        public override void Init()
        {           
            Register(eItemStyle.Normal, typeof(ItemStyleNormal));
            Register(eItemStyle.Label, typeof(ItemStyleLabel));
        }        

        public BaseItemListStyle Create(eItemStyle style)
        {
            return Create<BaseItemListStyle>(style);
        }
    }

    public class ItemIconStyleFactory : BaseItemStyleFactory
    {
        public static ItemIconStyleFactory It
        {
            get { return Singleton<ItemIconStyleFactory>.It; }
        }

        public override void Init()
        {           
            Register(eItemStyle.Normal, typeof(NormalIconStyle));            
        }

        public BaseItemIconStyle Create(eItemStyle style)
        {
            return Create<BaseItemIconStyle>(style);
        }
    }
} // namespace Phoenix
