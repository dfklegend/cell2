using System;

namespace Phoenix.Utils
{
    public static class BitUtil
    {
        public static int SetBits(ref int value, int bits)
        {
            return (value |= bits);
        }
        public static int ClrBits(ref int value, int bits)
        {
            return (value &= ~(bits));
        }
        public static bool IsBitOn(int value, int bit)
        {
            return (((value) & (bit)) != 0);
        }
        public static bool IsBitOff(int value, int bit)
        {
            return (((value) & (bit)) == 0);
        }
        public static bool BitMaskComp(int value, int bit, int mask)
        {
            return (((value) & (bit)) == mask);
        }
    }
} // namespace Phoenix.Utils