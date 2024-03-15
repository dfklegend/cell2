using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using System.IO;
using System.Text;
using Phoenix.Res;

namespace Phoenix.csv
{   
    /*
     * 数据表格的载入
     * 数据存储为utf-8编码
     * Excel表格的定义为:
     * 第一行,中文名字
     * 第二行,对应于T中实际的属性名
     * 
     * 
     * TODO: 看能否过滤掉key为null的数据
     * */
    public abstract class TableManager<TKey, TValue>
    {
        public abstract string TablePath();
        public abstract TKey MakeKey(TValue obj);
        // 是否需要map
        public virtual bool needMakeMap() { return true; }
        // 实际内容的行索引(-1是域对应)
        public abstract int indexRead();
        // the data arrays.
        TValue[] mItemArray;
        Dictionary<TKey, int> mKeyItemMap = new Dictionary<TKey, int>();

        // constructor.
        public TableManager()
        {            
        }

        void LoadTable()
        {
            Log.LogCenter.Asset.Debug("Begin load:" + TablePath());
            // load from excel txt file.
            try
            {
                mItemArray = CSVParser.ParseFromStr<TValue>(LoadFileContent(), indexRead());
            }
            catch (System.Exception e)
            {
                Debug.LogError(e);
                Log.LogCenter.Asset.Debug(e.ToString());
                mItemArray = new TValue[0];
            }

            if (mItemArray == null)
                mItemArray = new TValue[0];

            makeMap();
            
            Log.LogCenter.Asset.Debug(TablePath() + " read over!total row:" + mKeyItemMap.Count);
        }      

        private void makeMap()
        {
            if (!needMakeMap())
                return;
            this.mKeyItemMap.Clear();
            // build the key-value map.
            for (int i = 0; i < mItemArray.Length; i++)
            {
                TKey key = MakeKey(mItemArray[i]);
                if (key == null)
                {
                    //Debug.Log("find null key!");
                    continue;
                }
                try
                {
                    if(!mKeyItemMap.ContainsKey(key))
                        mKeyItemMap[key] = i;
                }
                catch (System.Exception e)
                {
#if UNITY_EDITOR
                    Log.LogCenter.Asset.Debug("skip empty key,line:{0}", i + indexRead() + 1);
#endif
                    Debug.LogError(e.Message +" "+ i);
                    continue;
                }
            }
        }

        // return false if @a <index> is out of range.
        public bool isValidIndex(int index)
        {
            if (index < 0 || index >= mItemArray.Length)
                return false;
            return true;
        }

        // get a item base the index.
        public TValue GetAt(int index)
        {
            if (isValidIndex(index))
            {
                return mItemArray[index];
            }
            return default(TValue);
        }

        // get a item base the key.
        public TValue GetItem(TKey key)
        {
            if( key == null )
                return default(TValue);
            int itemIndex;
            if (mKeyItemMap.TryGetValue(key, out itemIndex))
                return mItemArray[itemIndex];
            return default(TValue);
        }

        // get the item array.
        public TValue[] GetAllItem()
        {
            return mItemArray;
        }

        public void Load(bool bForce = false) 
        {
            if ( !bForce && mItemArray != null)
                return;
            LoadTable();
        }        

        protected virtual string LoadFileContent()
        {
            return CSVEnv.loader.LoadContent(TablePath());
        }
    }
}
