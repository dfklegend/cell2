using System;
using System.Collections.Generic;
using System.IO;
using System.Runtime.InteropServices;
using System.Text;

namespace Phoenix.Network
{
    public class MsgStream : MemoryStream
    {
        public MsgStream()
        {
        }

        public MsgStream(int capacity):base(capacity)
        {            
        }

        public MsgStream(byte[] bytes) : base(bytes)
        {
        }

        #region WriteStream
        public MsgStream Write(Int16 v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(Int32 v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(Int64 v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(UInt16 v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(UInt32 v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(UInt64 v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(float v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(double v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(bool v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(Char v)
        {
            return Write(BitConverter.GetBytes(v));
        }
        public MsgStream Write(Byte v)
        {
            base.WriteByte(v);
            return this;
        }
        public MsgStream Write(SByte v)
        {
            return Write((byte)v);
        }
        public MsgStream Write(string v, Encoding encoding)
        {
            return Write( encoding.GetBytes(v) );
        }
        public MsgStream Write(byte[] v)
        {
            base.Write(v, 0, v.Length);
            return this;
        }       
        #endregion

        #region ReadStream
        public SByte ReadInt8()
        {
            return ReadSByte();
        }
        public Int16 ReadInt16()
        {
            return BitConverter.ToInt16( ReadBytes(2), 0 );
        }
        public Int32 ReadInt32()
        {
            return BitConverter.ToInt32( ReadBytes(4), 0 );
        }
        public Int64 ReadInt64()
        {
            return BitConverter.ToInt64( ReadBytes(8), 0 );
        }
        public Byte ReadUInt8()
        {
            return ReadBYTE();
        }
        public UInt16 ReadUInt16()
        {
            return BitConverter.ToUInt16( ReadBytes(2), 0 );
        }
        public UInt32 ReadUInt32()
        {
            return BitConverter.ToUInt32( ReadBytes(4), 0 );
        }
        public UInt64 ReadUInt64()
        {
            return BitConverter.ToUInt64( ReadBytes(8), 0 );
        }
        public float ReadFloat()
        {
            return BitConverter.ToSingle( ReadBytes(4), 0 );
        }
        public double ReadDouble()
        {
            return BitConverter.ToDouble( ReadBytes(8), 0 );
        }
        public bool ReadBool()
        {
            return BitConverter.ToBoolean( ReadBytes(1), 0);
        }
        public char ReadChar()
        {
            return BitConverter.ToChar( ReadBytes(2), 0);
        }
        public SByte ReadSByte()
        {
            return (SByte)ReadBYTE();
        }
        public Byte ReadBYTE()
        {
            int b = base.ReadByte();
            if (b < 0)
                throw new Exception("MsgStream.ReadBYTE ReadEnd.");
            return (Byte)b;
        }
        public string ReadUTF8String()
        {
            return ReadString(Encoding.UTF8);
        }
        public string ReadUTF8String(UInt32 uLength)
        {
            return ReadString(uLength, Encoding.UTF8);
        }
        public string ReadString(Encoding srcEncoding)
        {
            List<byte> lt = new List<byte>();
            do 
            {
                int b = base.ReadByte();
                if (b > 0)
                    lt.Add((byte)b);
                else if(b == 0)
                    break;
                else
                    throw new Exception("MsgStream.ReadString ReadEnd.");
            } while (true);
            return srcEncoding.GetString(lt.ToArray());
        }
        public string ReadString(UInt32 uLength, Encoding srcEncoding)
        {
            return srcEncoding.GetString( ReadBytes(uLength) );
        }
        public byte[] ReadBytes(UInt32 uLength)
        {
            byte[] b = new byte[uLength];
            if( base.Read(b, 0, b.Length) <= 0 )
                throw new Exception("MsgStream.ReadBytes ReadEnd.");
            return b;
        }

        public MsgStream ReadStream()
        {
            UInt32 size = ReadUInt32();
            if( size == 0 )
                return new MsgStream();

            byte[] data = ReadBytes(size);
            MsgStream datas = new MsgStream();
            datas.Write(data, 0, data.Length);
            datas.Seek(0, System.IO.SeekOrigin.Begin);
            return datas;
        }

        #endregion

        #region 对应DynamicPacket的一些读写
        // DynamicPacket中的read/write
        // 16位长度
        public string DReadUTF8String16()
        {
            return DReadString16(Encoding.UTF8);
        }
        public string DReadString16(Encoding srcEncoding)
        {
            UInt16 uLength = ReadUInt16();
            if (uLength == 0)
                return string.Empty;
            return ReadString(uLength, srcEncoding);
        }
        public MsgStream DWriteUTF8String16(string str)
        {
            return DWriteString16(str, Encoding.UTF8);
        }
        public MsgStream DWriteString16(string str, Encoding dstEncoding)
        {
            byte[] byDst = dstEncoding.GetBytes(str);
            UInt16 len = (UInt16)(byDst.Length);
            Write((UInt16)len);
            if (len == 0)
            {
                //Write((byte)0);
                return this;
            }
            Write(byDst);
            //Write((byte)0);
            return this;
        }

        // 32位长度
        public string DReadUTF8String32()
        {
            return DReadString32(Encoding.UTF8);
        }
        public string DReadString32(Encoding srcEncoding)
        {
            UInt32 uLength = ReadUInt32();
            if (uLength == 0)
                return string.Empty;
            return ReadString(uLength, srcEncoding);
        }
        public MsgStream DWriteUTF8String32(string str)
        {
            return DWriteString32(str, Encoding.UTF8);
        }
        public MsgStream DWriteString32(string str, Encoding dstEncoding)
        {
            byte[] byDst = dstEncoding.GetBytes(str);
            UInt32 len = (UInt32)(byDst.Length);
            Write(len);
            if (len == 0)
            {                
                return this;
            }
            Write(byDst);            
            return this;
        }

        public void DReadVector<T>(out List<T> l)
            where T : struct
        {
            l = new List<T>();
            UInt16 uLength = ReadUInt16();
            T t = default(T);
            int sizeofT = Marshal.SizeOf(t);
            byte[] bytes = new byte[sizeofT];
            for (int iLoop = 0; iLoop < uLength; iLoop++)
            {
                if( base.Read(bytes, 0, sizeofT) <= 0)
                    throw new Exception("MsgStream.DReadVector ReadEnd.");
                if (typeof(T) == typeof(UInt64))
                {
                    object obj = BitConverter.ToUInt64(bytes, 0);
                    t = (T)obj;
                    l.Add(t);
                }
                else if (typeof(T) == typeof(UInt32))
                {
                    object obj = BitConverter.ToUInt32(bytes, 0);
                    t = (T)obj;
                    l.Add(t);
                }
                else
                {
                    throw new ArgumentNullException();
                }
            }
        }
        public MsgStream DWriteVector<T>(List<T> l)
            where T : struct
        {
            Write((UInt16)l.Count);
            for (int iLoop = 0; iLoop < l.Count; iLoop++)
            {
                if (typeof(T) == typeof(UInt64))
                {
                    object obj = l[iLoop];
                    Write(BitConverter.GetBytes((UInt64)obj));
                }
                else
                {
                    throw new ArgumentNullException();
                }
            }
            return this;
        }

        #endregion
    }
}
