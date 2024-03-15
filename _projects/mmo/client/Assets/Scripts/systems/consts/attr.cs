
namespace Phoenix.Game
{
    // 属性定义，和服务器对应
    public enum eAttrType
    {
        Invalid = -1,
        Level,
        Side,
        HP,
        HPMax,
        Energy,
        EnergyMax,
        AttackSpeed,
        CriticalRate,
        PhysicPower,
        MagicPower,
        PhysicArmor,
        MagicArmor,
        WeaponMinDmg,
        WeaponMaxDmg,


        MaxAttrNum,
    }

    public static partial class Strings
    {
        public static string[] ATTR_NAMES = new string[]
        {
            "-",
            "-",
            "-",
            "生命",
            "-",
            "能量",
            "攻速",
            "暴击",
            "物伤",
            "法伤",
            "护甲",
            "法抗",
            "武器最小伤害",
            "武器最大伤害",
        };
    }
}